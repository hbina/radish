package cmd

import (
	"fmt"
	"strings"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/zrank/
// ZRANK key member WITHSCORE
func ZrankCommand(c *pkg.Client, args [][]byte) {
	implZrankCommand(c, args, false)
}

func implZrankCommand(c *pkg.Client, args [][]byte, reverse bool) {
	if len(args) < 3 {
		c.Conn().WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
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
				c.Conn().WriteError(util.SyntaxErr)
				return
			}
		case "withscore":
			{
				withScore = true
			}
		}
	}

	maybeSet, _ := c.Db().Get(key)

	if maybeSet == nil {
		c.Conn().WriteNull()
		return
	}

	if maybeSet.Type() != types.ValueTypeZSet {
		c.Conn().WriteError(util.WrongTypeErr)
		return
	}

	set := maybeSet.(*types.ZSet)

	node, rank := set.FindNodeByLex(memberKey)

	if node == nil || node.Key != memberKey {
		if withScore {
			// We should have a null array
			c.Conn().WriteNullArray()
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
		c.Conn().WriteBulkString(fmt.Sprint(node.Score))
	} else {
		if reverse {
			c.Conn().WriteInt(set.Len() - rank)
		} else {
			c.Conn().WriteInt(rank - 1)
		}
	}
}
