package redis

// https://redis.io/commands/zdiffstore/
// ZDIFFSTORE numkeys key [key ...] [WITHSCORES]
func ZdiffstoreCommand(c *Client, args [][]byte) {
	implZSetSetOperationCommand(c, args, true, ZSetOperationDiff, false)
}
