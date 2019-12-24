package main

import (
	"bufio"
	"bytes"
	"net/http"
)

func httpHostHeader(br *bufio.Reader) *http.Request {
	const maxPeek = 4 << 16
	peekSize := 0
	for {
		peekSize++
		if peekSize > maxPeek {
			return nil
		}
		b, err := br.Peek(peekSize)
		if n := br.Buffered(); n > peekSize {
			b, _ = br.Peek(n)
			peekSize = n
		}
		if len(b) > 0 {
			if b[0] < 'A' || b[0] > 'Z' {
				return nil
			}
			if bytes.Index(b, crlfcrlf) != -1 || bytes.Index(b, lflf) != -1 {
				req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(b)))
				if err != nil {
					return nil
				}
				if len(req.Header["Host"]) > 1 {
					return nil
				}
				return req
			}
		}
		if err != nil {
			return nil
		}
	}
}

var (
	crlfcrlf = []byte("\r\n\r\n")
	lflf     = []byte("\n\n")
)
