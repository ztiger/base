package common

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

var (
	ErrPoolExhausted = errors.New("Pool exhausted!")
	errPoolClosed    = errors.New("Pool closed!")
)

// 池结构体
type Pool struct {
	Dial             func() (interface{}, error)
	TestOnBorrow     func(interface{}, time.Time) error
	RemovePooledItem func(interface{}) error

	maxIdle     int
	maxActive   int
	closed      bool
	active      int
	idleTimeout time.Duration
	mu          sync.Mutex
	idle        list.List
}

type idleItem struct {
	item interface{}
	time time.Time
}

// 生成一个池
func NewPool(maxIdle int, maxActive int, idleTimeout time.Duration) *Pool {
	return &Pool{maxIdle: maxIdle, maxActive: maxActive, idleTimeout: idleTimeout}
}

// 从连接池获取一个redis client
func (p *Pool) Get() (interface{}, error) {
	for i := 0; i < 2; i++ {
		item, err := p.get()
		if err == nil {
			return item, nil
		}
		<-time.After(500) //定时作用
	}
	return nil, errors.New("Get pooled item timeout.")
}

// 从连接池获取一个redis client
func (p *Pool) get() (interface{}, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return nil, errPoolClosed
	}

	for i, n := 0, p.idle.Len(); i < n; i++ {
		item := p.idle.Front()
		if item == nil {
			break
		}
		p.idle.Remove(item)
		idleItm := item.Value.(idleItem)
		testFunc := p.TestOnBorrow
		if testFunc == nil || testFunc(idleItm.item, idleItm.time) == nil {
			return idleItm.item, nil
		} else {
			p.active -= 1
			p.RemovePooledItem(idleItm.item)
		}
	}

	if p.maxActive > 0 && p.active >= p.maxActive {
		return nil, ErrPoolExhausted
	}

	dialFunc := p.Dial
	cli, err := dialFunc()
	if err != nil {
		return nil, err
	}
	p.active += 1
	return cli, nil
}

//把连接放回连接池
func (p *Pool) Put(item interface{}) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	idleItm := idleItem{item: item, time: time.Now()}
	p.idle.PushFront(idleItm)
	if p.idle.Len() > p.maxIdle {
		item = p.idle.Remove(p.idle.Back()).(idleItem).item
	} else {
		item = nil
	}
	if item != nil {
		p.active -= 1
		p.RemovePooledItem(item)
	}
	return nil
}

//关闭连接池
func (p *Pool) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.closed = true
	idle := p.idle
	p.idle.Init()
	p.active = 0
	for itm := idle.Front(); itm != nil; itm = itm.Next() {
		p.RemovePooledItem(itm.Value.(idleItem).item)
	}
	return nil
}

// 获取连接池当前的大小
func (p *Pool) Size() int {
	return p.active
}
