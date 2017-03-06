package im_utils

import (
	"errors"
	_ "fmt"
	"strconv"
	"vipcomm/redisf"

	log "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
)

// 绑定uid和pushId
func BindPushId(uid int64, pushid string) error {
	db := "pushid"
	key := "pushid|" + strconv.FormatInt(uid, 10)
	res, err := redis.String(redisf.Do(db, "SET", key, pushid))
	if err != nil || res != "OK" {
		return err
	} else {
		return nil
	}
}

//根据uid获取pushid
func GetPushIdByUid(uid int64) (string, error) {
	db := "pushid"
	key := "pushid|" + strconv.FormatInt(uid, 10)
	res, err := redis.String(redisf.Do(db, "GET", key))
	if err != nil {
		return "", err
	} else if res == "" {
		return "", errors.New("pushId is empty")
	} else {
		return res, nil
	}
}

//获取群成员的pushId列表
func GetPushIdByGrp(groupId int64) ([]string, error) {
	db := "pushid"
	key := "grp_pushid|" + strconv.FormatInt(groupId, 10)
	res, err := redis.Strings(redisf.Do(db, "LRANGE", key, 0, -1))
	return res, err
}

//用户退群删除pushId
func DelPushIdByGrp(uid int64, groupId int64) (bool, error) {
	db := "pushid"
	key := "grp_pushid|" + strconv.FormatInt(groupId, 10)
	v, _ := GetPushIdByUid(uid)
	res, err := redis.Bool(redisf.Do(db, "LREM", key, 0, v))
	return res, err
}

//用户加群增加pushId
func AddPushIdToGrp(uid int64, groupId int64) (bool, error) {
	db := "pushid"
	DelPushIdByGrp(uid, groupId)
	key := "grp_pushid|" + strconv.FormatInt(groupId, 10)
	v, err := GetPushIdByUid(uid)
	if err != nil {
		log.Errorf("get pushid faild uid:%v gid:%v err:%v", uid, groupId, err)
		return false, err
	}
	return redis.Bool(redisf.Do(db, "LPUSH", key, v))

}

//pushId变更，修改群pushId列表
func UpdatePushIdToGrp(uid int64, pushId string, groupId int64) (bool, error) {
	db := "pushid"
	DelPushIdByGrp(uid, groupId)
	key := "grp_pushid|" + strconv.FormatInt(groupId, 10)
	return redis.Bool(redisf.Do(db, "LPUSH", key, pushId))

}
