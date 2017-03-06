package models

import (
	"github.com/astaxie/beego/orm"
)

// 群版本信息
type DUserGrpVersion struct {
	Id      int64 `orm:"auto"`
	Uid     int64 `orm:"column(uid_)"`
	Version int64 `orm:"column(version)"`
}

func init() {
	orm.RegisterModel(new(DUserGrpVersion))
}

func IncrUserGrpVersion(uid int64, o orm.Ormer) (err error) {
	_, err = o.Raw("insert into tbl_user_grp_version(uid_, version, last_modify) values(?, 1, now()) on duplicate key update version = version + 1, last_modify = now()", uid).Exec()
	return
}

// 群列表版本信息
func GetUserGrpVersion(uid int64, o orm.Ormer) (version int64, err error) {
	var info DUserGrpVersion
	err = o.Raw("select version from tbl_user_grp_version where uid_=?", uid).QueryRow(&info)
	if err == nil {
		version = info.Version
	}
	return
}
