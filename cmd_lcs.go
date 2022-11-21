package redis

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

// https://redis.io/commands/lcs/
// LCS key1 key2 [LEN] [IDX] [MINMATCHLEN len] [WITHMATCHLEN]
func LcsCommand(c *Client, args [][]byte) {
	if len(args) < 3 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	isLen := false
	isIdx := false
	var minMatchLen *int
	isWithMatchLen := false

	// Parse the optional arguments
	for i := 3; i < len(args); i++ {
		arg := strings.ToLower(string(args[i]))
		switch arg {
		case "len":
			isLen = true
		case "idx":
			isIdx = true
		case "minmatchlen":
			if i+1 == len(args) {
				c.Conn().WriteError(SyntaxErr)
			}

			i++

			len64, err := strconv.ParseInt(string(args[i]), 10, 32)

			if err != nil || len64 < 0 {
				log.Println(err)
				c.Conn().WriteError(InvalidIntErr)
				return
			}

			len := int(len64)
			minMatchLen = &len
		case "withmatchlen":
			isWithMatchLen = true
		default:
			c.Conn().WriteError(SyntaxErr)
			return
		}
	}

	db := c.Db()
	key1 := string(args[1])
	key2 := string(args[2])

	maybeValueX, _ := db.GetOrExpire(key1, true)
	maybeValueY, _ := db.GetOrExpire(key2, true)

	if maybeValueY == nil ||
		maybeValueY.Type() != ValueTypeString ||
		maybeValueX == nil ||
		maybeValueX.Type() != ValueTypeString {
		c.Conn().WriteError(WrongTypeErr)
		return
	}

	valueX := maybeValueX.(*String)
	valueY := maybeValueY.(*String)

	// Add 1 to each axes to include the initial 0s
	lcs := make([][]int, 0, valueY.Len()+1)

	for i := 0; i < valueY.Len()+1; i++ {
		lcs = append(lcs, make([]int, valueX.Len()+1))
	}

	setLcs := func(x, y, value int) {
		lcs[y][x] = value
	}

	getLcs := func(x, y int) int {
		return lcs[y][x]
	}

	/* Start building the LCS table. */
	for y := 0; y <= valueY.Len(); y++ {
		for x := 0; x <= valueX.Len(); x++ {
			if y == 0 || x == 0 {
				setLcs(x, y, 0)
			} else if valueY.Get(y-1) == valueX.Get(x-1) {
				setLcs(x, y, getLcs(x-1, y-1)+1)
			} else {
				lcs1 := getLcs(x-1, y)
				lcs2 := getLcs(x, y-1)

				if lcs1 > lcs2 {
					setLcs(x, y, lcs1)
				} else {
					setLcs(x, y, lcs2)
				}
			}
		}
	}

	found := false
	valueXMatchIdx := make([]int, 0)
	valueYMatchIdx := make([]int, 0)
	xrange_end := valueX.Len()
	yrange_end := valueY.Len()

	// Work out the matches backward
	for x, y := valueX.Len(), valueY.Len(); x > 0 && y > 0; {
		if valueX.Get(x-1) == valueY.Get(y-1) {
			if !found {
				found = true
				xrange_end = x - 1
				yrange_end = y - 1
			}

			x--
			y--

			if x == 0 && y == 0 && found {
				valueXMatchIdx = append(valueXMatchIdx, 0, xrange_end)
				valueYMatchIdx = append(valueYMatchIdx, 0, yrange_end)
			}
		} else {
			if found {
				valueXMatchIdx = append(valueXMatchIdx, x, xrange_end)
				valueYMatchIdx = append(valueYMatchIdx, y, yrange_end)
				found = false
			}

			if getLcs(x-1, y) > getLcs(x, y-1) {
				x--
			} else {
				y--
			}

		}
	}

	if len(valueXMatchIdx)%2 != 0 ||
		len(valueYMatchIdx)%2 != 0 ||
		len(valueXMatchIdx) != len(valueYMatchIdx) {
		log.Println("valueXMatchIdx or valueYMatchIdx is not even or not equal:", valueXMatchIdx, valueYMatchIdx)
	}

	matchCount := len(valueXMatchIdx) / 2

	// Filter out matches that doesn't satisfy min match length
	if minMatchLen != nil {
		valueXMatchIdxN := make([]int, 0)
		valueYMatchIdxN := make([]int, 0)
		for i := 0; i < matchCount; i++ {
			if valueXMatchIdx[i*2+1]-valueXMatchIdx[i*2]+1 >= *minMatchLen {
				valueXMatchIdxN = append(valueXMatchIdxN,
					valueXMatchIdx[i*2],
					valueXMatchIdx[i*2+1])
				valueYMatchIdxN = append(valueYMatchIdxN,
					valueYMatchIdx[i*2],
					valueYMatchIdx[i*2+1])
			}
		}
		valueXMatchIdx = valueXMatchIdxN
		valueYMatchIdx = valueYMatchIdxN
		matchCount = len(valueXMatchIdxN) / 2
	}

	var result strings.Builder
	// TODO: Calculate the total length
	for i := matchCount - 1; i >= 0; i-- {
		result.WriteString(valueX.inner[valueXMatchIdx[i*2] : valueXMatchIdx[(i*2)+1]+1])
	}

	if isIdx {
		c.Conn().WriteArray(4)
		c.Conn().WriteBulkString("matches")
		c.Conn().WriteArray(matchCount)
		for i := 0; i < matchCount; i++ {
			if isWithMatchLen {
				c.Conn().WriteArray(3) // we are going to write the 2 matches and the match length
			} else {
				c.Conn().WriteArray(2) // we are going to write the 2 matches
			}
			c.Conn().WriteArray(2) // we are going to write the start and end of matchX
			c.Conn().WriteInt(valueXMatchIdx[i*2])
			c.Conn().WriteInt(valueXMatchIdx[i*2+1])
			c.Conn().WriteArray(2) // we are going to write the start and end of matchX
			c.Conn().WriteInt(valueYMatchIdx[i*2])
			c.Conn().WriteInt(valueYMatchIdx[i*2+1])
			if isWithMatchLen {
				c.Conn().WriteInt(valueXMatchIdx[i*2+1] - valueXMatchIdx[i*2] + 1)
			}
		}
		c.Conn().WriteBulkString("len")
		c.Conn().WriteInt(len(result.String()))
	} else if isLen {
		c.Conn().WriteInt(getLcs(valueX.Len(), valueY.Len()))
	} else {
		c.Conn().WriteBulkString(result.String())
	}
}
