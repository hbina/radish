package redis

// https://redis.io/commands/zintercard/
// ZINTERCARD numkeys key [key ...] [LIMIT limit]
func ZintercardCommand(c *Client, args [][]byte) {
	implZSetSetOperationCommand(c, args, false, ZSetOperationInter, true)
}
