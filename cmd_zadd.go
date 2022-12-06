package redis

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

const (
	ZaddCompareMode = iota
	ZaddCompareGt
	ZaddCompareLt
)

const (
	ZaddInsertMode = iota
	ZaddInsertNx
	ZaddInsertXx
)

// https://redis.io/commands/zadd/
// ZADD key [NX | XX] [GT | LT] [CH] [INCR] score member [score member ...]
func ZaddCommand(c *Client, args [][]byte) {
	if len(args) < 4 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])

	// Parse options
	optionCount := 0
	compareMode := ZaddCompareMode
	insertMode := ZaddInsertMode
	chEnabled := false
	incrEnabled := false

	for i := 2; i < len(args); i++ {
		arg := strings.ToLower(string(args[i]))

		// If arg is a number, then we have found score, meaning there are no options left
		_, err := strconv.ParseInt(arg, 10, 32)

		if err == nil {
			break
		}

		switch arg {
		case "xx":
			{
				if insertMode != ZaddInsertMode {
					c.Conn().WriteError(SyntaxErr)
					return
				}
				insertMode = ZaddInsertXx
				optionCount++
			}
		case "nx":
			{
				if insertMode != ZaddInsertMode || compareMode != ZaddCompareMode {
					c.Conn().WriteError(SyntaxErr)
					return
				}
				insertMode = ZaddInsertNx
				optionCount++
			}
		case "gt":
			{
				if compareMode != ZaddCompareMode || insertMode == ZaddInsertNx {
					c.Conn().WriteError(SyntaxErr)
					return
				}
				compareMode = ZaddCompareGt
				optionCount++
			}
		case "lt":
			{
				if compareMode != ZaddCompareMode || insertMode == ZaddInsertNx {
					c.Conn().WriteError(SyntaxErr)
					return
				}
				compareMode = ZaddCompareLt
				optionCount++
			}
		case "ch":
			{
				chEnabled = true
				optionCount++
			}
		case "incr":
			{
				incrEnabled = true
				optionCount++
			}
		}
	}

	// Cannot find any score member pairs
	if len(args)-(optionCount+2) == 0 {
		c.Conn().WriteError(WrongNumOfArgsErr)
		return
	}

	// Check if there are score member pairs before we even proceed
	if (len(args)-(optionCount+2))%2 == 1 {
		c.Conn().WriteError(SyntaxErr)
		return
	}

	// Validate that all the scores are valid floats
	for i := 2 + optionCount; i < len(args); i += 2 {
		score, err := strconv.ParseFloat(string(args[i]), 64)
		if err != nil || math.IsNaN(score) {
			c.Conn().WriteError(InvalidFloatErr)
			return
		}
	}

	// Redis does not support multiple score-element pair when doing INCR option
	// for some reasons...
	if incrEnabled && len(args)-optionCount-2 > 2 {
		c.Conn().WriteError(fmt.Sprintf("ERR %s option supports a single increment-element pair", "INCR"))
		return
	}

	maybeSet := c.Db().Get(key)

	if maybeSet == nil {
		maybeSet = NewZSet()
	}

	if maybeSet.Type() != ValueTypeZSet {
		c.Conn().WriteError(WrongTypeErr)
		return
	}

	set := maybeSet.Value().(*SortedSet)

	addedCount := 0
	var newScore *float64 = nil
	for i := 2 + optionCount; i+1 < len(args); i += 2 {
		// SAFETY: We already validated all scores to be valid
		score, _ := strconv.ParseFloat(string(args[i]), 64)
		member := string(args[i+1])
		old := set.GetByKey(member)

		if old != nil && incrEnabled {
			score += old.Score()
		}

		if (old != nil && insertMode == ZaddInsertNx) ||
			(old == nil && insertMode == ZaddInsertXx) ||
			(old != nil && compareMode == ZaddCompareGt && score <= old.Score()) ||
			(old != nil && compareMode == ZaddCompareLt && score >= old.Score()) {
			continue
		}

		added := set.AddOrUpdate(member, score)

		if added || (chEnabled && old != nil && old.Score() != score) {
			addedCount++
		}

		// When INCR is enabled, only 1 pair of score-member can be specified
		if incrEnabled {
			newScore = &score
			break
		}
	}

	c.Db().Set(key, NewZSetFromSs(set), time.Time{})

	if incrEnabled {
		if newScore == nil {
			c.Conn().WriteNull()
		} else {
			c.Conn().WriteString(fmt.Sprint(*newScore))
		}
	} else {
		c.Conn().WriteInt(addedCount)
	}
}
