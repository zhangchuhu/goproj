package protocol

//群聊拉取聊天信息的结构体
type GroupMsgData struct {
	GroupId int64  `json:"group_id"`
	Uid     int64  `json:"uid"`
	SeqId   int64  `json:"seq_id"`
	MsgId   string `json:"msg_id"`
	MsgType int32  `json:"msg_type"`
	MsgBody string `json:"msg_body"`
}
