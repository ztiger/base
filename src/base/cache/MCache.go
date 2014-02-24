//@简单的内存缓存，用于缓存小量的，且读比较频繁的数据
//@v 0.0.1

package cache

import (
	"errors"
	"sync"
	"time"
)

var (
	DefaultEvery int = 60 //clock time of recycling the expired cache items in memory
)

type MemoryItem struct {
	val        interface{}
	Lastaccess time.Time
	expired    int64
}

type MCache struct {
	lock  sync.RWMutex
	dur   time.Duration
	items map[string]*MemoryItem
	Every int
}

//生成一个内存缓存对象
func NewMCache() *MCache {
	cache := &MCache{items: make(map[string]*MemoryItem)}
	return cache
}

//获取值
func (cache *MCache) Get(key string) interface{} {
	cache.lock.RLock()
	defer cache.lock.RUnlock()
	item, ok := cache.items[key]
	if !ok {
		return nil
	}
	if (time.Now().Unix() - item.Lastaccess.Unix()) > item.expired {
		return nil
	}
	return item.val
}

//存数据
func (cache *MCache) Put(key string, val interface{}, expired int64) error {
	cache.lock.Lock()
	defer cache.lock.Unlock()
	item := MemoryItem{
		val:        val,
		Lastaccess: time.Now(),
		expired:    expired,
	}
	cache.items[key] = &t
	return nil
}

//删除key
func (cache *MCache) Delete(key string) error {
	cache.lock.Lock()
	defer cache.lock.Unlock()
	if _, ok := cache[key]; !ok {
		return errors.New("Key not exist.")
	}
	delete(cache.items, key)
	_, valid := cache.items[key]
	if valid {
		return errors.New("Delete key error")
	}
	return nil
}

//判断key 是否存在
func (cache *MCache) IsExist(key string) bool {
	cache.lock.RLock()
	defer cache.lock.RUnlock()
	_, ok := cache.items[key]
	return ok
}

//清除缓存
func (cache *MCache) ClearAll() error {
	cache.lock.Lock()
	defer cache.lock.Unlock()
	cache.items = make(map[string]*MemoryItem)
	return nil
}
