package iyzico

import (
	"bytes"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"time"
)

type Method string

const (
	POST Method = "POST"
	GET  Method = "GET"
)

func HTTP(method Method, url string, headers *map[string]string, body interface{}) ([]byte, error) {
	json, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(string(method), url, bytes.NewBuffer(json))
	if err != nil {
		return nil, err
	}

	if headers != nil {
		for key, value := range *headers {
			req.Header[key] = []string{value}
		}
	}

	client := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 5 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
