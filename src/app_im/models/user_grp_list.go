package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"strconv"
)

// 群列表
type DUserGrpList struct {
	Id     int64 `orm:"auto"`
	Uid    int64 `orm:"column(uid_)"`
	Gid    int64 `orm:"column(gid_)"`
	Flag   int   `orm:"column(flag_)"`
	Notify int   `orm:"column(notify)"`
	Block  int   `orm:"column(block)"`
}

func init() {
	orm.RegisterModel(new(DUserGrpList))
}

func genUgName(uid int64) string {
	return fmt.Sprintf("tbl_user_grp_%d", uid%1)
}

/*
* 更新用户群数据（加入的群）
* @param op    0 退群， 1 加群
 */
func SetUserGrpList(uid int64, gid int64, op int, o orm.Ormer) (version int64, err error) {
	_, err = o.Raw("insert into "+genUgName(uid)+"(uid_, gid_, flag_, last_modify) values(?, ?, ?, now()) on duplicate key update flag_ = ?, last_modify=now()",
		uid, gid, op, op).Exec()
	if err != nil {
		return
	}

	// 更新群列表版本号
	if err = IncrUserGrpVersion(uid, o); err != nil {
		return
	}

	version, err = GetUserGrpVersion(uid, o)
	return
}

/*
* @return flag   0 不在群里，1 在群里
 */
func GetUserGrpListInfo(uid int64, gid int64, o orm.Ormer) (flag int, err error) {
	var res orm.ParamsList
	sql := "select flag_ from " + genUgName(uid) + " where uid_ = ? and gid_ = ?"
	n, err := o.Raw(sql, uid, gid).ValuesFlat(&res)
	if err == nil && n > 0 {
		flag, _ = strconv.Atoi(res[0].(string))
	} else if err == orm.ErrNoRows {
		err = nil
	}
	return
}

func GetUserGrpList(uid int64, o orm.Ormer) (glist []int64, err error) {
	var res orm.ParamsList
	sql := "select gid_ from " + genUgName(uid) + " where uid_=? and flag_ = 1"
	n, err := o.Raw(sql, uid).ValuesFlat(&res)
	if err == nil && n > 0 {
		for _, i := range res {
			gid, _ := strconv.ParseInt(i.(string), 10, 64)
			glist = append(glist, gid)
		}
	}
	return
}
