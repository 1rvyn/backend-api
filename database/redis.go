package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisInstance struct {
	Client *redis.Client
}

var ctx = context.Background()

var Redis RedisInstance

var RedisHost = os.Getenv("REDIS_HOST")
var RedisPassword = os.Getenv("REDIS_PASSWORD")

func ConnectRedis() {
	fmt.Println("Connecting to Redis")
	client := redis.NewClient(&redis.Options{
		Addr:     RedisHost,
		Password: RedisPassword, // no password set
		//Addr:     "localhost:6379",
		//Password: "",
		DB: 0, // use default DB
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to Redis")

	Redis = RedisInstance{Client: client}

}

func (r *RedisInstance) Put(key string, value string) error {
	fmt.Printf("Putting value %s in Redis for key %s)\n", value, key)
	err := r.Client.Set(ctx, key, value, 0).Err()
	if err != nil {
		fmt.Println("Error putting value in Redis", err)
		return err
	}
	fmt.Printf("Successfully added %s with value %s\n", key, value)
	return nil
}

func (r *RedisInstance) Get(key string) (string, error) {
	fmt.Printf("Getting value from Redis for key %s)\n", key)
	value, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		fmt.Println("Error getting value from Redis\n", err)
		return "", err
	}
	fmt.Printf("Successfully retrieved value %s for key %s\n", value, key)
	return value, nil
}

func (r *RedisInstance) PutHMap(key string, value map[string]interface{}) error {
	// fmt.Printf("Putting HashMap %s in Redis with the key %s)\n", value, key)
	err := r.Client.HMSet(ctx, key, value).Err()
	if err != nil {
		fmt.Println("Error putting |HashMap| in Redis", err)
		return err
	}

	err = r.Client.Expire(ctx, key, 24*time.Hour).Err() // Set TTL to 24 hours
	if err != nil {
		fmt.Println("Error setting |TTL| in Redis", err)
		return err
	}
	//fmt.Printf("Successfully added %s with value %v\n", key, value)
	return nil
}

func (r *RedisInstance) GetHMap(key string) (map[string]string, error) {
	fmt.Printf("Getting HashMAP from Redis for key %s)\n", key)
	value, err := r.Client.HGetAll(ctx, key).Result()
	if err != nil {
		fmt.Println("Error getting |HashMAP| from Redis", err)
		return nil, err
	}
	//fmt.Printf("Successfully retrieved value %v for key %s\n", value, key)
	return value, nil
}
