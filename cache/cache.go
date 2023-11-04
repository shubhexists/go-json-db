package cache

import (
	"github.com/patrickmn/go-cache"
	"time"
)
// We will be caching the data for model methods like ReadAll, Read, etc.

// This function would cache the data for a specified time provided by the user.
//Start cache
func StartCache(expiration int, purge int) *cache.Cache {
	return cache.New(time.Duration(expiration)*time.Minute, time.Duration(purge)*time.Minute)
}

//Set cache
func SetCache(c *cache.Cache, key string, value interface{}) {
	c.Set(key, value, cache.DefaultExpiration)
}

//Get cache
func GetCache(c *cache.Cache, key string) (interface{}, bool) {
	return c.Get(key)
}

//Delete cache
func DeleteCache(c *cache.Cache, key string) {
	c.Delete(key)
}

//Flush cache - Deletes all the keys
func FlushCache(c *cache.Cache) {
	c.Flush()
}
