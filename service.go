package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"

	"go.uber.org/zap"
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
		logger.Error("error occured when parse http request.", zap.Error(err))
		return
	}
	h := rq.Host
	if h == "" {
		logger.Info("host is empty")
		return
	}
	k := genHash(rq)
	ret := make(chan net.IP)
	go resolveName(h+".", ret)
	ip := <-ret
	if ip == nil {
		logger.Info("ip not found for " + h)
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
		logger.Error("error occured when listen "+h+":80 ,"+ip.String(), zap.Error(err))
		return
	}
	defer func() {
		rConn.Close()
		c.Close()
	}()
	if err := rq.WriteProxy(rConn); err != nil {
		logger.Error("error on writing byte", zap.Error(err))
		return
	}
	br := bufio.NewReader(rConn)
	res, err := http.ReadResponse(br, rq)
	if err != nil {
		logger.Error("response parse error", zap.Error(err))
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
		logger.Error("tcp listen error", zap.Error(err))
		return
	}
	for {
		c, err := ln.Accept()
		if err != nil {
			logger.Error("error occured when connect to client", zap.Error(err))
			return
		}
		go serveConn(c)
	}
}
