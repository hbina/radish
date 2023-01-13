package util

import (
	"fmt"
	"net"
	"strconv"
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

func (c *Conn) WriteAll(in []byte) error {
	t := 0

	for t != len(in) {
		c, err := c.conn.Write(in[t:])

		if err != nil {
			return err
		}

		t += c
	}

	return nil
}

func (c *Conn) Close() error {
	return c.conn.Close()
}

func (c *Conn) WriteString(value string) error {
	return c.WriteAll([]byte(fmt.Sprintf("+%s\r\n", value)))
}

func (c *Conn) WriteError(value string) error {
	return c.WriteAll([]byte(fmt.Sprintf("-%s\r\n", value)))
}

func (c *Conn) WriteBulkString(value string) error {
	return c.WriteAll([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)))
}

func (c *Conn) WriteInt(value int) error {
	return c.WriteAll([]byte(fmt.Sprintf(":%d\r\n", value)))
}

func (c *Conn) WriteFloat32(value float32) error {
	return c.WriteAll([]byte(fmt.Sprintf(",%s\r\n", strconv.FormatFloat(float64(value), 'f', -1, 32))))
}

func (c *Conn) WriteFloat64(value float64) error {
	return c.WriteAll([]byte(fmt.Sprintf(",%s\r\n", fmt.Sprint(value))))
}

func (c *Conn) WriteInt64(value int64) error {
	return c.WriteAll([]byte(fmt.Sprintf(":%d\r\n", value)))
}

func (c *Conn) WriteArray(value int) error {
	return c.WriteAll([]byte(fmt.Sprintf("*%d\r\n", value)))
}

func (c *Conn) WriteMap(value int) error {
	return c.WriteAll([]byte(fmt.Sprintf("%%%d\r\n", value/2)))
}

func (c *Conn) WriteSet(value int) error {
	return c.WriteAll([]byte(fmt.Sprintf("~%d\r\n", value)))
}

func (c *Conn) WriteNull() error {
	return c.WriteAll([]byte("_\r\n"))
}

func (c *Conn) WriteNullBulk() error {
	return c.WriteAll([]byte("$-1\r\n"))
}

func (c *Conn) WriteNullArray() error {
	return c.WriteAll([]byte("*-1\r\n"))
}
