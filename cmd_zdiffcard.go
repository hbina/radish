package redis

// https://redis.io/commands/zdiffcard/
// ZDIFFCARD numkeys key [key ...] [WITHSCORES]
func ZdiffcardCommand(c *Client, args [][]byte) {
	implZSetSetOperationCommand(c, args, false, ZSetOperationDiff, true)
}
