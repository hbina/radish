package redis

import (
	"fmt"
	"strings"
)

// https://redis.io/commands/zrank/
// ZRANK key member WITHSCORE
func ZrankCommand(c *Client, args [][]byte) {
	implZrankCommand(c, args, false)
}

func implZrankCommand(c *Client, args [][]byte, reverse bool) {
	if len(args) < 3 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	memberKey := string(args[2])
	withScore := false

	// Parse additional option
	for i := 3; i < len(args); i++ {
		switch strings.ToLower(string(args[i])) {
		default:
			{
				c.Conn().WriteError(SyntaxErr)
				return
			}
		case "withscore":
			{
				withScore = true
			}
		}
	}

	maybeSet := c.Db().Get(key)

	if maybeSet == nil {
		c.Conn().WriteNull()
		return
	}

	if maybeSet.Type() != ValueTypeZSet {
		c.Conn().WriteError(WrongTypeErr)
		return
	}

	set := maybeSet.(*ZSet)

	node, rank := set.inner.findNodeByLex(memberKey)

	if node == nil || node.key != memberKey {
		if withScore {
			// We should have a null array
			c.Conn().WriteRaw([]byte("*-1\r\n"))
		} else {
			c.Conn().WriteNull()
		}
		return
	}

	if withScore {
		c.Conn().WriteArray(2)
		if reverse {
			c.Conn().WriteBulkString(fmt.Sprint(set.Len() - rank))
		} else {
			c.Conn().WriteBulkString(fmt.Sprint(rank - 1))
		}
		c.Conn().WriteBulkString(fmt.Sprint(node.score))
	} else {
		if reverse {
			c.Conn().WriteInt(set.Len() - rank)
		} else {
			c.Conn().WriteInt(rank - 1)
		}
	}
}
