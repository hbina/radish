package redis

import (
	"sync/atomic"

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
	db DatabaseId

	redis *Redis
}

// NewClient creates new client and adds it to the redis.
func (r *Redis) NewClient(conn redcon.Conn) *Client {
	c := &Client{
		conn:     conn,
		redis:    r,
		clientId: r.NextClientId(),
	}
	return c
}

// NextClientId atomically gets and increments a counter to return the next client id.
func (r *Redis) NextClientId() ClientId {
	id := atomic.AddUint64(&r.nextClientId, 1)
	return ClientId(id)
}

// Clients gets the current connected clients.
func (r *Redis) Clients() Clients {
	return r.clients
}

func (r *Redis) getClients() Clients {
	return r.clients
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
func (c *Client) SelectDb(db DatabaseId) {
	c.db = db
}

// DbId gets the clients selected database id.
func (c *Client) DbId() DatabaseId {
	return c.db
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
