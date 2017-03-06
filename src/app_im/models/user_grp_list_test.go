package models

import (
	"testing"
	_ "vipcomm/mysql"

	"github.com/astaxie/beego/orm"
)

func Test_AddBuddy1(t *testing.T) {
	o := orm.NewOrm()
	o.Using("app_im")
	glist, err := GetUserGrpList(1234, o)
	if err != nil {
		t.Error("获取群列表失败")
		return
	} else {
		for _, gid := range glist {
			t.Log("gid: ", gid)
		}
	}

	// 批量获取群信息
	ginfos, err := BatchGetGrpInfo(glist, o)
	if err != nil {
		t.Error("批量获取群信息失败")
		return
	} else {
		for _, ginfo := range ginfos {
			t.Log("ginfo: ", ginfo)
		}
	}
}
