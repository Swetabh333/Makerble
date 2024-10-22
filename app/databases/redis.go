package databases

import (
	"context"
	"crypto/tls"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"strings"
)

// Function to connect to redis
func ConnectToRedis() (*redis.Client, error) {

	var tlsConfig *tls.Config
	if strings.Contains(os.Getenv("REDIS_URL"), "upstash.io") {
		tlsConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	client := redis.NewClient(&redis.Options{
		Addr:      os.Getenv("REDIS_URL"),
		Password:  os.Getenv("REDIS_PASSWORD"),
		DB:        0, // Use default DB
		TLSConfig: tlsConfig,
	})
	err := client.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	log.Println("Connected to Redis successfully")
	return client, nil
}
