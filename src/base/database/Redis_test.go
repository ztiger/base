package database

import (
	//"menteslibres.net/gosexy/redis"
	"testing"
)

func TestRedisPool(t *testing.T) {
	redisPool := NewRedisPool("127.0.0.1", 6379, 10, 10, 0)
	if redisPool == nil {
		t.Fatal("...")
	}

	client, err := redisPool.Get()
	if err != nil {
		t.Fatal(err)
	}

	if _, err := client.Set("test", "aaaa"); err != nil {
		t.Fatal(err)
	}

	if redisPool.Size() != 1 {
		t.Fatal(redisPool.Size())
	}

	redisPool.Close()
	if redisPool.Size() != 0 {
		t.Fatal(redisPool.Size())
	}

}
