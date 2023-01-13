package pkg

import (
	"fmt"

	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// A connected Client.
type Client struct {
	conn  *util.Conn
	dbId  uint64
	redis *Redis
	R3    bool
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
	c.R3 = false
}

func (c *Client) UseResp3() {
	c.R3 = true
}

func (c *Client) WriteToConn(nodes []*types.SortedSetNode, withScores bool) bool {
	if c.R3 {
		if withScores {
			ok := c.Conn().WriteArray(len(nodes))

			if !ok {
				return false
			}

			for _, node := range nodes {
				ok = c.Conn().WriteArray(2)

				if !ok {
					return false
				}

				ok = c.Conn().WriteBulkString(node.Key)

				if !ok {
					return false
				}

				ok = c.Conn().WriteFloat64(node.Score)

				if !ok {
					return false
				}
			}
		} else {
			ok := c.Conn().WriteArray(len(nodes))

			if !ok {
				return false
			}

			for _, node := range nodes {
				ok = c.Conn().WriteBulkString(node.Key)

				if !ok {
					return false
				}
			}
		}
	} else {
		if withScores {
			ok := c.Conn().WriteArray(len(nodes) * 2)

			if !ok {
				return false
			}

			for _, node := range nodes {
				if !ok {
					return false
				}

				ok = c.Conn().WriteBulkString(node.Key)

				if !ok {
					return false
				}

				ok = c.Conn().WriteBulkString(fmt.Sprint(node.Score))

				if !ok {
					return false
				}
			}
		} else {
			ok := c.Conn().WriteArray(len(nodes))

			if !ok {
				return false
			}

			for _, node := range nodes {
				ok = c.Conn().WriteBulkString(node.Key)

				if !ok {
					return false
				}
			}
		}
	}

	return true
}
