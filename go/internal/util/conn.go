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

func (c *Conn) Close() error {
	return c.conn.Close()
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

func (c *Conn) WriteString(value string) bool {
	err := c.WriteAll([]byte(fmt.Sprintf("+%s\r\n", value)))

	if err != nil {
		c.HandleWriteError(err)
		return false
	}

	return true
}

func (c *Conn) WriteError(value string) bool {
	err := c.WriteAll([]byte(fmt.Sprintf("-%s\r\n", value)))

	if err != nil {
		c.HandleWriteError(err)
		return false
	}

	return true
}

func (c *Conn) WriteBulkString(value string) bool {
	err := c.WriteAll([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)))

	if err != nil {
		c.HandleWriteError(err)
		return false
	}

	return true
}

func (c *Conn) WriteInt(value int) bool {
	err := c.WriteAll([]byte(fmt.Sprintf(":%d\r\n", value)))

	if err != nil {
		c.HandleWriteError(err)
		return false
	}

	return true
}

func (c *Conn) WriteFloat32(value float32) bool {
	err := c.WriteAll([]byte(fmt.Sprintf(",%s\r\n", strconv.FormatFloat(float64(value), 'f', -1, 32))))

	if err != nil {
		c.HandleWriteError(err)
		return false
	}

	return true
}

func (c *Conn) WriteFloat64(value float64) bool {
	err := c.WriteAll([]byte(fmt.Sprintf(",%s\r\n", fmt.Sprint(value))))

	if err != nil {
		c.HandleWriteError(err)
		return false
	}

	return true
}

func (c *Conn) WriteInt64(value int64) bool {
	err := c.WriteAll([]byte(fmt.Sprintf(":%d\r\n", value)))

	if err != nil {
		c.HandleWriteError(err)
		return false
	}

	return true
}

func (c *Conn) WriteArray(value int) bool {
	err := c.WriteAll([]byte(fmt.Sprintf("*%d\r\n", value)))

	if err != nil {
		c.HandleWriteError(err)
		return false
	}

	return true
}

func (c *Conn) WriteMap(value int) bool {
	err := c.WriteAll([]byte(fmt.Sprintf("%%%d\r\n", value/2)))

	if err != nil {
		c.HandleWriteError(err)
		return false
	}

	return true
}

func (c *Conn) WriteSet(value int) bool {
	err := c.WriteAll([]byte(fmt.Sprintf("~%d\r\n", value)))

	if err != nil {
		c.HandleWriteError(err)
		return false
	}

	return true
}

func (c *Conn) WriteNull() bool {
	err := c.WriteAll([]byte("_\r\n"))

	if err != nil {
		c.HandleWriteError(err)
		return false
	}

	return true
}

func (c *Conn) WriteNullBulk() bool {
	err := c.WriteAll([]byte("$-1\r\n"))

	if err != nil {
		c.HandleWriteError(err)
		return false
	}

	return true
}

func (c *Conn) WriteNullArray() bool {
	err := c.WriteAll([]byte("*-1\r\n"))

	if err != nil {
		c.HandleWriteError(err)
		return false
	}

	return true
}

func (c *Conn) HandleWriteError(err error) {
	Logger.Printf("Failed to write to connection: '%s'\n", err)

	err = c.conn.Close()

	if err != nil {
		Logger.Printf("Unable to close connection: '%s'\n", err)
	}
}
