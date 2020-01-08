package main

import (
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var logger *zap.Logger

func main() {
	ops := zap.NewProductionConfig()
	ops.OutputPaths = []string{"stdout"}
	logger, _ = ops.Build()
	serveEcho()
	serveSurrogate()
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	log.Printf("Signal (%s) received, stopping\n", s)
}
