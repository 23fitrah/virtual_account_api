package config

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	"github.com/olivere/elastic/v7"
)

type ESConfig struct {
	URL      string
	Username string
	Password string
	Timeout  time.Duration
}

var (
	esClient     *elastic.Client
	esClientLock sync.RWMutex
)

func SetEsClient(client *elastic.Client) {
	esClientLock.Lock()
	defer esClientLock.Unlock()
	esClient = client
}

func GetEsClient() *elastic.Client {
	esClientLock.RLock()
	defer esClientLock.RUnlock()
	return esClient
}

func NewElasticsearchConnection(config ESConfig) (*elastic.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	options := []elastic.ClientOptionFunc{
		elastic.SetURL(config.URL),
		elastic.SetSniff(false),
	}

	if config.Username != "" && config.Password != "" {
		options = append(options, elastic.SetBasicAuth(config.Username, config.Password))
	}

	client, err := elastic.NewClient(options...)
	if err != nil {
		logEntry := map[string]interface{}{
			"level":     "error",
			"timestamp": time.Now().Format(time.RFC3339),
			"message":   "Failed to connect to Elasticsearch: ",
			"error":     err.Error(),
		}
		jsonBytes, _ := sonic.Marshal(logEntry)

		return nil, fmt.Errorf(string(jsonBytes))
	}

	info, code, err := client.Ping(config.URL).Do(ctx)
	if err != nil {
		logEntry := map[string]interface{}{
			"level":     "error",
			"timestamp": time.Now().Format(time.RFC3339),
			"message":   "Failed to ping Elasticsearch",
			"error":     err.Error(),
		}
		jsonBytes, _ := sonic.Marshal(logEntry)

		return nil, fmt.Errorf(string(jsonBytes))
	}

	logEntry := map[string]interface{}{
		"level":     "info",
		"timestamp": time.Now().Format(time.RFC3339),
		"message":   fmt.Sprintf("Elasticsearch returned with code %d and version %s", code, info.Version.Number),
	}
	jsonBytes, _ := sonic.Marshal(logEntry)

	fmt.Println(string(jsonBytes))

	SetEsClient(client)
	return client, nil
}
