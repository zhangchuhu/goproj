package protocol

// 全量同步好友列表
// 请求地址: /buddy/full_sync/
type PFullSyncBuddy struct {
	Uid int64 `json:"uid"`
}

type PFullSyncBuddyRes struct {
	Version int64   `json:"version"`
	Bids    []int64 `json:"bids"`     // 好友
	SBids   []int64 `json:"sbids"`    // 单方好友
	Black   []int64 `json:"black"`    // 黑名单
	BeBlack []int64 `json:"be_black"` // 被黑名单
}

// 请求添加好友
// 请求地址：/buddy/addbuddy/
type PAddBuddy struct {
	Uid int64  `json:"uid"` // 请求方
	Bid int64  `json:"bid"` // 被请求
	Msg string `json:"msg"` // 附带信息
}

type PAddBuddyRes struct {
	Uid     int64 `json:"uid"`
	Bid     int64 `json:"bid"`
	Flag    int   `json:"flag"`    // 单方好友1  互为好友2
	Version int64 `json:"version"` // 好友列表、黑名单等用户资料版本号
}

// 请求加好友通知
// uri : 1
type PAddBuddyNotify struct {
	Uid   int64  `json:"uid"` // 请求方
	Bid   int64  `json:"bid"` // 被请求方
	Msg   string `json:"msg"`
	SeqId int64  `json:"seq_id"` // 请求seqid
}

// 同意、拒绝添加好友
// 请求地址：/buddy/addbuddy_answer/
type PAddBuddyAnswer struct {
	Uid   int64 `json:"uid"`    // 请求加好友方
	Bid   int64 `json:"bid"`    // 被请求方
	Agree int   `json:"agree"`  // 1 同意， 0 拒绝
	SeqId int64 `json:"seq_id"` // 请求seqid
}

type PAddBuddyAnswerRes struct {
	Uid     int64 `json:"uid"`
	Bid     int64 `json:"bid"`
	Agree   int   `json:"agree"`   // 1 同意， 0 拒绝
	Version int64 `json:"version"` // 好友列表、黑名单等用户资料版本号
}

// 好友数据变化
// uri: 2
type PBuddyChangeNotify struct {
	Uid     int64 `json:"uid"`
	Bid     int64 `json:"bid"`
	Flag    int   `json:"flag"`    // 单方好友1  互为好友2
	Version int64 `json:"version"` // 好友列表、黑名单等用户资料版本号
}

// 删除好友
// 请求地址： /buddy/del_buddy/
type PDelBuddy struct {
	Uid int64 `json:"uid"`
	Bid int64 `json:"bid"`
}

type PDelBuddyRes struct {
	Uid     int64 `json:"uid"`
	Bid     int64 `json:"bid"`
	Version int64 `json:"version"` // 好友列表、黑名单等用户资料版本号
}
