package redis

import (
	"github.com/tidwall/redcon"
)

// TODO Client flags
const (
// client is a master
// client is a slave
// ...
)

// A connected Client.
type Client struct {
	clientId ClientId
	// The client connection.
	conn redcon.Conn

	// Selected database (default 0)
	dbId DatabaseId

	redis *Redis
}

// Redis gets the redis instance.
func (c *Client) Redis() *Redis {
	return c.redis
}

// ClientId get the client id.
func (c *Client) ClientId() ClientId {
	return c.clientId
}

func (c *Client) Conn() redcon.Conn {
	return c.conn
}

// SelectDb selects the clients database.
func (c *Client) SelectDb(dbId DatabaseId) {
	c.dbId = dbId
}

// DbId gets the clients selected database id.
func (c *Client) DbId() DatabaseId {
	return c.dbId
}

// Db gets the clients selected database.
func (c *Client) Db() *RedisDb {
	return c.Redis().RedisDb(c.DbId())
}

// Disconnects and removes a Client.
func (c *Client) FreeClient() {
	c.Conn().Close() // TODO should we log on error?
	delete(c.Redis().getClients(), c.ClientId())
}
