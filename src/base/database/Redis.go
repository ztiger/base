package database

import (
	"base/common"
	"errors"
	"menteslibres.net/gosexy/redis"
	"time"
)

var (
	ErrPoolExhausted = errors.New("Connection pool exhausted!")
	errPoolClosed    = errors.New("Connection pool closed!")
)

//redis 连接池
type RedisPool struct {
	host string
	port uint
	pool *common.Pool
}

//生成一个redis pool
func NewRedisPool(host string, port uint, maxIdle int, maxActive int, idleTimeout time.Duration) *RedisPool {
	pool := common.NewPool(maxIdle, maxActive, idleTimeout)
	redisPool := &RedisPool{host: host, port: port, pool: pool}

	redisPool.pool.Dial = redisPool.dial
	redisPool.pool.TestOnBorrow = redisPool.testOnBorrow
	redisPool.pool.RemovePooledItem = redisPool.removePooledItem
	return redisPool
}

//生成一个redis连接
func (p *RedisPool) dial() (interface{}, error) {
	client := redis.New()
	err := client.Connect(p.host, p.port)
	if err != nil {
		return nil, errors.New("Could not connect to redis server.")
	}
	return client, nil
}

//测试连接是否正常
func (p *RedisPool) testOnBorrow(item interface{}, time time.Time) error {
	return nil
}

//关闭连接
func (p *RedisPool) removePooledItem(item interface{}) error {
	var redisClient *redis.Client
	switch v := item.(type) {
	case *redis.Client:
		redisClient = v
	default:
		return errors.New("")
	}
	redisClient.Quit()
	return nil
}

// 从连接池获取一个redis client
func (p *RedisPool) Get() (*redis.Client, error) {
	pooledItem, err := p.pool.Get()
	if err != nil {
		return nil, errors.New("Get pooled redis client timeout.")
	}

	var redisClient *redis.Client
	switch v := pooledItem.(type) {
	case *redis.Client:
		redisClient = v
	default:
		return nil, errors.New("")
	}
	return redisClient, nil
}

//把连接放回连接池
func (p *RedisPool) Put(client *redis.Client) error {
	p.pool.Put(client)
	return nil
}

//关闭连接池
func (p *RedisPool) Close() error {
	p.pool.Close()
	return nil
}

// 获取连接池当前的大小
func (p *RedisPool) Size() int {
	return p.pool.Size()
}
