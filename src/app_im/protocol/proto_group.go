package protocol

// 创建群   /group/create_grp/
type PCreateGrp struct {
	Uid        int64  `json:"uid"`         // 群主uid
	GrpName    string `json:"grp_name"`    // 群名称 （必填）
	Logo       string `json:"logo"`        // 群logo （必填）
	Location   string `json:"location"`    // 位置
	Desc       string `json:"desc"`        // 群公告
	Privilege  int    `json:"privilege"`   // 0 普通群  1 私密群  （必填）
	Password   string `json:"password"`    // 进群密码
	InviteType int    `json:"invite_type"` // 0 无需审核、1 管理员审核、 2 少于50人直接进群 （必填）
}

type PCreateGrpRes struct {
	PCreateGrp
	Gid     int64 `json:"gid"`     // 群ID
	Version int64 `json:"version"` // 群主的群列表version
}

// 群list变化通知
// uri: 3
type PGrpListNotify struct {
	Uid     int64 `json:"uid"`
	Gid     int64 `json:"gid"`
	Op      int   `json:"op"`      // 0 退出、1 加入
	Version int64 `json:"version"` // 群资料版本信息
}

/*
* 加入群  /group/join_grp/
* 返回值
*        -101 群不存在
*        -103 进群需要验证密码，密码不对
*        -104 群人数超过上限
 */
type PJoinGrp struct {
	Gid    int64  `json:"gid"`
	Passwd string `json:"passwd"` // 加群密码，如果有的话
}

type PJoinGrpRes struct {
	PJoinGrp
	Version int64 `json:"version"` // 返回-1，表示正在等管理员验证
}

/*
* 退出群  /group/exit_grp/
 */
type PExitGrp struct {
	Gid int64 `json:"gid"`
}

type PExitGrpRes struct {
	PExitGrp
	Version int64 `json:"version"` // 申请人群资料版本信息
}

/*
* 管理员验证加群结果
* 返回值
*        -101 群不存在
*        -105 非管理员
*        -104 人数超过群的最大限制
 */
type JoinCheckGrp struct {
	Uid   int64 `json:"uid"` // 申请入群的人
	Gid   int64 `json:"gid"`
	Agree int   `json:"agree"` // 1 同意、0 不同意
}

type JoinCheckGrpRes struct {
	JoinCheckGrp
}

// 群成员列表变化通知
// uri: 4
type PGrpMemberNotify struct {
	Uid int64 `json:"uid"`
	Gid int64 `json:"gid"`
	Op  int   `json:"op"` // 0 退出、1 加入
}

// 邀请加入群  /group/invite/
type PInviteJoinGrp struct {
	Uid   int64 `json:"uid"`    // 邀请方
	ToUid int64 `json:"to_uid"` // 被邀请方
	Gid   int64 `json:"gid"`
}

type PInviteJoinGrpRes struct {
	PInviteJoinGrp
}

// 邀请进群通知
// 5
type PInviteJoinGrpNotify struct {
	PInviteJoinGrp
}

// 拉取群成员列表  /group/sync_member/
type PSyncGrpMembers struct {
	Gid int64 `json:"gid"`
}

type PSyncGrpMembersRes struct {
	Gid     int64   `json:"gid"`
	Members []int64 `json:"members"`
}

// 全量拉取群信息  /group/full_sync/
type PFullSyncGrpList struct {
}

type PGrpInfo struct {
	Gid        int64  `json:"gid"`
	Name       string `json:"name"`
	Logo       string `json:"logo"`        // 群logo （必填）
	Location   string `json:"location"`    // 位置 格式："经度,维度"
	Desc       string `json:"desc"`        // 群公告
	Privilege  int    `json:"privilege"`   // 0 普通群  1 私密群  （必填）
	Password   int    `json:"password"`    // 是否有进群密码
	InviteType int    `json:"invite_type"` // 0 无需审核、1 管理员审核、 2 少于50人直接进群 （必填）
}

type PFullSyncGrpListRes struct {
	GInfo   []PGrpInfo `json:"ginfo"`
	Version int64      `json:"version"`
}

// 加入群请求，通知管理员验证
// uri: 6
type PJoinGrpCheckNotify struct {
	Uid int64 `json:"uid"`
	Gid int64 `json:"gid"`
}

// 群角色
type PGrpMemberInfo struct {
	Role int `json:"role"` // 0 普通成员、1 管理员、2 群主
}
