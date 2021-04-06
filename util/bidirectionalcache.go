package util

import (
	"github.com/patrickmn/go-cache"
	"time"
)

type BiCache struct {
	Forward  *cache.Cache
	Backward *cache.Cache
}

func NewBiCache(defaultExpiration, cleanupInterval time.Duration) *BiCache {
	return &BiCache{
		Forward:  cache.New(defaultExpiration, cleanupInterval),
		Backward: cache.New(defaultExpiration, cleanupInterval),
	}
}

func (c *BiCache) ContainsKey(k string) bool {
	_, ok := c.Forward.Get(k)
	return ok
}

func (c *BiCache) ContainsValue(v string) bool {
	_, ok := c.Backward.Get(v)
	return ok
}

// No checking of duplicate keys / values bc its being used in the context of UUID to Codes. UUIDs are guaranteed to be
// unique and there's a 1 / 1.8 billion (1 / 35^6) chance of the codes being the same so it should be alright?
func (c *BiCache) Set(k string, v string, d time.Duration) {
	c.Forward.Set(k, v, d)
	c.Backward.Set(v, k, d)
}

func (c *BiCache) SetDefault(k string, v string) {
	c.Forward.SetDefault(k, v)
	c.Backward.SetDefault(v, k)
}

func (c *BiCache) GetValue(k string) (string, bool) {
	v, ok := c.Forward.Get(k)
	return v.(string), ok
}

func (c *BiCache) GetKey(v string) (string, bool) {
	k, ok := c.Backward.Get(v)
	return k.(string), ok
}

func (c *BiCache) GetValueWithExpiration(k string) (string, time.Time, bool) {
	v, expirationTime, ok := c.Forward.GetWithExpiration(k)
	return v.(string), expirationTime, ok
}

func (c *BiCache) GetKeyWithExpiration(v string) (string, time.Time, bool) {
	k, expirationTime, ok := c.Backward.GetWithExpiration(v)
	return k.(string), expirationTime, ok
}

func (c *BiCache) GetKeyAndDelete(v string) (string, bool) {
	k, ok := c.Backward.Get(v)
	if ok {
		c.Delete(k.(string), v)
	}

	return k.(string), ok
}

func (c *BiCache) GetValueAndDelete(k string) (string, bool) {
	v, ok := c.Forward.Get(k)
	if ok {
		c.Delete(k, v.(string))
	}

	return v.(string), ok
}

func (c *BiCache) Delete(k, v string) {
	c.Forward.Delete(k)
	c.Backward.Delete(v)
}

func (c *BiCache) Flush() {
	c.Forward.Flush()
	c.Backward.Flush()
}
