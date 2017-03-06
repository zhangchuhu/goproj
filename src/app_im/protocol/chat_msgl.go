package protocol

type HttpCommRes struct {
	Code string `json:"code"`
}
type HttpRes struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
type AppMsg struct {
	Uri  string      `json:"uri"`
	Data interface{} `json:"data"`
}

//系统状态栏推送协议
type Android struct {
	Title   string      `json:"title"`
	Message string      `json:"message"`
	PayLoad interface{} `json:"payload"`
}
type Apn struct {
	Alert   string      `json:"altert"`
	Badge   int32       `json:"badge"`
	Sound   string      `json:"sound"`
	PayLoad interface{} `json:"payload"`
}

type SysMsg struct {
	JAndroid Android `json:"android"`
	JApn     Apn     `json:"apn"`
}

//单聊未读取消息通知
type PChatMsgNotify struct {
	MsgId int64 `json:"msgid"`
}

//聊天消息
type ChatMsg struct {
	Id      int64  `json:"id"`
	FromUid int64  `json:"fromUid"`
	ToUid   int64  `json:"toUid"`
	MsgType int32  `json:"msgType"`
	MsgBody string `json:"msgBody"`
}

type PostMsg struct {
	Code string `json:"code"`
}

// 绑定push_id
type PBindPush struct {
	Uid    int64  `json:"uid"`
	PushId string `json:"pushid"`
}
