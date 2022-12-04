package redis

// https://redis.io/commands/zrevrank/
// ZREVRANK key member WITHSCORE
func ZrevrankCommand(c *Client, args [][]byte) {
	implZrankCommand(c, args, true)
}
