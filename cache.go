package main

import (
	"log"

	"github.com/go-redis/redis/v7"
)

var client *redis.Client

func initRedis() {
	client = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	if client == nil {
		panic("redis: connection error")
	}
}

func setCache(key []byte, content []byte) {
	err := client.Set(string(key), string(content), 0).Err()
	if err != nil {
		log.Println(err)
	}
}

func getCache(key []byte) ([]byte, error) {
	va, err := client.Get(string(key)).Result()
	if err != nil {
		return nil, err
	}
	return []byte(va), nil
}
