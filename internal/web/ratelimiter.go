package ratelimiter

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang/glog"
)

const maxPerHour = 500

func getIPFromRemoteAddr(remoteAddr string) string {
	comps := strings.Split(remoteAddr, ":")
	if len(comps) != 2 {
		// If the address is not parsable, just return localhost
		return "localhost"
	}
	return comps[0]
}

type hourRate struct {
	hour  int
	count int
}

type RateLimiter struct {
	rates map[string]*hourRate
	mu    sync.RWMutex
}

func New() *RateLimiter {
	return &RateLimiter{rates: make(map[string]*hourRate)}
}

func (r *RateLimiter) ShouldAllow(req *http.Request) bool {
	if !strings.Contains(req.URL.Path, "menus") {
		return true
	}
	hour := time.Now().Hour()
	ip := getIPFromRemoteAddr(req.RemoteAddr)
	r.mu.Lock()
	rate, ok := r.rates[ip]
	if !ok || rate.hour != hour {
		r.rates[ip] = &hourRate{hour, 0}
		rate = r.rates[ip]
	}
	var requestRate = rate.count + 1
	rate.count = requestRate
	r.mu.Unlock()
	glog.Infof("RateLimiter %s {hour: %v, count: %v} Allow: %b", ip, hour, requestRate, requestRate <= maxPerHour)
	return requestRate <= maxPerHour
}
