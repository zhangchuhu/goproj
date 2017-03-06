package im_utils

import (
	_ "fmt"
	"strconv"
	"vipcomm/redisf"

	"github.com/garyburd/redigo/redis"
)

// 记录已读的消息Id
func SetMsgId(uid int64, msgId int64) error {
	db := "chat_msg"
	key := "msgid|" + strconv.FormatInt(uid, 10)
	res, err := redis.String(redisf.Do(db, "SET", key, msgId))
	if err != nil || res != "OK" {
		return err
	} else {
		return nil
	}
}

//查询用户已读取的消息ID
func GetMsgIdByUid(uid int64) (int64, error) {
	db := "chat_msg"
	key := "msgid|" + strconv.FormatInt(uid, 10)
	res, err := redis.Int64(redisf.Do(db, "GET", key))
	if err != nil {
		return 0, err
	} else {
		return res, nil
	}
}
