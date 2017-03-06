package models

import (
	"fmt"
	"testing"
	_ "vipcomm/mysql"

	"github.com/astaxie/beego/orm"
)

func Test_BindPushId(t *testing.T) {
	o := orm.NewOrm()
	o.Using("default") // 指定库名

	uid := int64(1234)
	pushId := "vxh5v5KuKSy1234"

	// bind pushId
	err := SavePushId(uid, pushId, o)
	if err != nil {
		t.Error("绑定pushId 失败 ", err)
		return
	}

	// get PushId
	var id string
	id, err = GetPushIdByUid(uid, o)
	if err != nil {
		t.Error("获取pushId失败 ", err)
		return
	}
	fmt.Println("************", id, "**********************")
}
