package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"net/http"
	"vipcomm/confs"
	"vipcomm/mylog"
	"vipcomm/redisf"
)

var (
	log   = mylog.Logger
	count = 0
)

func redisHandler(w http.ResponseWriter, r *http.Request) {
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

	http.HandleFunc("/redis_test", redisHandler)
	http.ListenAndServe("0.0.0.0:4000", nil)
}
