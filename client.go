package redis

import (
	"fmt"
	"net"
)

type Conn struct {
	conn net.Conn
}

func (c *Conn) WriteResp(resp *RespArray) {}

func (c *Conn) WriteError(value string) {
	c.conn.Write([]byte(fmt.Sprintf("-%s\r\n", value)))
}

func (c *Conn) WriteArray(value int) {
	c.conn.Write([]byte(fmt.Sprintf("*%d\r\n", value)))
}

func (c *Conn) WriteString(value string) {
	c.conn.Write([]byte(fmt.Sprintf("+%s\r\n", value)))
}

func (c *Conn) WriteBulkString(value string) {
	c.conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)))
}

func (c *Conn) WriteInt(value int) {
	c.conn.Write([]byte(fmt.Sprintf(":%d\r\n", value)))
}

func (c *Conn) WriteInt64(value int64) {
	c.conn.Write([]byte(fmt.Sprintf(":%d\r\n", value)))
}

func (c *Conn) WriteNull() {
	c.conn.Write([]byte("$-1\r\n"))
}

func (c *Conn) WriteNullArray() {
	c.conn.Write([]byte("*-1\r\n"))
}

func (c *Conn) WriteRaw(value []byte) {
	c.conn.Write(value)
}

// A connected Client.
type Client struct {
	conn  *Conn
	dbId  uint64
	redis *Redis
}

func (c *Client) Read(buffer []byte) (int, error) {
	return c.conn.conn.Read(buffer)
}

func (c *Client) Redis() *Redis {
	return c.redis
}

func (c *Client) Conn() *Conn {
	return c.conn
}

// SetDb sets the database that this client interacts with.
func (c *Client) SetDb(dbId uint64) {
	c.dbId = dbId
}

// Db gets the clients selected database.
func (c *Client) Db() *Db {
	return c.redis.GetDb(c.dbId)
}
