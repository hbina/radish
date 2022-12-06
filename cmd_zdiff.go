package redis

// https://redis.io/commands/zdiff/
// ZDIFF numkeys key [key ...] [WITHSCORES]
func ZdiffCommand(c *Client, args [][]byte) {
	implZSetSetOperationCommand(c, args, false, ZSetOperationDiff, false)
}
