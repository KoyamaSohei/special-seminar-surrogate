package main

import (
	"bytes"
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/go-redis/redis/v7"
)

var client *redis.Client

func initRedis() {
	h := os.Getenv("REDIS_HOST")
	p := os.Getenv("REDIS_PASS")
	client = redis.NewClient(&redis.Options{
		Addr:     h + ":6379",
		Password: p,
		DB:       0,
	})
	if client == nil {
		panic("redis: connection error")
	}
}

func setCache(key []byte, res *http.Response, b *bytes.Buffer) {
	println(string(b.Bytes()))
	var e time.Duration = 0
	if s := res.Header.Get("Surrogate-Control"); s != "" {
		re := regexp.MustCompile(`max-age=[0-9]+`)
		n, err := time.ParseDuration(re.FindString(s) + "ms")
		if err == nil {
			e = n
		}
	}
	ks := base64.StdEncoding.EncodeToString(key)
	err := client.Set(ks, b, e).Err()
	if err != nil {
		log.Println(err)
	}
}

func getCache(key []byte) ([]byte, error) {
	ks := base64.StdEncoding.EncodeToString(key)
	va, err := client.Get(ks).Result()
	if err != nil {
		return nil, err
	}
	return []byte(va), nil
}
