package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/bytedance/sonic"
)

var AppLogger *log.Logger
var logWriter *DailyLogWriter

var logChan = make(chan string, 1000)

type DailyLogWriter struct {
	mu          sync.Mutex
	currentDate string
	file        *os.File
	dir         string
	filename    string
}

func NewDailyLogWriter(dir, filename string) *DailyLogWriter {
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}
	return &DailyLogWriter{
		dir:      dir,
		filename: filename,
	}
}

func (w *DailyLogWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	today := time.Now().Format("2006-01-02")
	if today != w.currentDate {
		if err := w.rotate(today); err != nil {
			return 0, err
		}
	}

	return w.file.Write(p)
}

func (w *DailyLogWriter) rotate(today string) error {
	if w.file != nil {
		w.file.Close()
	}

	filename := fmt.Sprintf("%s-%s.json", w.filename, today)
	path := filepath.Join(w.dir, filename)

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	w.file = file
	w.currentDate = today
	return nil
}

func init() {
	logWriter = NewDailyLogWriter("logs", "app")

	multiWriter := io.MultiWriter(os.Stdout, logWriter)

	log.SetOutput(multiWriter)

	AppLogger = log.New(multiWriter, "", 0)

	go startLogWorker()
}

func startLogWorker() {
	for msg := range logChan {
		AppLogger.Println(msg)
	}
}

func LogJSON(data interface{}) {
	jsonBytes, err := sonic.Marshal(data)
	if err != nil {
		fallback := fmt.Sprintf(`{"level":"error","message":"Failed to marshal log data: %v"}`, err)
		select {
		case logChan <- fallback:
		default:
			// Jika channel penuh, drop log atau print ke stderr darurat (opsional)
			fmt.Println("Log channel full, dropping message")
		}
		return
	}

	select {
	case logChan <- string(jsonBytes):
	default:
		// Opsi: Drop log jika antrian penuh agar performa aplikasi utama tidak terganggu
		fmt.Println("Log channel full, dropping log message")
	}
}

func LogError(message string, err error) {
	logEntry := map[string]interface{}{
		"level":     "error",
		"timestamp": time.Now().Format(time.RFC3339),
		"message":   message,
		"error":     err.Error(),
	}
	LogJSON(logEntry)
}

func LogInfo(message string) {
	logEntry := map[string]interface{}{
		"level":     "info",
		"timestamp": time.Now().Format(time.RFC3339),
		"message":   message,
	}

	LogJSON(logEntry)
}
