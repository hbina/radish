package redis

import (
	"github.com/tidwall/redcon"
)

// A connected Client.
type Client struct {
	// The client connection.
	conn redcon.Conn

	// Selected database (default 0)
	dbId uint64

	redis *Redis
}

// Redis gets the redis instance.
func (c *Client) Redis() *Redis {
	return c.redis
}

func (c *Client) Conn() redcon.Conn {
	return c.conn
}

// SelectDb selects the clients database.
func (c *Client) SelectDb(dbId uint64) {
	c.dbId = dbId
}

// DbId gets the clients selected database id.
func (c *Client) DbId() uint64 {
	return c.dbId
}

// Db gets the clients selected database.
func (c *Client) Db() *Db {
	return c.Redis().RedisDb(c.DbId())
}
