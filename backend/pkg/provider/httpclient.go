package provider

import (
	"net"
	"net/http"
	"sync"
	"time"
)

var (
	defaultClient *http.Client
	once          sync.Once
)

// HTTPClient returns a shared, tuned *http.Client.
func HTTPClient() *http.Client {
	once.Do(func() {
		transport := &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   10,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   5 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}
		defaultClient = &http.Client{
			Timeout:   7 * time.Second,
			Transport: transport,
		}
	})
	return defaultClient
}
