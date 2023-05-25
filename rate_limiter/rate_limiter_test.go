package main

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

type RateLimiter struct {
	rate int
	burst int
	tokens int
	ticker *time.Ticker
	mtx sync.Mutex
}

func NewRateLimiter(rate int, burst int) *RateLimiter {
	rl := &RateLimiter{
		rate: rate,
		burst: burst,
		tokens: burst,
		ticker: time.NewTicker(time.Second * time.Duration(rate)),
	}

	go func() {
		for range rl.ticker.C {
			rl.mtx.Lock()
			if rl.tokens < rl.burst {
				rl.tokens++
			}
			rl.mtx.Unlock()
		}
	}()

	return rl
}

func (rl *RateLimiter) Allow() bool {
	rl.mtx.Lock()
	defer rl.mtx.Unlock()

	if rl.tokens > 0 {
		rl.tokens--
		return true
	}
	return false
}

func TestRateLimiter(t *testing.T) {
	rl := NewRateLimiter(5, 5)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !rl.Allow() {
			http.Error(w, "Too many Requests", http.StatusTooManyRequests)
			return
		}
		w.Write([]byte("Hello World"))
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	for i:= 0;  i < 10; i++ {
		resp, err := http.Get(server.URL)
		if err != nil {
			t.Fatal(err)
		}

		if i < 5 && resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200 for request %d, got %d", i+1, resp.StatusCode)
		} else if i >= 5 && resp.StatusCode != http.StatusTooManyRequests {
			t.Errorf("Expected status 429 for request %d, got %d", i+1, resp.StatusCode)
		}

		time.Sleep(500 * time.Millisecond)
	}
}
