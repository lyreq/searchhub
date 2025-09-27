package provider

import (
	"context"
	"errors"
	"io"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"time"
)

func init() { rand.Seed(time.Now().UnixNano()) }

type RetryConfig struct {
	MaxRetries int
	BaseDelay  time.Duration
	MaxDelay   time.Duration
}

var defaultRetry = RetryConfig{
	MaxRetries: 4, // toplam 5 deneme: 0..4
	BaseDelay:  200 * time.Millisecond,
	MaxDelay:   5 * time.Second,
}

func Do(ctx context.Context, providerName string, req *http.Request, client *http.Client) (*http.Response, error) {
	if client == nil {
		client = HTTPClient()
	}
	lim := Limiter(providerName)
	cfg := defaultRetry

	var lastErr error
	for attempt := 0; attempt <= cfg.MaxRetries; attempt++ {
		// rate-limit
		if err := lim.Wait(ctx); err != nil {
			return nil, err
		}

		resp, err := client.Do(req)
		if err != nil {
			if isRetriableErr(err) && attempt < cfg.MaxRetries && ctx.Err() == nil {
				time.Sleep(backoffDelay(cfg, attempt))
				lastErr = err
				continue
			}
			return nil, err
		}

		// 2xx -> başarı
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return resp, nil
		}

		// retriable değilse veya hakkımız bitti ise doğrudan döndür
		if !isRetriableStatus(resp.StatusCode) || attempt >= cfg.MaxRetries {
			return resp, nil
		}

		// gövdeyi boşaltıp kapat (retry öncesi)
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()

		// Retry-After varsa saygı göster
		delay := retryAfterDelay(resp)
		if delay <= 0 {
			delay = backoffDelay(cfg, attempt)
		}
		if delay > cfg.MaxDelay {
			delay = cfg.MaxDelay
		}

		t := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			t.Stop()
			return nil, ctx.Err()
		case <-t.C:
			// continue
		}
	}
	if lastErr != nil {
		return nil, lastErr
	}
	return nil, context.DeadlineExceeded
}

// ---- helpers ----

func isRetriableStatus(code int) bool {
	if code == http.StatusTooManyRequests { // 429
		return true
	}
	return code == 500 || code == 502 || code == 503 || code == 504
}

func isRetriableErr(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	var ne net.Error
	if errors.As(err, &ne) {
		return ne.Timeout() || ne.Temporary()
	}
	return errors.Is(err, io.ErrUnexpectedEOF)
}

func retryAfterDelay(resp *http.Response) time.Duration {
	val := resp.Header.Get("Retry-After")
	if val == "" {
		return 0
	}
	// seconds
	if secs, err := strconv.Atoi(val); err == nil && secs >= 0 {
		return time.Duration(secs) * time.Second
	}
	// HTTP-date
	if when, err := http.ParseTime(val); err == nil {
		d := time.Until(when)
		if d < 0 {
			return 0
		}
		return d
	}
	return 0
}

func backoffDelay(cfg RetryConfig, attempt int) time.Duration {
	// expo backoff + jitter (±30%)
	d := cfg.BaseDelay << attempt
	if d > cfg.MaxDelay {
		d = cfg.MaxDelay
	}
	jitter := 0.3
	min := float64(d) * (1 - jitter)
	max := float64(d) * (1 + jitter)
	return time.Duration(min + rand.Float64()*(max-min))
}
