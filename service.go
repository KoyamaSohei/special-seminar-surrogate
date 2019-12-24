package main

import (
	"bufio"
	"io"
	"log"
	"net"
)

func serveConn(c net.Conn) {
	br := bufio.NewReader(c)
	rq := httpHostHeader(br)
	h := rq.Host
	if rq == nil || h == "" {
		return
	}
	log.Println(h)
	ret := make(chan net.IP)
	go resolveName(h+".", ret)
	ip := <-ret
	if ip == nil {
		log.Println("ip not found")
		return
	}
	if n := br.Buffered(); n > 0 {
		peeked, _ := br.Peek(br.Buffered())
		ca, err := getCache(peeked)
		if err != nil {
			handleConn(c, peeked, ip, h)
		} else {
			c.Write(ca)
		}
	}
}

func handleConn(c net.Conn, p []byte, ip net.IP, h string) {
	log.Println(string(p))
	rConn, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: ip, Port: 80, Zone: h})
	if err != nil {
		log.Println(err)
		return
	}
	defer rConn.Close()
	if _, err := rConn.Write(p); err != nil {
		log.Println(err)
		return
	}
	buf := make([]byte, 0xffff)
	ok := true
	ca := make([]byte, 0)
	for {
		n, err := rConn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Println(err)
				ok = false
			}
			break
		}
		b := buf[:n]
		n, err = c.Write(b)
		ca = append(ca, b...)
		if err != nil {
			log.Println(err)
			ok = false
			break
		}
	}
	if ok && len(ca) > 0 {
		go setCache(p, ca)
	}
}

func serveSurrogate() {
	initRedis()
	ln, err := net.Listen("tcp", "0.0.0.0:80")
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
