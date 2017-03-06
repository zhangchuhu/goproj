package models

import (
	"github.com/astaxie/beego/orm"
)

type DBuddyVersion struct {
	Id      int64 `orm:"auto"`
	Uid     int64 `orm:"column(uid_)"`
	Version int64 `orm:"column(version_)"`
}

func init() {
	orm.RegisterModel(new(DBuddyVersion))
}

// 设置版本号
func IncBuddyVersion(uid int64, o orm.Ormer) error {
	_, err := o.Raw("insert into buddyinfo_version(uid_, version_, last_modify) values(?, 1, now()) on duplicate key update version_=version_+1, last_modify=now()", uid).Exec()
	return err
}

func GetBuddyVersion(uid int64, o orm.Ormer) (version int64, err error) {
	var info DBuddyVersion
	err = o.Raw("select uid_, version_ from buddyinfo_version where uid_=?", uid).QueryRow(&info)
	if err == nil {
		version = info.Version
	}
	return
}
