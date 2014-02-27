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
func NewMCache(every int, dur time.Duration) *MCache {
	cache := &MCache{dur: dur, Every: every, items: make(map[string]*MemoryItem)}
	go cache.checkExpiration()
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
	cache.items[key] = &item
	return nil
}

//删除key
func (cache *MCache) Delete(key string) error {
	cache.lock.Lock()
	defer cache.lock.Unlock()
	if _, ok := cache.items[key]; !ok {
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

//检查健值过期协程
func (cache *MCache) checkExpiration() {
	if cache.Every < 1 {
		return
	}

	for {
		<-time.After(cache.dur)
		if cache.items == nil {
			return
		}
		for key, _ := range cache.items {
			cache.item_expired(key)
		}
	}
}

func (cache *MCache) item_expired(key string) bool {
	cache.lock.Lock()
	defer cache.lock.Unlock()
	item, ok := cache.items[key]
	if !ok {
		return true
	}
	sec := time.Now().Unix() - item.Lastaccess.Unix()
	if sec >= item.expired {
		delete(cache.items, key)
		return true
	}
	return false
}
