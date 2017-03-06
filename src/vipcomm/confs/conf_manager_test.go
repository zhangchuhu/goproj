package confs

import (
	"reflect"
	"testing"
	"vipcomm/mylog"
)

func Test_Confs(t *testing.T) {
	defer mylog.FlushLog()
	t.Log("DbConfs 开始测试 ... ")
	if DbConfs != nil && DbConfs != reflect.Zero(reflect.TypeOf(DbConfs)).Interface() {
		t.Log(DbConfs.Strings("all"))
	} else {
		t.Log("DbConfs is nil")
	}
}
