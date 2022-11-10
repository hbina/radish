package redis

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

var dbId int64 = 10

func CreateTestClient() *redis.Client {
	c := redis.NewClient(&redis.Options{
		Addr: "localhost:6380",
		DB:   int(atomic.AddInt64(&dbId, 1)),
	})
	return c
}

func TestPingCommand(t *testing.T) {
	c := CreateTestClient()

	s, err := c.Ping().Result()
	assert.Equal(t, "PONG", s)
	assert.NoError(t, err)

	pingCmd := redis.NewStringCmd("ping", "Hello, redis server!")
	c.Process(pingCmd)
	s, err = pingCmd.Result()
	assert.Equal(t, "Hello, redis server!", s)
	assert.NoError(t, err)
}

func TestSetGetCommand(t *testing.T) {
	c := CreateTestClient()

	s, err := c.Set("k", "v", 0).Result()
	assert.Equal(t, "OK", s)
	assert.NoError(t, err)

	s, err = c.Set("k2", nil, 0).Result()
	assert.Equal(t, "OK", s)
	assert.NoError(t, err)

	s, err = c.Set("k3", "v", 1*time.Hour).Result()
	assert.Equal(t, "OK", s)
	assert.NoError(t, err)

	s, err = c.Get("k").Result()
	assert.Equal(t, "v", s)
	assert.NoError(t, err)
}

func TestDelCommand(t *testing.T) {
	c := CreateTestClient()

	i, err := c.Del("k", "k3").Result()
	assert.Equal(t, i, int64(2))
	assert.NoError(t, err)

	i, err = c.Del("abc").Result()
	assert.Equal(t, i, int64(1))
	assert.NoError(t, err)
}

func TestTtlCommand(t *testing.T) {
	c := CreateTestClient()

	s, err := c.Set("aKey", "hey", 1*time.Minute).Result()
	assert.Equal(t, "OK", s)
	assert.NoError(t, err)
	s, err = c.Set("bKey", "hallo", 0).Result()
	assert.Equal(t, "OK", s)
	assert.NoError(t, err)

	ttl, err := c.TTL("aKey").Result()
	assert.True(t, ttl.Seconds() > 55 && ttl.Seconds() < 61, "ttl: %d", ttl)
	assert.NoError(t, err)

	ttl, err = c.TTL("none").Result()
	assert.Equal(t, time.Duration(-2000000000), ttl)
	assert.NoError(t, err)

	ttl, err = c.TTL("bKey").Result()
	assert.NoError(t, err)
	assert.Equal(t, time.Duration(-1000000000), ttl)
}

func TestExpiry(t *testing.T) {
	c := CreateTestClient()

	s, err := c.Set("x", "val", 10*time.Millisecond).Result()
	assert.NoError(t, err)
	assert.Equal(t, "OK", s)

	time.Sleep(20 * time.Millisecond)

	_, err = c.Get("x").Result()
	assert.Error(t, err)
	assert.Equal(t, err, redis.Nil)
}

func TestZaddCommand(t *testing.T) {
	c := CreateTestClient()

	{
		s, err := c.ZAdd("myzset", redis.Z{
			Score: 1, Member: "one",
		}).Result()
		assert.NoError(t, err)
		assert.Equal(t, "(integer) 1", s)
	}

	{
		s, err := c.ZAdd("myzset", redis.Z{
			Score: 1, Member: "uno",
		}).Result()
		assert.NoError(t, err)
		assert.Equal(t, "(integer) 1", s)
	}

	{
		s, err := c.ZAdd("myzset",
			redis.Z{
				Score: 2, Member: "two",
			}, redis.Z{
				Score: 3, Member: "three",
			}).Result()
		assert.NoError(t, err)
		assert.Equal(t, "(integer) 2", s)
	}
}
