package models

import (
	"testing"

	"github.com/astaxie/beego/orm"
)

func buddyExist(bid int64, binfos []int64) bool {
	for _, v := range binfos {
		if v == bid {
			return true
		}
	}
	return false
}

func Test_AddBuddy(t *testing.T) {
	o := orm.NewOrm()
	o.Using("app_im") // 指定库名

	uid := int64(1234)
	bid := int64(12345)

	// get version
	ver, err := GetBuddyVersion(uid, o)
	if err != nil {
		t.Error("获取好友版本失败 ", err)
		return
	}

	// add buddy
	err = SetBuddy(uid, bid, 1, o)
	if err != nil {
		t.Error("添加好友失败 ", err)
		return
	}

	ver1, err := GetBuddyVersion(uid, o)
	if err != nil || ver1 != ver+1 {
		t.Error("添加好友后版本号不一致 ", err, ver1, ver)
		return
	}

	// get buddy
	bids, sbids, err := GetAllBuddys(uid, o)
	if err != nil || !buddyExist(bid, sbids) {
		t.Error("查询好友失败 ", err)
		return
	}

	err = SetBuddy(uid, bid, 2, o)
	if err != nil {
		t.Error("添加好友失败 ", err)
		return
	}

	bids, sbids, err = GetAllBuddys(uid, o)
	t.Log(bids, sbids)
}
