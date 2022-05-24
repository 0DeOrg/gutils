package raftutils

/**
 * @Author: lee
 * @Description:
 * @File: cache
 * @Date: 2022-05-23 2:21 下午
 */

import (
	"encoding/json"
	"io"
	"sync"
)

type logCache struct {
	data map[string]string
	sync.RWMutex
}

func newLogCache() *logCache {
	ret := &logCache{}
	ret.data = make(map[string]string, 16)
	return ret
}

func (c *logCache) Get(key string) string {
	c.RLock()
	ret := c.data[key]
	c.RUnlock()
	return ret
}

func (c *logCache) Set(key string, value string) {
	c.Lock()
	defer c.Unlock()
	c.data[key] = value
}

// Marshal serializes cache data
func (c *logCache) Marshal() ([]byte, error) {
	c.RLock()
	defer c.RUnlock()
	dataBytes, err := json.Marshal(c.data)
	return dataBytes, err
}

// UnMarshal deserializes cache data
func (c *logCache) UnMarshal(serialized io.ReadCloser) error {
	var newData map[string]string
	if err := json.NewDecoder(serialized).Decode(&newData); err != nil {
		return err
	}

	c.Lock()
	defer c.Unlock()
	c.data = newData

	return nil
}
