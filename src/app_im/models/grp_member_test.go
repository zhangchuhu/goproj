package models

import (
	"testing"
	"vipcomm/mysql"
	"vipcomm/zkclient"
	"github.com/astaxie/beego/orm"
)

func init() {
	conn, _:= zkclient.ZkClient("")
    defer  conn.Close()

	mysql.InitDbOrm(conn, []string {"app_im"}, "app_im", "105")
}

func Test_UpdateGrpMemberFlag(t *testing.T) {
	o := orm.NewOrm()
	o.Using("app_im")
	err := UpdateGrpMemberFlag(11111, 1111, 0, o)
	if err != nil {
		t.Error("更新失败")
		return
	}
}

func Test_GetGrpMemberByUid(t *testing.T) {
	o := orm.NewOrm()
	t.Log(GetGrpMemberByUid(4, 33, o))
}
