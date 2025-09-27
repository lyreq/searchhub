package provider

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

var (
	rlMu         sync.RWMutex
	limiters     = map[string]*rate.Limiter{}
	defaultRPM   = 60
	defaultBurst = 5
)

func SetDefaultRPM(rpm int) {
	if rpm <= 0 {
		return
	}
	rlMu.Lock()
	defaultRPM = rpm
	rlMu.Unlock()
}

// SetRateLimit sets (or overwrites) a provider-specific limiter.
func SetRateLimit(providerName string, rpm int, burst int) {
	if rpm <= 0 {
		rpm = defaultRPM
	}
	if burst <= 0 {
		burst = defaultBurst
	}
	rlMu.Lock()
	defer rlMu.Unlock()
	// tokens per minute -> interval between tokens
	interval := time.Minute / time.Duration(rpm)
	limiters[providerName] = rate.NewLimiter(rate.Every(interval), burst)
}

// Limiter returns the limiter for providerName, creating a default if missing.
func Limiter(providerName string) *rate.Limiter {
	rlMu.RLock()
	lim, ok := limiters[providerName]
	rlMu.RUnlock()
	if ok {
		return lim
	}
	SetRateLimit(providerName, defaultRPM, defaultBurst)
	rlMu.RLock()
	l := limiters[providerName]
	rlMu.RUnlock()
	return l
}
