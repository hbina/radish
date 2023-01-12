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
		c.WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
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
				c.WriteError(util.SyntaxErr)
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
		c.WriteNullBulk()
		return
	}

	if maybeSet.Type() != types.ValueTypeZSet {
		c.WriteError(util.WrongTypeErr)
		return
	}

	set := maybeSet.(*types.ZSet)

	node, rank := set.FindNodeByLex(memberKey)

	if node == nil || node.Key != memberKey {
		if withScore {
			// We should have a null array
			c.WriteNullArray()
		} else {
			c.WriteNullBulk()
		}
		return
	}

	if withScore {
		c.WriteArray(2)
		if reverse {
			c.WriteBulkString(fmt.Sprint(set.Len() - rank))
		} else {
			c.WriteBulkString(fmt.Sprint(rank - 1))
		}
		c.WriteBulkString(fmt.Sprint(node.Score))
	} else {
		if reverse {
			c.WriteInt(set.Len() - rank)
		} else {
			c.WriteInt(rank - 1)
		}
	}
}
