package confs

import (
	"github.com/astaxie/beego/config"
	"vipcomm/mylog"
)

var (
	DbConfs    config.Configer
	RedisConfs config.Configer
	RpcConfs   config.Configer
)

func init() {
	defer mylog.FlushLog()
	InitDbConf("../conf/mysql.conf")
	InitRedisConf("../conf/redis.conf")
	InitRpcConf("../conf/rpc.conf")
}

// 数据库
func InitDbConf(confFile string) error {
	var err error
	DbConfs, err = config.NewConfig("ini", confFile)
	if err != nil {
		mylog.Logger.Error(err)
	}
	return err
}

// redis
func InitRedisConf(confFile string) error {
	var err error
	RedisConfs, err = config.NewConfig("ini", confFile)
	if err != nil {
		mylog.Logger.Error(err)
	}
	return err
}

// rpc
func InitRpcConf(confFile string) error {
	var err error
	RpcConfs, err = config.NewConfig("ini", confFile)
	if err != nil {
		mylog.Logger.Error(err)
	}
	return err
}
