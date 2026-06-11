package utils

import (
	"context"
	"time"

	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
)

func GetOrSet[T any](
	ctx context.Context,
	rdb *redis.Client,
	key string,
	ttl time.Duration,
	fetchFunc func() (T, error),
) (T, error) {
	var data T

	val, err := rdb.Get(ctx, key).Result()
	if err == nil {
		errUnmarshal := sonic.Unmarshal([]byte(val), &data)
		if errUnmarshal == nil {
			return data, nil
		}
	} else if err != redis.Nil {
		// log redis
	}

	data, err = fetchFunc()
	if err != nil {
		return data, err
	}

	marshaledData, errMarshal := sonic.Marshal(data)

	go func() {
		if errMarshal == nil {
			rdb.Set(context.Background(), key, marshaledData, ttl)
		}
	}()

	return data, nil
}

func GetOrSetWithValidation[T any](
	ctx context.Context,
	rdb *redis.Client,
	key string,
	ttl time.Duration,
	fetchFunc func() (T, error),
	validateFunc func(T) bool,
) (T, error) {
	var data T
	var cacheHit bool

	val, err := rdb.Get(ctx, key).Result()
	if err == nil {
		if errUnmarshal := sonic.Unmarshal([]byte(val), &data); errUnmarshal == nil {
			cacheHit = true
		}
	}

	if cacheHit {
		isValid := validateFunc(data)
		if isValid {
			return data, nil
		}
	}

	data, err = fetchFunc()
	if err != nil {
		return data, err
	}

	marshaledData, errMarshal := sonic.Marshal(data)

	go func() {
		if errMarshal == nil {
			rdb.Set(context.Background(), key, marshaledData, ttl)
		}
	}()

	return data, nil
}

func InvalidateCache(ctx context.Context, rdb *redis.Client, key string) error {
	return rdb.Del(ctx, key).Err()
}
