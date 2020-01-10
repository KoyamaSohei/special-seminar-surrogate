package main

import (
	"bytes"
	"encoding/base64"
	"net/http"
	"time"

	"go.uber.org/zap"
)

func setCache(key []byte, res *http.Response, b *bytes.Buffer) {
	logger.Info(string(b.Bytes()))
	var e time.Duration = 0
	if s := res.Header.Get("Surrogate-Control"); s != "" {
		n, err := time.ParseDuration(maxage.FindString(s) + "s")
		if err == nil {
			e = n
		}
	}
	logger.Info("expire at " + e.String())
	ks := base64.StdEncoding.EncodeToString(key)
	err := client.Set(ks, b, e).Err()
	if err != nil {
		logger.Error("error occured when set cache", zap.Error(err))
	}
}

func getCache(key []byte) ([]byte, error) {
	ks := base64.StdEncoding.EncodeToString(key)
	va, err := client.Get(ks).Result()
	if err != nil {
		logger.Error("cache not found", zap.Error(err))
		return nil, err
	}
	logger.Info("cache found")
	return []byte(va), nil
}
