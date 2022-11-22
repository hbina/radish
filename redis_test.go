package redis

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

var dbId int64 = 0
var server *Redis = Default()
var port string = fmt.Sprintf("localhost:%s", "6380")

func CreateTestClient() *redis.Client {
	c := redis.NewClient(&redis.Options{
		Addr: port,
		DB:   int(atomic.AddInt64(&dbId, 1)),
	})
	return c
}

func init() {
	go server.Run(":6380")
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
	assert.Equal(t, i, int64(0))
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
		assert.Equal(t, int64(1), s)
	}

	{
		s, err := c.ZAdd("myzset", redis.Z{
			Score: 1, Member: "uno",
		}).Result()
		assert.NoError(t, err)
		assert.Equal(t, int64(1), s)
	}

	{
		s, err := c.ZAdd("myzset",
			redis.Z{
				Score: 2, Member: "two",
			}, redis.Z{
				Score: 3, Member: "three",
			}).Result()
		assert.NoError(t, err)
		assert.Equal(t, int64(2), s)
	}
}

// Tests taken from "SRANDMEMBER with <count> - $type"
// in unit/types/set.tcl from redis
func TestSrandMember(t *testing.T) {
	c := CreateTestClient()

	for i := 0; i < 50; i++ {
		s, err := c.SAdd("testsrandmember", fmt.Sprint(i)).Result()
		assert.NoError(t, err)
		assert.True(t, s > 0)
	}

	var sizes = []int{5, 45}

	for _, size := range sizes {
		_, err := c.Del("testrandmember2").Result()
		assert.NoError(t, err)

		// Iterate many times to increase probability of succeeding
		for i := 0; i < 1000; i++ {
			s, err := c.SRandMemberN("testsrandmember", int64(size)).Result()
			assert.NoError(t, err)

			for _, v := range s {
				_, err := c.SAdd("testsrandmember2", v).Result()
				assert.NoError(t, err)
			}

			c1, err := c.SCard("testsrandmember").Result()
			assert.NoError(t, err)
			c2, err := c.SCard("testsrandmember2").Result()
			assert.NoError(t, err)
			if c1 == c2 {
				break
			}
		}

		c1, err := c.SCard("testsrandmember").Result()
		assert.NoError(t, err)
		c2, err := c.SCard("testsrandmember2").Result()
		assert.NoError(t, err)
		assert.Equal(t, c1, c2, fmt.Sprintf("Failed when size = %d", size))
	}

}
