package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	log "github.com/cihub/seelog"
)

type GroupMsgInfo struct {
	Id      int64  `orm:"column(id)" json:"-"`
	GroupId int64  `orm:"column(group_id)" json:"groupId"`
	Uid     int64  `orm:"column(uid)" json:"uid"`
	SeqId   int64  `orm:"column(seq_id)" json:"seqId"`
	MsgId   string `orm:"column(msg_id)" json:"msgId"`
	MsgType int32  `orm:"column(msg_type)" json:"msgType"`
	MsgBody string `orm:"column(msg_body)" json:"msgBody"`
}

func genGroupMsgTableName(group_id int64) string {
	return fmt.Sprintf("group_msg_%d", group_id%1)
}

func init() {
	orm.RegisterModel(new(GroupMsgInfo))
}

func InsertGroupMsg(group_id int64, uid int64, msg_id string, msg_type int32, msg_body string, o orm.Ormer) error {
	sql := "insert into " + genGroupMsgTableName(group_id) +
		" (group_id, uid,  msg_id, msg_type, msg_body, last_modify) values(?, ?, ?, ?, ?, now() )"

	_, err := o.Raw(sql, group_id, uid, msg_id, msg_type, msg_body).Exec()
	if err != nil {
		log.Errorf("[models.InsertGroupMsg] sql exec err group_id=%v, uid=%v, msg_id=%v, msg_type=%v, msg_body=%v", group_id, uid, msg_id, msg_type, msg_body)
	}

	return err
}

//从from_seq_id 到 to_seq_id的消息
func BatchGetGroupMsgByInterval(group_id int64, from_seq_id int64, to_seq_id int64, o orm.Ormer) (*[]GroupMsgInfo, error) {
	var groupInfos = []GroupMsgInfo{}
	sql := "select * from " + genGroupMsgTableName(group_id) + " where group_id = ? and seq_id > ? and seq_id < ?"

	_, err := o.Raw(sql, group_id, from_seq_id, to_seq_id).QueryRows(&groupInfos)
	if err != nil {
		log.Errorf("[models.BatchGetGroupMsg] err group_id=%v old_seq_id=%v new_seq_id=%v", group_id, from_seq_id, to_seq_id)
	}
	return &groupInfos, err
}

func BatchGetGroupMsg(group_id int64, seq_id int64, o orm.Ormer) (*[]GroupMsgInfo, error) {
	var groupInfos = []GroupMsgInfo{}
	sql := "select * from " + genGroupMsgTableName(group_id) + " where group_id = ? and seq_id > ?"

	_, err := o.Raw(sql, group_id, seq_id).QueryRows(&groupInfos)
	if err != nil {
		log.Errorf("[models.BatchGetGroupMsg] err group_id=%v seq_id=%v ", group_id, seq_id)
	}
	return &groupInfos, err
}

func UpdateGroupMsgSeqId(msg_id string, group_id int64, seq_id int64, o orm.Ormer) error {
	sql := "update " + genGroupMsgTableName(group_id) + " set seq_id = ? where msg_id = ?"

	_, err := o.Raw(sql, seq_id, msg_id).Exec()
	if err != nil {
		log.Errorf("update err msg_id=%v, group_id=%v, seq_id=%v", msg_id, group_id, seq_id)
	}
	return err
}
