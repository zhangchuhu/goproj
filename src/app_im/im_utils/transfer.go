package im_utils

import (
	log "github.com/cihub/seelog"
	"vipcomm"
)

//应用内透传
// 通知单个uid
func AppNotifyPeer(uid int64, req interface{}) error {
	pushIds := []string{}
	pushId, err := GetPushIdByUid(uid)
	if err == nil {
		pushIds = append(pushIds, pushId)
		vipcomm.PushClientPack(pushIds, req)
		log.Infof("[AppNotifyPeer] success,req:%v", req)
		return nil
	}
	log.Infof("[AppNotifyPeer] err:%v", err)
	return err
}

// 通知uids
func AppNotifyByUids(uids []int64, req interface{}) error {
	for _, uid := range uids {
		AppNotifyPeer(uid, req)
	}
	return nil
}

// 通知群
func AppNotifyByGrp(gid int64, req interface{}) error {
	pushIds, _ := GetPushIdByGrp(gid)
	err := vipcomm.PushClientPack(pushIds, req)
	return err
}

//状态栏通知

// 通知单个uid
func SysNotifyPeer(uid int64, req interface{}) error {
	pushIds := []string{}
	pushId, err := GetPushIdByUid(uid)
	if err == nil {
		pushIds = append(pushIds, pushId)
		vipcomm.SysNotify(pushIds, req)
		return nil
	}
	return err
}

// 通知uids
func SysNotifyByUids(uids []int64, req interface{}) error {
	for _, uid := range uids {
		SysNotifyPeer(uid, req)
	}
	return nil
}

// 通知群
func SysNotifySysByGrp(gid int64, req interface{}) error {
	pushIds, _ := GetPushIdByGrp(gid)
	err := vipcomm.SysNotify(pushIds, req)
	return err
}
