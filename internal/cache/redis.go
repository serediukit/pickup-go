package cache

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"pickup-srv/util"
)

type Cacher interface {
	Close() error
	Get(context.Context, string) (string, error)
	Set(context.Context, string, interface{}, time.Duration) error
}

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient() *RedisClient {
	host := util.GetEnv("REDIS_HOST", "redis")
	port := util.GetEnv("REDIS_PORT", "6379")
	addr := fmt.Sprintf("%s:%s", host, port)

	var client *redis.Client
	var err error

	const maxRetries = 10

	for i := 0; i < maxRetries; i++ {
		client = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: "",
			DB:       0,
		})

		ctx := context.Background()

		_, err = client.Ping(ctx).Result()
		if err == nil {
			log.Println("Redis connected")
			break
		}

		log.Printf("Redis connection attempt %d/%d failed: %v", i, maxRetries, err)
		time.Sleep(time.Duration(i+1) * 3 * time.Second)
	}

	if err != nil {
		log.Fatalf("Could not connect to Redis after %d attempts: %v", maxRetries, err)
	}

	return &RedisClient{client: client}
}

func (rc *RedisClient) Close() error {
	return rc.client.Close()
}

func (rc *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return rc.client.Get(ctx, key).Result()
}

func (rc *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return rc.client.Set(ctx, key, value, expiration).Err()
}
