package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTakeBytesUntilClrf(t *testing.T) {
	{
		ok := []byte("+OK\r\n")
		str, leftover := TakeBytesUntilClrf(ok)
		assert.Equal(t, str, []byte("+OK"))
		assert.Empty(t, leftover)
	}
	{
		ok := []byte("-Error message\r\n")
		str, leftover := TakeBytesUntilClrf(ok)
		assert.Equal(t, str, []byte("-Error message"))
		assert.Empty(t, leftover)
	}
}

func TestDisplayRespReply(t *testing.T) {
	// Simple strings
	res, leftover := CreateRespReply([]byte("+OK\r\n"))
	assert.Equal(t, "OK", res)
	assert.Empty(t, leftover)

	// Errors

	res, leftover = CreateRespReply([]byte("-Error message\r\n"))
	assert.Equal(t, "Error message", res)
	assert.Empty(t, leftover)

	// Integers
	res, leftover = CreateRespReply([]byte(":1\r\n"))
	assert.Equal(t, "1", res)
	assert.Empty(t, leftover)
	res, leftover = CreateRespReply([]byte(":-100\r\n"))
	assert.Equal(t, "-100", res)
	assert.Empty(t, leftover)
	res, leftover = CreateRespReply([]byte(":32\r\n"))
	assert.Equal(t, "32", res)
	assert.Empty(t, leftover)
	res, leftover = CreateRespReply([]byte(":-0\r\n"))
	assert.Equal(t, "0", res)
	assert.Empty(t, leftover)

	// Bulk strings
	res, leftover = CreateRespReply([]byte("$5\r\nhello\r\n"))
	assert.Equal(t, "hello", res)
	assert.Empty(t, leftover)
	res, leftover = CreateRespReply([]byte("$0\r\n\r\n"))
	assert.Equal(t, "", res)
	assert.Empty(t, leftover)
	res, leftover = CreateRespReply([]byte("$-1\r\n"))
	assert.Equal(t, "(nil)", res)
	assert.Empty(t, leftover)

}

func TestCreateRespReplyFromRespArray(t *testing.T) {
	// Arrays
	_, leftover := CreateRespReply([]byte("*2\r\n$7\r\nmatches\r\n*2\r\n+hello\r\n+world\r\n"))
	assert.Empty(t, leftover)

	_, leftover = CreateRespReply([]byte("*0\r\n"))
	assert.Empty(t, leftover)

	_, leftover = CreateRespReply([]byte("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"))
	assert.Empty(t, leftover)

	_, leftover = CreateRespReply([]byte("*3\r\n$7\r\nmatches\r\n*2\r\n+hello\r\n+world\r\n$6\r\nwhat??\r\n"))
	assert.Empty(t, leftover)
}
