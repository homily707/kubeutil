package client

import "sync"

type store interface {
	Set(key string, value interface{})
	Get(key string) interface{}
	Contains(key string) bool
}

type concurrentHashMap struct {
	sync.RWMutex
	data map[string]interface{}
}

func NewConcurrentHashMap() concurrentHashMap {
	return concurrentHashMap{data: make(map[string]interface{})}
}

func (c concurrentHashMap) Set(key string, value interface{}) {
	c.Lock()
	defer c.Unlock()
	c.data[key] = value
}

func (c concurrentHashMap) Get(key string) interface{} {
	c.RLock()
	defer c.RUnlock()
	return c.data[key]
}

func (c concurrentHashMap) Contains(key string) bool {
	c.RLock()
	defer c.RUnlock()
	_, ok := c.data[key]
	return ok
}
