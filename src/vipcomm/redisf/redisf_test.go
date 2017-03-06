package redisf

import (
	"github.com/garyburd/redigo/redis"
	"testing"
	"vipcomm"
)

func TestRedisCache(t *testing.T) {
	conn, err := vipcomm.ZkClient("")
	if err != nil {
		t.Error("get zookeeper client faild")
		return
	}
	defer conn.Close()

	key := "redis_example"
	err = InitRedis(conn, key, "97")
	if err != nil {
		t.Error("init redis client faild")
		return
	}

	// set
	res, err := redis.String(Do(key, "SET", "test_abcd", "1"))
	if err != nil || res != "OK" {
		t.Error("SET test faild")
	}
	t.Log("SET      ", res, err)

	// get
	v2, err := redis.Int(Do(key, "GET", "test_abcd"))
	if err != nil || v2 != 1 {
		t.Error("GET test_abcd 返回值不对", v2, err)
	}
	t.Log("GET      ", v2, err)

	// exists
	v1, err := redis.Bool(Do(key, "EXISTS", "test_abcd"))
	if err != nil || v1 != true {
		t.Error("EXISTS test faild")
	}
	t.Log("Exists   ", v1, err)

	v3, err := redis.String(Do("vip_logic", "GET", "test_abcd"))
	if err == nil {
		t.Error("不存在配置未返回错误", v3, err)
	}
	t.Log("不存在配置情况 ", v3, err)

}
