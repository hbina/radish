package cmd

import "github.com/hbina/radish/internal/pkg"

func GenerateCommands() map[string]*pkg.Command {
	arr := []*pkg.Command{
		pkg.NewCommand("ping", PingCommand, pkg.CMD_READONLY),
		pkg.NewCommand("set", SetCommand, pkg.CMD_WRITE),
		pkg.NewCommand("get", GetCommand, pkg.CMD_READONLY),
		pkg.NewCommand("del", DelCommand, pkg.CMD_WRITE),
		pkg.NewCommand("ttl", TtlCommand, pkg.CMD_READONLY),
		pkg.NewCommand("lpush", LPushCommand, pkg.CMD_WRITE),
		pkg.NewCommand("rpush", RPushCommand, pkg.CMD_WRITE),
		pkg.NewCommand("lpop", LPopCommand, pkg.CMD_WRITE),
		pkg.NewCommand("rpop", RPopCommand, pkg.CMD_WRITE),
		pkg.NewCommand("lrange", LRangeCommand, pkg.CMD_READONLY),
		pkg.NewCommand("config", ConfigCommand, pkg.CMD_WRITE),
		pkg.NewCommand("info", InfoCommand, pkg.CMD_READONLY),
		pkg.NewCommand("select", SelectCommand, pkg.CMD_WRITE),
		pkg.NewCommand("flushall", FlushAllCommand, pkg.CMD_WRITE),
		pkg.NewCommand("function", FunctionCommand, pkg.CMD_WRITE),
		pkg.NewCommand("incr", IncrCommand, pkg.CMD_WRITE),
		pkg.NewCommand("incrby", IncrByCommand, pkg.CMD_WRITE),
		pkg.NewCommand("incrbyfloat", IncrByFloatCommand, pkg.CMD_WRITE),
		pkg.NewCommand("decr", DecrCommand, pkg.CMD_WRITE),
		pkg.NewCommand("decrby", DecrByCommand, pkg.CMD_WRITE),
		pkg.NewCommand("decrbyfloat", DecrByFloatCommand, pkg.CMD_WRITE),
		pkg.NewCommand("object", ObjectCommand, pkg.CMD_READONLY),
		pkg.NewCommand("sadd", SaddCommand, pkg.CMD_WRITE),
		pkg.NewCommand("smembers", SmembersCommand, pkg.CMD_WRITE),
		pkg.NewCommand("smismember", SmismemberCommand, pkg.CMD_WRITE),
		pkg.NewCommand("zadd", ZaddCommand, pkg.CMD_WRITE),
		pkg.NewCommand("dump", DumpCommand, pkg.CMD_READONLY),
		pkg.NewCommand("exists", ExistsCommand, pkg.CMD_READONLY),
		pkg.NewCommand("restore", RestoreCommand, pkg.CMD_WRITE),
		pkg.NewCommand("pttl", PttlCommand, pkg.CMD_READONLY),
		pkg.NewCommand("debug", DebugCommand, pkg.CMD_READONLY),
		pkg.NewCommand("srem", SremCommand, pkg.CMD_WRITE),
		pkg.NewCommand("sintercard", SintercardCommand, pkg.CMD_READONLY),
		pkg.NewCommand("sinter", SinterCommand, pkg.CMD_READONLY),
		pkg.NewCommand("sinterstore", SinterstoreCommand, pkg.CMD_WRITE),
		pkg.NewCommand("scard", ScardCommand, pkg.CMD_READONLY),
		pkg.NewCommand("sismember", SismemberCommand, pkg.CMD_READONLY),
		pkg.NewCommand("sunion", SunionCommand, pkg.CMD_READONLY),
		pkg.NewCommand("sunionstore", SunionstoreCommand, pkg.CMD_WRITE),
		pkg.NewCommand("sdiff", SdiffCommand, pkg.CMD_READONLY),
		pkg.NewCommand("sdiffstore", SdiffstoreCommand, pkg.CMD_WRITE),
		pkg.NewCommand("spop", SpopCommand, pkg.CMD_WRITE),
		pkg.NewCommand("srandmember", SrandmemberCommand, pkg.CMD_READONLY),
		pkg.NewCommand("smove", SmoveCommand, pkg.CMD_WRITE),
		pkg.NewCommand("watch", WatchCommand, pkg.CMD_READONLY),
		pkg.NewCommand("multi", MultiCommand, pkg.CMD_READONLY),
		pkg.NewCommand("exec", ExecCommand, pkg.CMD_READONLY),
		pkg.NewCommand("flushdb", FlushDbCommand, pkg.CMD_WRITE),
		pkg.NewCommand("dbsize", DbSizeCommand, pkg.CMD_READONLY),
		pkg.NewCommand("setx", SetXCommand, pkg.CMD_WRITE),
		pkg.NewCommand("setnx", SetNxCommand, pkg.CMD_WRITE),
		pkg.NewCommand("expire", ExpireCommand, pkg.CMD_WRITE),
		pkg.NewCommand("setex", SetexCommand, pkg.CMD_WRITE),
		pkg.NewCommand("getex", GetexCommand, pkg.CMD_READONLY),
		pkg.NewCommand("getdel", GetdelCommand, pkg.CMD_WRITE),
		pkg.NewCommand("mget", MgetCommand, pkg.CMD_WRITE),
		pkg.NewCommand("getset", GetsetCommand, pkg.CMD_WRITE),
		pkg.NewCommand("mset", MsetCommand, pkg.CMD_READONLY),
		pkg.NewCommand("msetnx", MsetnxCommand, pkg.CMD_WRITE),
		pkg.NewCommand("strlen", StrlenCommand, pkg.CMD_READONLY),
		pkg.NewCommand("setbit", SetbitCommand, pkg.CMD_WRITE),
		pkg.NewCommand("getbit", GetbitCommand, pkg.CMD_READONLY),
		pkg.NewCommand("setrange", SetrangeCommand, pkg.CMD_WRITE),
		pkg.NewCommand("getrange", GetrangeCommand, pkg.CMD_READONLY),
		pkg.NewCommand("lcs", LcsCommand, pkg.CMD_READONLY),
		pkg.NewCommand("zrange", ZrangeCommand, pkg.CMD_READONLY),
		pkg.NewCommand("type", TypeCommand, pkg.CMD_READONLY),
		pkg.NewCommand("zcard", ZcardCommand, pkg.CMD_READONLY),
		pkg.NewCommand("zscore", ZscoreCommand, pkg.CMD_READONLY),
		pkg.NewCommand("zincrby", ZincrbyCommand, pkg.CMD_WRITE),
		pkg.NewCommand("zrem", ZremCommand, pkg.CMD_WRITE),
		pkg.NewCommand("zrevrange", ZrevrangeCommand, pkg.CMD_READONLY),
		pkg.NewCommand("zrank", ZrankCommand, pkg.CMD_READONLY),
		pkg.NewCommand("zrevrank", ZrevrankCommand, pkg.CMD_READONLY),
		pkg.NewCommand("zrangebyscore", ZrangebyscoreCommand, pkg.CMD_READONLY),
		pkg.NewCommand("zrevrangebyscore", ZrevrangebyscoreCommand, pkg.CMD_READONLY),
		pkg.NewCommand("zcount", ZcountCommand, pkg.CMD_READONLY),
		pkg.NewCommand("zrangebylex", ZrangebylexCommand, pkg.CMD_READONLY),
		pkg.NewCommand("zrevrangebylex", ZrevrangebylexCommand, pkg.CMD_READONLY),
		pkg.NewCommand("zlexcount", ZlexcountCommand, pkg.CMD_READONLY),
		pkg.NewCommand("zremrangebyscore", ZremrangebyscoreCommand, pkg.CMD_WRITE),
		pkg.NewCommand("zremrangebylex", ZremrangebylexCommand, pkg.CMD_WRITE),
		pkg.NewCommand("zremrangebyrank", ZremrangebyrankCommand, pkg.CMD_WRITE),
		pkg.NewCommand("zinter", ZinterCommand, pkg.CMD_READONLY),
		pkg.NewCommand("zintercard", ZintercardCommand, pkg.CMD_READONLY),
		pkg.NewCommand("zinterstore", ZinterstoreCommand, pkg.CMD_WRITE),
		pkg.NewCommand("zunion", ZunionCommand, pkg.CMD_READONLY),
		pkg.NewCommand("zunioncard", ZunioncardCommand, pkg.CMD_READONLY),
		pkg.NewCommand("zunionstore", ZunionstoreCommand, pkg.CMD_WRITE),
		pkg.NewCommand("zdiff", ZdiffCommand, pkg.CMD_READONLY),
		pkg.NewCommand("zdiffcard", ZdiffcardCommand, pkg.CMD_READONLY),
		pkg.NewCommand("zdiffstore", ZdiffstoreCommand, pkg.CMD_WRITE),
		pkg.NewCommand("hello", HelloCommand, pkg.CMD_WRITE),
		pkg.NewCommand("zpopmin", ZpopminCommand, pkg.CMD_WRITE),
		pkg.NewCommand("zpopmax", ZpopmaxCommand, pkg.CMD_WRITE),
		pkg.NewCommand("zmpop", ZmpopCommand, pkg.CMD_WRITE),
	}

	res := make(map[string]*pkg.Command, len(arr))

	for _, r := range arr {
		res[r.Name] = r
	}

	return res
}

func GenerateBlockingCommands() map[string]*pkg.BlockingCommand {
	arr := []*pkg.BlockingCommand{
		pkg.NewBlockingCommand("bzmpop", BzmpopCommand, pkg.CMD_WRITE),
	}

	res := make(map[string]*pkg.BlockingCommand, len(arr))

	for _, r := range arr {
		res[r.Name] = r
	}

	return res
}
