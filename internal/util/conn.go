package util

import (
	"fmt"
	"net"
)

type Conn struct {
	conn net.Conn
}

func NewConn(conn net.Conn) *Conn {
	return &Conn{
		conn: conn,
	}
}

func (c *Conn) Read(buffer []byte) (int, error) {
	return c.conn.Read(buffer)
}

func (c *Conn) WriteResp(resp Resp2) error {
	return resp.WriteToConn(c.conn)
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

func (c *Conn) WriteError(value string) {
	c.conn.Write([]byte(fmt.Sprintf("-%s\r\n", value)))
}

func (c *Conn) WriteRaw(value []byte) {
	c.conn.Write(value)
}
