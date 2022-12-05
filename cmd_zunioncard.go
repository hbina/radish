package redis

// https://redis.io/commands/zunioncard/
// ZUNIONCARD numkeys key [key ...] [LIMIT limit]
func ZunioncardCommand(c *Client, args [][]byte) {
	implZSetSetOperationCommand(c, args, false, ZSetOperationUnion, true)
}
