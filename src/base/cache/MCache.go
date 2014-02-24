//@简单的内存缓存，用于缓存小量的，且读比较频繁的数据
//@v 0.0.1

package cache

import (
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

func (cache *MCache) Put(key string, val interface{}, expired int64) error {
	cache.lock.Lock()
	defer cache.lock.Unlock()

}
