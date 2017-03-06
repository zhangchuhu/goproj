package vipcomm

import (
	"encoding/json"
	"vipcomm/mysql"
	"vipcomm/redisf"
	"vipcomm/zkclient"

	"github.com/astaxie/beego"
	log "github.com/cihub/seelog"
)

var (
	GroupId = "0" // 机房id
)

type ZkServerConf struct {
	Redis        []string `json:"redis"`
	RedisCluster []string `json:"redis_cluster"`
	Db           []string `json:"db"`         // 数据库
	DefautDb     string   `json:"default_db"` // 默认数据库
	RedisMq      []string `json:"redis_mq"`   // redis消息队列
	Check        string   `json:"check"`      // 检查主备机房
	Thrift       []string `json:"thrift"`
	Desc         string   `json:desc"`
}

func InitServer() error {
	// appName
	var path = "/vip_servers/" + beego.BConfig.AppName
	conn, err := zkclient.ZkClient("")
	if err != nil {
		log.Errorf("get zookeeper conn faild err:%v", err)
		return err
	}
	defer conn.Close()

	value, _, err := conn.Get(path)
	if err != nil {
		log.Errorf("get conf faild. path:%v err:%v", path, err)
		return err
	}

	var sconf ZkServerConf
	err = json.Unmarshal(value, &sconf)
	if err != nil {
		log.Errorf("unmarshal conf faild. path:%v value:%v err:%v", path, string(value), err)
		return err
	}

	// 机房id
	if GroupId = beego.AppConfig.String("groupId"); GroupId == "" {
		GroupId = "0"
	}
	log.Infof("group id: %v", GroupId)

	// 初始化mysql
	if len(sconf.Db) > 0 {
		if err = mysql.InitDbOrm(conn, sconf.Db, sconf.DefautDb, GroupId); err != nil {
			return err
		}
	}

	// init redis
	if len(sconf.Redis) > 0 {
		if err = redisf.InitRedis(conn, sconf.Redis, GroupId); err != nil {
			return err
		}
	}
	//init rediscluster
	if len(sconf.RedisCluster) > 0 {
		log.Infof("*******redis cluste******: %v", sconf.RedisCluster)
		if err = redisf.InitRedisCluster(conn, sconf.RedisCluster); err != nil {
			return err
		}
	}
	return nil
}
