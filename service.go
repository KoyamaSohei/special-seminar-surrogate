package main

import (
	"bufio"
	"io"
	"log"
	"net"
)

func serveConn(c net.Conn) {
	br := bufio.NewReader(c)
	hh := httpHostHeader(br)
	log.Println(hh)
	ret := make(chan net.IP)
	go resolveName(hh+".", ret)
	ip := <-ret
	if ip == nil {
		log.Println("ip not found")
		return
	}
	if n := br.Buffered(); n > 0 {
		peeked, _ := br.Peek(br.Buffered())
		handleConn(c, peeked, ip, hh)
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
	for {
		n, err := rConn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			break
		}
		b := buf[:n]
		n, err = c.Write(b)
		if err != nil {
			log.Println(err)
			break
		}
	}

}

func serveSurrogate() {
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
