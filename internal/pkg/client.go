package pkg

import "github.com/hbina/radish/internal/util"

// A connected Client.
type Client struct {
	conn  *util.Conn
	dbId  uint64
	redis *Redis
}

func (c *Client) Read(buffer []byte) (int, error) {
	return c.conn.Read(buffer)
}

func (c *Client) Redis() *Redis {
	return c.redis
}

func (c *Client) Conn() *util.Conn {
	return c.conn
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

func (c *Client) UseResp2() {
	c.conn.UseResp2()
}

func (c *Client) UseResp3() {
	c.conn.UseResp3()
}
