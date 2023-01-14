package test

import (
	"testing"

	"github.com/hbina/radish/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestTakeBytesUntilClrf(t *testing.T) {
	{
		in := []byte("+OK\r\n")
		str, leftover, ok := util.TakeBytesUntilClrf(in)
		assert.True(t, ok)
		assert.Equal(t, []byte("+OK"), str)
		assert.Empty(t, leftover)
	}
	{
		in := []byte("-Error message\r\n")
		str, leftover, ok := util.TakeBytesUntilClrf(in)
		assert.True(t, ok)
		assert.Equal(t, []byte("-Error message"), str)
		assert.Empty(t, leftover)
	}
	{
		in := []byte("\r\n")
		str, leftover, ok := util.TakeBytesUntilClrf(in)
		assert.True(t, ok)
		assert.Equal(t, []byte(""), str)
		assert.Empty(t, leftover)
	}
	{
		in := []byte("\n")
		str, leftover, ok := util.TakeBytesUntilClrf(in)
		assert.False(t, ok)
		assert.Equal(t, []byte("\n"), str)
		assert.Empty(t, leftover)
	}
	{
		in := []byte("\r")
		str, leftover, ok := util.TakeBytesUntilClrf(in)
		assert.False(t, ok)
		assert.Equal(t, []byte("\r"), str)
		assert.Empty(t, leftover)
	}
	{
		in := []byte("-Error message\n")
		str, leftover, ok := util.TakeBytesUntilClrf(in)
		assert.False(t, ok)
		assert.Equal(t, []byte("-Error message\n"), str)
		assert.Empty(t, leftover)
	}
}

func TestStringifyGoodRespString(t *testing.T) {
	res, ok, leftover := util.StringifyRespBytes([]byte("+OK\r\n"))
	assert.Equal(t, "OK", res)
	assert.True(t, ok)
	assert.Empty(t, leftover)
}

func TestStringifyBadRespString(t *testing.T) {
	res, ok, leftover := util.StringifyRespBytes([]byte("+OK\n"))
	assert.Empty(t, res)
	assert.False(t, ok)
	assert.Empty(t, leftover)
}

func TestStringifyGoodRespError(t *testing.T) {
	res, ok, leftover := util.StringifyRespBytes([]byte("-Error message\r\n"))
	assert.Equal(t, "Error message", res)
	assert.True(t, ok)
	assert.Empty(t, leftover)
}

func TestStringifyBadRespError(t *testing.T) {
	res, ok, leftover := util.StringifyRespBytes([]byte("-Error message\r"))
	assert.Empty(t, res)
	assert.False(t, ok)
	assert.Empty(t, leftover)
}

func TestStringifyGoodRespInteger(t *testing.T) {
	// Integers
	res, ok, leftover := util.StringifyRespBytes([]byte(":1\r\n"))
	assert.Equal(t, "(integer) 1", res)
	assert.True(t, ok)
	assert.Empty(t, leftover)
	res, ok, leftover = util.StringifyRespBytes([]byte(":-100\r\n"))
	assert.Equal(t, "(integer) -100", res)
	assert.True(t, ok)
	assert.Empty(t, leftover)
	res, ok, leftover = util.StringifyRespBytes([]byte(":32\r\n"))
	assert.Equal(t, "(integer) 32", res)
	assert.True(t, ok)
	assert.Empty(t, leftover)
	res, ok, leftover = util.StringifyRespBytes([]byte(":-0\r\n"))
	assert.Equal(t, "(integer) 0", res)
	assert.True(t, ok)
	assert.Empty(t, leftover)
}

func TestStringifyBadRespInteger(t *testing.T) {
	// Integers
	res, ok, leftover := util.StringifyRespBytes([]byte(":1ddd\r\n"))
	assert.Empty(t, res)
	assert.False(t, ok)
	assert.Empty(t, leftover)

	res, ok, leftover = util.StringifyRespBytes([]byte(":-1d00\r\n"))
	assert.Empty(t, res)
	assert.False(t, ok)
	assert.Empty(t, leftover)

	res, ok, leftover = util.StringifyRespBytes([]byte(":v32\r\n"))
	assert.Empty(t, res)
	assert.False(t, ok)
	assert.Empty(t, leftover)

	res, ok, leftover = util.StringifyRespBytes([]byte(":-v0\r\n"))
	assert.Empty(t, res)
	assert.False(t, ok)
	assert.Empty(t, leftover)
}

func TestStringifyRespGoodBulkString(t *testing.T) {
	res, ok, leftover := util.StringifyRespBytes([]byte("$5\r\nhello\r\n"))
	assert.Equal(t, "\"hello\"", res)
	assert.True(t, ok)
	assert.Empty(t, leftover)
	res, ok, leftover = util.StringifyRespBytes([]byte("$0\r\n\r\n"))
	assert.Equal(t, "\"\"", res)
	assert.True(t, ok)
	assert.Empty(t, leftover)
	res, ok, leftover = util.StringifyRespBytes([]byte("$-1\r\n"))
	assert.Equal(t, "(nil)", res)
	assert.True(t, ok)
	assert.Empty(t, leftover)
}

func TestStringifyRespBadBulkString(t *testing.T) {
	res, ok, leftover := util.StringifyRespBytes([]byte("$5\r\nlo\r\n"))
	assert.Empty(t, res)
	assert.False(t, ok)
	assert.Empty(t, leftover)

	res, ok, leftover = util.StringifyRespBytes([]byte("$0\r\nddd\r\n"))
	assert.Empty(t, res)
	assert.False(t, ok)
	assert.Empty(t, leftover)

	res, ok, leftover = util.StringifyRespBytes([]byte("$2\r\nddd\r\n"))
	assert.Empty(t, res)
	assert.False(t, ok)
	assert.Empty(t, leftover)
}

func TestStringifyGoodRespArray(t *testing.T) {
	res, ok, leftover := util.StringifyRespBytes([]byte("*0\r\n"))
	assert.Equal(t, "(empty)", res)
	assert.True(t, ok)
	assert.Empty(t, leftover)

	res, ok, leftover = util.StringifyRespBytes([]byte("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"))
	assert.Equal(t, "1) \"hello\"\n2) \"world\"", res)
	assert.True(t, ok)
	assert.Empty(t, leftover)

	res, ok, leftover = util.StringifyRespBytes([]byte("*2\r\n+hello\r\n*2\r\n+world\r\n+good\r\n"))
	assert.Equal(t, "1) \"hello\"\n2) 1) \"world\"\n   2) \"good\"", res)
	assert.True(t, ok)
	assert.Empty(t, leftover)

	res, ok, leftover = util.StringifyRespBytes([]byte("*3\r\n+a\r\n*2\r\n*2\r\n+b\r\n+c\r\n+d\r\n+e\r\n"))
	assert.Equal(t, "1) \"a\"\n2) 1) 1) \"b\"\n      2) \"c\"\n   2) \"d\"\n3) \"e\"", res)
	assert.True(t, ok)
	assert.Empty(t, leftover)

	res, ok, leftover = util.StringifyRespBytes([]byte("*4\r\n$7\r\nmatches\r\n*4\r\n*2\r\n*2\r\n:3\r\n:4\r\n*2\r\n:12\r\n:13\r\n*2\r\n*2\r\n:2\r\n:2\r\n*2\r\n:8\r\n:8\r\n*2\r\n*2\r\n:1\r\n:1\r\n*2\r\n:4\r\n:4\r\n*2\r\n*2\r\n:0\r\n:0\r\n*2\r\n:0\r\n:0\r\n$3\r\nlen\r\n:5\r\n"))
	assert.Equal(t, "1) \"matches\"\n2) 1) 1) 1) (integer) 3\n         2) (integer) 4\n      2) 1) (integer) 12\n         2) (integer) 13\n   2) 1) 1) (integer) 2\n         2) (integer) 2\n      2) 1) (integer) 8\n         2) (integer) 8\n   3) 1) 1) (integer) 1\n         2) (integer) 1\n      2) 1) (integer) 4\n         2) (integer) 4\n   4) 1) 1) (integer) 0\n         2) (integer) 0\n      2) 1) (integer) 0\n         2) (integer) 0\n3) \"len\"\n4) (integer) 5", res)
	assert.True(t, ok)
	assert.Empty(t, leftover)
}

func TestStringifyBadRespArray(t *testing.T) {
	res, ok, leftover := util.StringifyRespBytes([]byte("*3\r\n"))
	assert.Empty(t, res)
	assert.False(t, ok)
	assert.Empty(t, leftover)

	res, ok, leftover = util.StringifyRespBytes([]byte("*2\r\n$5\r\nhlo\r\n$5\r\nworld\r\n"))
	assert.Empty(t, res)
	assert.False(t, ok)
	assert.Empty(t, leftover)
}
