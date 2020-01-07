package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
)

func genHash(req *http.Request) []byte {
	h := sha256.New()
	s := fmt.Sprintf("%s--%s--%s", req.Host, req.Method, req.URL.RequestURI())
	h.Write([]byte(s))
	return h.Sum(nil)
}

func serveConn(c net.Conn) {
	br := bufio.NewReader(c)
	rq, err := http.ReadRequest(br)
	if err != nil {
		return
	}
	h := rq.Host
	if h == "" {
		return
	}
	k := genHash(rq)
	ret := make(chan net.IP)
	go resolveName(h+".", ret)
	ip := <-ret
	if ip == nil {
		log.Println("ip not found")
		return
	}
	ca, err := getCache(k)
	if err != nil {
		handleConn(c, ip, h, rq, k)
	} else {
		c.Write(ca)
	}
}

func handleConn(c net.Conn, ip net.IP, h string, rq *http.Request, key []byte) {
	rConn, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: ip, Port: 80, Zone: h})
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		rConn.Close()
		c.Close()
	}()
	if err := rq.WriteProxy(rConn); err != nil {
		log.Println(err)
		return
	}
	br := bufio.NewReader(rConn)
	res, err := http.ReadResponse(br, rq)
	if err != nil {
		log.Println(err)
		return
	}
	b := new(bytes.Buffer)
	wt := io.MultiWriter(c, b)
	res.Write(wt)
	if rq.Method == "GET" {
		go setCache(key, res, b)
	}

}

func serveSurrogate() {
	initRedis()
	h := os.Getenv("PROXY_HOST")
	p := os.Getenv("PROXY_PORT")
	if h == "" {
		h = "0.0.0.0"
	}
	if p == "" {
		p = "80"
	}
	ln, err := net.Listen("tcp", h+":"+p)
	if err != nil {
		log.Fatal(err)
		return
	}
	for {
		c, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
			return
		}
		go serveConn(c)
	}
}
