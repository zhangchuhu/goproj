package im_utils

import (
	_ "fmt"
	"testing"
	"vipcomm/mylog"
	"vipcomm/redisf"
	"vipcomm/zkclient"
)

func Test_PushId(t *testing.T) {
	defer mylog.FlushLog()
	conn, err := zkclient.ZkClient("")
	if err != nil {
		t.Error("get zookeeper client faild")
		return
	}
	defer conn.Close()

	key := []string{"pushid"}
	err = redisf.InitRedis(conn, key, "97")
	if err != nil {
		t.Error("init redis client faild")
		return
	}

	AddPushIdToGrp(1, 1000)
}
