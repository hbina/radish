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

func (c *Client) Close() error {
	return c.conn.Close()
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

func (c *Client) WriteToConn(nodes []*types.SortedSetNode, withScores bool, handleSingleElement bool) {
	err := func() error {
		if c.R3 && withScores {
			if len(nodes) == 1 && handleSingleElement {
				res := nodes[0]
				err := c.Conn().WriteArray(2)

				if err != nil {
					return err
				}

				err = c.Conn().WriteBulkString(res.Key)

				if err != nil {
					return err
				}

				err = c.Conn().WriteFloat64(res.Score)

				if err != nil {
					return err
				}
			} else {
				err := c.Conn().WriteArray(len(nodes))

				if err != nil {
					return err
				}

				for _, node := range nodes {
					err = c.Conn().WriteArray(2)

					if err != nil {
						return err
					}

					err = c.Conn().WriteBulkString(node.Key)

					if err != nil {
						return err
					}

					err = c.Conn().WriteFloat64(node.Score)

					if err != nil {
						return err
					}
				}
			}
		} else {
			if withScores {
				err := c.Conn().WriteArray(len(nodes) * 2)

				if err != nil {
					return err
				}

				for _, node := range nodes {
					err = c.Conn().WriteBulkString(node.Key)

					if err != nil {
						return err
					}

					err = c.Conn().WriteBulkString(fmt.Sprint(node.Score))

					if err != nil {
						return err
					}
				}

			} else {
				err := c.Conn().WriteArray(len(nodes))

				if err != nil {
					return err
				}

				for _, node := range nodes {
					err = c.Conn().WriteBulkString(node.Key)

					if err != nil {
						return err
					}
				}
			}
		}

		return nil
	}()

	if err != nil {
		util.Logger.Printf("Failed to write to connection: '%s'\n", err)

		err := c.Conn().Close()

		if err != nil {
			util.Logger.Printf("Unable to close connection: '%s'\n", err)
		}
	}
}
