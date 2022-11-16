package redis

import (
	"fmt"
	"time"
)

// https://redis.io/commands/setex/
// SETEX key seconds value
func SetexCommand(c *Client, args [][]byte) {
	if len(args) != 4 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	seconds := string(args[2])
	value := string(args[3])

	newTtl, err := ParseExpiryTime(seconds, uint64(time.Second))

	if err != nil {
		c.Conn().WriteError(InvalidIntErr)
		return
	}

	// This is currently buggy because genericSetCommand will perform some writes to the connection that we don't want to happen.
	// We could add a stub connection here but the better solution is for this function to return Result<T, Error>
	// Oh well...
	genericSetCommand(c, key, value, time.Time{}, SetWriteMode, false)
	genericExpireCommand(c, key, newTtl, ExpireMode)
}
