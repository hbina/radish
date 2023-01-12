package pkg

import (
	"fmt"
	"net"
)

// A connected Client.
type Client struct {
	conn  net.Conn
	dbId  uint64
	redis *Redis
	r3    bool
}

func (c *Client) Read(buffer []byte) (int, error) {
	return c.conn.Read(buffer)
}

func (c *Client) Redis() *Redis {
	return c.redis
}

func (c *Client) DbId() uint64 {
	return c.dbId
}

// SetDb sets the database that this client interacts with.
func (c *Client) SetDb(dbId uint64) {
	c.dbId = dbId
}

// Db gets the clients selected database.
func (c *Client) Db() *Db {
	return c.redis.GetDb(c.dbId)
}

func (c *Client) WriteError(value string) {
	c.conn.Write([]byte(fmt.Sprintf("-%s\r\n", value)))
}

func (c *Client) WriteArray(value int) {
	c.conn.Write([]byte(fmt.Sprintf("*%d\r\n", value)))
}

func (c *Client) WriteMap(value int) {
	if c.r3 {
		// To escape % we need %%
		c.conn.Write([]byte(fmt.Sprintf("%%%d\r\n", value)))
	} else {
		c.conn.Write([]byte(fmt.Sprintf("*%d\r\n", value*2)))
	}
}

func (c *Client) WriteSet(value int) {
	if c.r3 {
		c.conn.Write([]byte(fmt.Sprintf("~%d\r\n", value)))
	} else {
		c.conn.Write([]byte(fmt.Sprintf("*%d\r\n", value)))
	}
}

func (c *Client) WriteSimpleString(value string) {
	c.conn.Write([]byte(fmt.Sprintf("+%s\r\n", value)))
}

func (c *Client) WriteBulkString(value string) {
	c.conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)))
}

func (c *Client) WriteInt(value int) {
	c.conn.Write([]byte(fmt.Sprintf(":%d\r\n", value)))
}

func (c *Client) WriteDouble(value float64) {
	if c.r3 {
		c.conn.Write([]byte(fmt.Sprintf(",%f\r\n", value)))
	} else {
		c.conn.Write([]byte(fmt.Sprintf("+%f\r\n", value)))
	}
}

func (c *Client) WriteInt64(value int64) {
	c.conn.Write([]byte(fmt.Sprintf(":%d\r\n", value)))
}

func (c *Client) WriteNullBulk() {
	if c.r3 {
		c.conn.Write([]byte("_\r\n"))
	} else {
		c.conn.Write([]byte("$-1\r\n"))
	}
}

func (c *Client) WriteNullArray() {
	if c.r3 {
		c.conn.Write([]byte("_\r\n"))
	} else {
		c.conn.Write([]byte("*-1\r\n"))
	}
}

func (c *Client) WriteNull() {
	c.conn.Write([]byte("_\r\n"))
}

func (c *Client) WriteRaw(value []byte) {
	c.conn.Write(value)
}

func (c *Client) UseResp2() {
	c.r3 = false
}

func (c *Client) UseResp3() {
	c.r3 = true
}
