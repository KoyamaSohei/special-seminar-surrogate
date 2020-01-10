package main

import (
	"github.com/go-redis/redis/v7"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"regexp"
	"syscall"
)

var (
	logger *zap.Logger
	client *redis.Client
	maxage *regexp.Regexp
	ip     *regexp.Regexp
	ds     string
)

func initConf() {
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
	maxage = regexp.MustCompile(`max-age=[0-9]+`)
	ds = os.Getenv("DNS_SERVER")
	if ds == "" {
		panic("DNS_SERVER not found")
	}
}

func main() {
	ops := zap.NewProductionConfig()
	ops.OutputPaths = []string{"stdout"}
	logger, _ = ops.Build()
	initConf()
	serveSurrogate()
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	log.Printf("Signal (%s) received, stopping\n", s)
}
