package server

import (
	"github.com/benbjohnson/immutable"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

const perMinLimit = 1000

// RateLimiter object
type RateLimiter struct {
	ServerMap          sync.Map
	Ticker             *time.Ticker
	WhitelistedServers *immutable.List
	CloseChan          chan struct{}
}

// Creates a new throttler with whitelisted servers
func NewThrottler() *RateLimiter {
	throttler := newEmptyThrottler()
	throttler.WhitelistedServers = getWhitelistedServers()
	return throttler
}

// Creates an empty throttler
func newEmptyThrottler() *RateLimiter {
	return &RateLimiter{
		ServerMap: sync.Map{},
		Ticker:    time.NewTicker(60 * time.Second),
		CloseChan: make(chan struct{}, 1),
	}
}

func getWhitelistedServers() *immutable.List {
	builder := immutable.NewListBuilder()
	builder.Append("hi")
	return builder.List()
}

// Increases the rate of throttling by 1
// Returns whether new increased throttle is greater than the threshold
func (t *RateLimiter) IncreaseRate(id string) bool {
	newRate := t.addRate(id)
	return newRate > perMinLimit
}

// Clears rates over given time period
// Also contains close chan
func (t *RateLimiter) Run() {
	go func() {
		for {
			select {
			case <-t.Ticker.C:
				t.clearRates()
			case <-t.CloseChan:
				t.Ticker.Stop()
				t.ServerMap.Range(func(key, value interface{}) bool {
					_, ok := t.ServerMap.LoadAndDelete(key)
					return ok
				})
				return
			}
		}
	}()
}

// Sends to close chan
func (t *RateLimiter) Close() {
	t.CloseChan <- struct{}{}
}

func (t *RateLimiter) addRate(id string) int {
	currRateInterface, loadOk := t.ServerMap.Load(id)
	if !loadOk {
		t.ServerMap.Store(id, 0)
		return 0
	} else {
		currRate, castOk := currRateInterface.(int)

		if !castOk {
			logrus.Fatalf("Unable to cast value to int (shouldn't see this)")
		}

		newRate := currRate + 1
		t.ServerMap.Store(id, newRate)
		return newRate
	}
}

func (t *RateLimiter) clearRates() {
	logrus.Info("Clearing limiter...")
	t.ServerMap.Range(func(key, value interface{}) bool {
		t.ServerMap.Store(key, 0)
		return true
	})
}
