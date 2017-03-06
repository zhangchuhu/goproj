package models

import (
	_ "vipcomm/mylog"

	"github.com/astaxie/beego/orm"
	log "github.com/cihub/seelog"
)

type PushIdInfo struct {
	Uid    int64
	PushId string
}

// 根据uid获取pushid
func GetPushIdByUid(uid int64, o orm.Ormer) (string, error) {
	var res PushIdInfo
	err := o.Raw("select uid,pushId from push_id where uid = ?", uid).QueryRow(&res)
	if err != nil {
		log.Errorf("get chat msgs faild. uid:%v  err:%v", uid, err)
	}
	log.Infof("======GetPushId. uid:%v pushId:%v", res.Uid, res.PushId)
	return res.PushId, err
}

// 根据群id获取群的pushid
func GetPushIdByGrp(gid int64) []string {
	return []string{"vip_app_im_2001", "vip_app_im_2002"}
}

// 绑定uid跟pushId
func SavePushId(uid int64, pushId string, o orm.Ormer) error {
	sql := "insert into push_id " +
		"(uid, pushId, last_modify) values(?, ?, now()) on duplicate key update pushId=?,last_modify=now() "
	_, err := o.Raw(sql, uid, pushId, pushId).Exec()
	if err != nil {
		log.Errorf("save push id faild. uid:%v err:%v", uid, err)
	}
	return err
}
