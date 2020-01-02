package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	serveEcho()
	serveSurrogate()
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	log.Printf("Signal (%s) received, stopping\n", s)
}
