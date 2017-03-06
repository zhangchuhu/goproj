package models

import (
	"app_im/protocol"
	"fmt"

	"github.com/astaxie/beego/orm"
	log "github.com/cihub/seelog"
)

func genTableName1(uid int64) string {
	return fmt.Sprintf("chat_msg_%d", uid%1)
}

func init() {
	orm.RegisterModel(new(protocol.ChatMsg))
}

// 新增消息
func SaveChatMsg(seq string, from_uid int64, to_uid int64, msg_type int32, msg_body string, o orm.Ormer) (int64, error) {
	sql := "insert into " + genTableName1(to_uid) +
		"( seq, from_uid, to_uid, msg_type, msg_body, last_modify) values(?, ?, ?, ?, ?, now()) on duplicate key update last_modify=now() "
	result, err := o.Raw(sql, seq, from_uid, to_uid, msg_type, msg_body).Exec()
	if err != nil {
		log.Errorf("add chat msg faild. from_uid:%v to_uid:%v err:%v", from_uid, to_uid, err)
	}
	id, _ := result.LastInsertId()
	return id, err
}

// 获取未读取的聊天消息
func GetChatMsg(uid int64, id int64, o orm.Ormer) ([]protocol.ChatMsg, error) {
	var chatmsgs = []protocol.ChatMsg{}
	_, err := o.Raw("select id,from_uid,to_uid,msg_type,msg_body from "+genTableName1(uid)+" where to_uid = ? and id = ? order by id asc", uid, id).QueryRows(&chatmsgs)
	if err != nil {
		log.Errorf("get chat msgs faild. uid:%v id %v err:%v", uid, id, err)
	}
	return chatmsgs, err
}
