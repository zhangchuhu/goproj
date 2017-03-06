package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	log "github.com/cihub/seelog"
)

// 好友信息表
type BuddyInfo struct {
	Id   int64 `orm:"auto"`
	Uid  int64 `orm:"column(uid_)"`
	Bid  int64 `orm:"column(bid_)"`
	Flag int   `orm:"column(flag_)"` // 0 deleted  1 单方好友 2 互为好友
}

func init() {
	orm.RegisterModel(new(BuddyInfo))
}

// buddyinfo table name
func genTableName(uid int64) string {
	return fmt.Sprintf("buddyinfo_%d", uid%1)
}

// 添加好友
func SetBuddy(uid int64, bid int64, flag int, o orm.Ormer) (err error) {
	sql := "insert into " + genTableName(uid) +
		"(uid_, bid_, flag_, last_modify) values(?, ?, ?, now()) on duplicate key update flag_=?, last_modify=now()"
	_, err = o.Raw(sql, uid, bid, flag, flag).Exec()
	if err != nil {
		return
	}

	// 修改版本号
	err = IncBuddyVersion(uid, o)
	return
}

func GetBuddy(uid int64, bid int64, o orm.Ormer) (info BuddyInfo, err error) {
	err = o.Raw("select uid_, bid_, flag_ from "+genTableName(uid)+" where uid_ = ? and bid_ = ?", uid, bid).QueryRow(&info)
	return
}

func GetAllBuddys(uid int64, o orm.Ormer) (bids []int64, sbids []int64, err error) {
	bids = make([]int64, 0)
	sbids = make([]int64, 0)
	var binfos = []BuddyInfo{}
	_, err = o.Raw("select uid_, bid_, flag_ from "+genTableName(uid)+" where uid_ = ? and flag_ > 0",
		uid).QueryRows(&binfos)
	if err != nil {
		log.Errorf("get buddys faild. uid:%v err:%v", uid, err)
	}

	for _, i := range binfos {
		if i.Flag == 1 {
			sbids = append(sbids, i.Bid)
		} else if i.Flag == 2 {
			bids = append(bids, i.Bid)
		}
	}
	return
}
