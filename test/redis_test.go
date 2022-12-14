package test

import (
	"fmt"
	"net"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-redis/redis"
	radish "github.com/hbina/radish"
	"github.com/stretchr/testify/assert"
)

var dbId int64 = 0
var port int = 6381

func CreateTestClient() *redis.Client {
	c := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("localhost:%d", port),
		DB:   int(atomic.AddInt64(&dbId, 1)),
	})
	return c
}

func init() {
	go radish.Run(port, false)
	time.Sleep(1 * time.Second)
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

	// Delete non-existent key returns 0
	{
		i, err := c.Del("abc").Result()
		assert.Equal(t, i, int64(0))
		assert.NoError(t, err)
	}

	// Create some keys then delete it
	{
		s, err := c.Set("k", "v", 0).Result()
		assert.Equal(t, "OK", s)
		assert.NoError(t, err)

		s, err = c.Set("k2", nil, 0).Result()
		assert.Equal(t, "OK", s)
		assert.NoError(t, err)

		_, err = c.Get("k2").Result()
		assert.NoError(t, err)

		i, err := c.Del("k", "k2").Result()
		assert.Equal(t, i, int64(2))
		assert.NoError(t, err)
	}

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

func TestRestoreCommand(t *testing.T) {
	c := CreateTestClient()

	{
		s, err := c.Set("dump-restore-string", "bar", time.Duration(0)).Result()
		assert.NoError(t, err)
		assert.Equal(t, "OK", s)

		dump, err := c.Dump("dump-restore-string").Result()
		assert.NoError(t, err)
		assert.NotEmpty(t, dump)

		s, err = c.Restore("dump-restore-string", time.Duration(0), dump).Result()
		assert.Equal(t, "BUSYKEY Target key name already exists.", err.Error())
		assert.Empty(t, s)

		i, err := c.Del("dump-restore-string").Result()
		assert.NoError(t, err)
		assert.Equal(t, int64(1), i)

		s, err = c.Restore("dump-restore-string", time.Duration(0), dump).Result()
		assert.NoError(t, err)
		assert.Equal(t, "OK", s)
	}

	{
		s1, err := c.ZAdd("dump-restore-zset", redis.Z{
			Score:  1,
			Member: "a",
		}, redis.Z{
			Score:  2,
			Member: "b",
		}, redis.Z{
			Score:  3,
			Member: "c",
		}, redis.Z{
			Score:  4,
			Member: "d",
		}).Result()
		assert.NoError(t, err)
		assert.Equal(t, int64(4), s1)

		dump, err := c.Dump("dump-restore-zset").Result()
		assert.NoError(t, err)
		assert.NotEmpty(t, dump)

		s2, err := c.Restore("dump-restore-zset", time.Duration(0), dump).Result()
		assert.Equal(t, "BUSYKEY Target key name already exists.", err.Error())
		assert.Empty(t, s2)

		i, err := c.Del("dump-restore-zset").Result()
		assert.NoError(t, err)
		assert.Equal(t, int64(1), i)

		s3, err := c.Restore("dump-restore-zset", time.Duration(0), dump).Result()
		assert.NoError(t, err)
		assert.Equal(t, "OK", s3)
	}
}

func TestBadRespCommand(t *testing.T) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", "localhost:6381")
	assert.NoError(t, err)

	tcpConn, err := net.DialTCP("tcp", nil, tcpAddr)
	assert.NoError(t, err)

	writeCount, err := tcpConn.Write([]byte("*-10\r\n"))
	assert.NoError(t, err)
	assert.Equal(t, 6, writeCount)

	c := CreateTestClient()

	s, err := c.Set("foo", "bar", 0).Result()
	assert.NoError(t, err)
	assert.Equal(t, "OK", s)

	s, err = c.Get("foo").Result()
	assert.NoError(t, err)
	assert.Equal(t, "bar", s)

	err = tcpConn.Close()
	assert.NoError(t, err)
}

func TestZaddCommad(t *testing.T) {

	// ZADD CH option changes return value to all changed elements
	{
		c := CreateTestClient()
		_, err := c.ZAdd("ztmp2", redis.Z{
			Score:  10,
			Member: "x",
		}, redis.Z{
			Score:  20,
			Member: "y",
		}, redis.Z{
			Score:  30,
			Member: "z",
		}).Result()
		assert.NoError(t, err)

		s, err := c.ZAdd("ztmp2", redis.Z{
			Score:  11,
			Member: "x",
		}, redis.Z{
			Score:  21,
			Member: "y",
		}, redis.Z{
			Score:  30,
			Member: "z",
		}).Result()
		assert.NoError(t, err)
		assert.Equal(t, int64(0), s)

		s, err = c.ZAddCh("ztmp2", redis.Z{
			Score:  12,
			Member: "x",
		}, redis.Z{
			Score:  22,
			Member: "y",
		}, redis.Z{
			Score:  30,
			Member: "z",
		}).Result()
		assert.NoError(t, err)
		assert.Equal(t, int64(2), s)
	}

	// ZRANGE basics
	{
		c := CreateTestClient()
		_, err := c.ZAdd("ztmp1", redis.Z{
			Score:  1,
			Member: "a",
		}, redis.Z{
			Score:  2,
			Member: "b",
		}, redis.Z{
			Score:  3,
			Member: "c",
		}, redis.Z{
			Score:  4,
			Member: "d",
		}).Result()
		assert.NoError(t, err)

		res, err := c.ZRange("ztmp1", 0, -1).Result()
		assert.NoError(t, err)
		assert.Equal(t, []string{"a", "b", "c", "d"}, res)

		res, err = c.ZRange("ztmp1", 0, -2).Result()
		assert.NoError(t, err)
		assert.Equal(t, []string{"a", "b", "c"}, res)

		res, err = c.ZRange("ztmp1", 1, -1).Result()
		assert.NoError(t, err)
		assert.Equal(t, []string{"b", "c", "d"}, res)

		res, err = c.ZRange("ztmp1", 1, -2).Result()
		assert.NoError(t, err)
		assert.Equal(t, []string{"b", "c"}, res)

		res, err = c.ZRange("ztmp1", -2, -1).Result()
		assert.NoError(t, err)
		assert.Equal(t, []string{"c", "d"}, res)

		res, err = c.ZRange("ztmp1", -2, -2).Result()
		assert.NoError(t, err)
		assert.Equal(t, []string{"c"}, res)

		// out of range start index
		res, err = c.ZRange("ztmp1", -5, 2).Result()
		assert.NoError(t, err)
		assert.Equal(t, []string{"a", "b", "c"}, res)

		res, err = c.ZRange("ztmp1", -5, 1).Result()
		assert.NoError(t, err)
		assert.Equal(t, []string{"a", "b"}, res)

		res, err = c.ZRange("ztmp1", 5, -1).Result()
		assert.NoError(t, err)
		assert.Equal(t, []string{}, res)

		res, err = c.ZRange("ztmp1", 5, -2).Result()
		assert.NoError(t, err)
		assert.Equal(t, []string{}, res)
	}
}

func TestZremrangebyScoreCommand(t *testing.T) {

	{
		c := CreateTestClient()

		_, err := c.ZAdd("zset", redis.Z{
			Score:  1,
			Member: "a",
		}, redis.Z{
			Score:  2,
			Member: "b",
		}, redis.Z{
			Score:  3,
			Member: "c",
		}, redis.Z{
			Score:  4,
			Member: "d",
		}, redis.Z{
			Score:  5,
			Member: "e",
		}).Result()
		assert.NoError(t, err)

		_, err = c.ZRemRangeByScore("zset", "2", "4").Result()
		assert.NoError(t, err)
	}
}

func TestZrankCommand(t *testing.T) {
	{
		c := CreateTestClient()

		_, err := c.ZAdd("zset", redis.Z{
			Score:  10,
			Member: "x",
		}, redis.Z{
			Score:  20,
			Member: "y",
		}, redis.Z{
			Score:  30,
			Member: "z",
		}).Result()
		assert.NoError(t, err)

		r1, err := c.ZRank("zset", "x").Result()
		assert.NoError(t, err)
		assert.Equal(t, int64(0), r1)
		r1, err = c.ZRank("zset", "y").Result()
		assert.NoError(t, err)
		assert.Equal(t, int64(1), r1)
		r1, err = c.ZRank("zset", "z").Result()
		assert.NoError(t, err)
		assert.Equal(t, int64(2), r1)

		r1, err = c.ZRevRank("zset", "x").Result()
		assert.NoError(t, err)
		assert.Equal(t, int64(2), r1)
		r1, err = c.ZRevRank("zset", "y").Result()
		assert.NoError(t, err)
		assert.Equal(t, int64(1), r1)
		r1, err = c.ZRevRank("zset", "z").Result()
		assert.NoError(t, err)
		assert.Equal(t, int64(0), r1)

		r2, err := c.ZRank("zset", "foo").Result()
		assert.Error(t, err)
		assert.Equal(t, int64(0), r2)
		r2, err = c.ZRevRank("zset", "foo").Result()
		assert.Error(t, err)
		assert.Equal(t, int64(0), r2)
	}
}
