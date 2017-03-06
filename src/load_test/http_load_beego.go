package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
	"vipcomm/confs"
	"vipcomm/mylog"
	"vipcomm/redisf"
)

var (
	log   = mylog.Logger
	count = 0
)

type RedisTestController struct {
	beego.Controller
}

func (this *RedisTestController) RedisHandler() {
	res, err := redis.String(redisf.Do("redis_load", "GET", "1234"))
	if err != nil {
		//fmt.Printf("error >> %v\n", err)
	} else {
		//fmt.Println(res)
	}

	count++
	if count%1000 == 0 {
		fmt.Println("1000" + res)
	}
}

func main() {
	defer mylog.FlushLog()
	confs.InitRedisConf("redis.conf")
	redisf.InitRedis()

	beego.Router("/redis_test_bee", &RedisTestController{}, "GET:RedisHandler")
	beego.Run()
}
