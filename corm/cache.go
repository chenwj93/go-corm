package corm

import (
	"time"
	"encoding/json"
	"sync"
	"utils"
)

type CacheValue struct {
	value     []byte
	catchTime time.Time
	len       int
}

type Cache struct {
	cache map[string]CacheValue
	mutex sync.RWMutex
}

func NewCache() Cache{
	return Cache{make(map[string]CacheValue), sync.RWMutex{}}
}

func (c *Cache)SetCache(key interface{}, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	keyJ, e := json.Marshal(key)
	valJ, e := json.Marshal(value)
	if e != nil{
		return
	}
	k := utils.Md5(string(keyJ))
	c.cache[k] = CacheValue{value:valJ, catchTime:time.Now(), len:len(valJ)}
	return
}

func (c *Cache)GetCache(key interface{})(val []byte, ok bool){
	keyJ, e := json.Marshal(key)
	if e != nil{
		return nil, false
	}
	k := utils.Md5(string(keyJ))

	c.mutex.RLock()
	cacheVal, ok := c.cache[k]
	c.mutex.RUnlock()
	if ok {
		return cacheVal.value, ok
	}
	return nil, ok
}

func (c *Cache)EmptyCache(){
	c.mutex.Lock()
	for k := range c.cache{
		delete(c.cache, k)
	}
	c.mutex.Unlock()
	return
}

func GenKey(keys ...interface{}) interface{}{
	key := make([]interface{}, 2)
	for _, k := range keys {
		key = append(key, k)
	}
	return key
}