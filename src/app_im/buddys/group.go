package main

import (
	"app_im/im_utils"
	"app_im/models"
	"app_im/protocol"
	"app_im/token"
	"encoding/json"
	"github.com/astaxie/beego/orm"
	log "github.com/cihub/seelog"
	"vipcomm"
)

// 群角色
type GroupRole int

const (
	KGrpMember  GroupRole = iota // 普通成员
	KGrpManager                  // 群管理
	KGrpOwner                    // 群主
)

var (
	kMaxGrpMembers = 500 // 群最大人数
)

type GroupController struct {
	token.BaseController
}

/*
* op      0 退出   1 加入群
* flag    是否通知
 */
func changeUserGrp(uid int64, gid int64, op int, flag bool, o orm.Ormer) (ver int64, err error) {
	ver, err = models.SetUserGrpList(uid, gid, op, o)
	if err == nil {
		im_utils.AppNotifyPeer(uid, protocol.AppMsg{
			Uri: vipcomm.Uri_GrpListChange,
			Data: protocol.PGrpListNotify{
				Uid:     uid,
				Gid:     gid,
				Op:      op,
				Version: ver,
			},
		})
	}
	return
}

// 创建群
func (this *GroupController) CreateGrp() {
	var req protocol.PCreateGrp
	var res protocol.PCreateGrpRes

	// 最终处理
	var err error
	var errStr = "success"
	var resCode = vipcomm.RES_OK
	defer func() {
		log.Infof("%v res:%v. uid:%v name:%v ver:%v gid:%v err:%v",
			errStr, resCode, req.Uid, req.GrpName, res.Version, res.Gid, err)
		this.Data["json"] = protocol.HttpRes{
			Code: resCode,
			Msg:  errStr,
			Data: res,
		}
		this.ServeJSON()
	}()

	// 校验参数
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &req)
	req.Uid = this.Uid
	if err != nil || req.Uid == 0 {
		resCode = vipcomm.RES_INVALID_ARG
		errStr = "invalid args"
		return
	}

	// 创建群Id
	o := orm.NewOrm()
	o.Using("app_im")
	res.PCreateGrp = req
	res.Password = ""
	res.Gid, err = models.NewGrp(&req, o)
	if err != nil {
		resCode = vipcomm.RES_DB_ERROR
		errStr = "create group faild"
		return
	}

	// 更新群成员
	if err = models.SetGrpMember(req.Uid, res.Gid, 1, int(KGrpOwner), o); err != nil {
		log.Errorf("add grp member faild. uid:%v gid:%v", req.Uid, res.Gid)
	}

	// 更新群成员缓存
	im_utils.AddPushIdToGrp(req.Uid, res.Gid)

	// 更新自己的群数据
	res.Version, err = changeUserGrp(req.Uid, res.Gid, 1, false, o)
	if err != nil {
		resCode = vipcomm.RES_DB_ERROR
		errStr = "set self grp list faild"
		return
	}
}

/* 加入 / 退出群
* @param uid
* @param gid
* @param opFlag   0 退出群  1 加入群
* @param notify   true 通知自己
 */
func (this *GroupController) joinOrLeave(uid, gid int64,
	opFlag int, notify bool, o orm.Ormer) (version int64, errStr, resCode string, err error) {
	// 已经在群里或者退出群，不做处理
	userInfo, err := models.GetGrpMemberByUid(uid, gid, o)
	if err == nil && userInfo.Flag == opFlag {
		version, err = models.GetUserGrpVersion(uid, o)
		if err != nil {
			errStr = "already in or existed group. get version faild"
			resCode = vipcomm.RES_DB_ERROR
		} else {
			errStr = "already in or existed group"
		}
		return
	}

	// 群主只能够解散群
	if userInfo.Role == int(KGrpOwner) {
		errStr = "owener not allow exit."
		resCode = vipcomm.RES_INVALID_ARG
		return
	}

	// 修改自己群列表
	version, err = changeUserGrp(uid, gid, opFlag, notify, o)
	if err != nil {
		resCode = vipcomm.RES_DB_ERROR
		errStr = "set self group list faild"
		return
	}

	// 修改群成员，新加入群为普通成员
	if err = models.SetGrpMember(uid, gid, opFlag, int(KGrpMember), o); err != nil {
		resCode = vipcomm.RES_DB_ERROR
		errStr = "update group member faild"
		return
	}

	// 更新群成员缓存
	if opFlag == 1 {
		im_utils.AddPushIdToGrp(uid, gid)
	} else if opFlag == 0 {
		im_utils.DelPushIdByGrp(uid, gid)
	}

	im_utils.AppNotifyByGrp(gid, protocol.AppMsg{
		Uri: vipcomm.Uri_GrpMemberChange,
		Data: protocol.PGrpMemberNotify{
			Uid: uid,
			Gid: gid,
			Op:  opFlag,
		},
	})
	return
}

// 请求加群
func (this *GroupController) JoinGrp() {
	var req protocol.PJoinGrp
	res := protocol.PJoinGrpRes{}
	uid := this.Uid

	// 最终处理
	var err error
	var errStr = "success"
	var resCode = vipcomm.RES_OK
	defer func() {
		log.Infof("%v res:%v. uid:%v gid:%v ver:%v err:%v",
			errStr, resCode, uid, req.Gid, res.Version, err)
		this.Data["json"] = protocol.HttpRes{
			Code: resCode,
			Msg:  errStr,
			Data: res,
		}
		this.ServeJSON()
	}()

	// 校验参数
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &req)
	res.PJoinGrp = req
	if err != nil || uid == 0 || req.Gid == 0 {
		resCode = vipcomm.RES_INVALID_ARG
		errStr = "invalid args"
		return
	}

	// 校验群是否存在
	o := orm.NewOrm()
	o.Using("app_im")
	grpInfo, err := models.GetGrpInfo(req.Gid, o)
	if err != nil && err == orm.ErrNoRows {
		resCode = vipcomm.RES_NOT_FOUND
		errStr = "not found group"
		return
	} else if err != nil {
		resCode = vipcomm.RES_DB_ERROR
		errStr = "get group faild"
		return
	}

	// 校验进群密码
	if grpInfo.Password != "" && grpInfo.Password != req.Passwd {
		resCode = vipcomm.RES_INVALID_PASSWORD
		errStr = "password invalid"
		return
	}

	// 校验群的最大人数
	memNumber, err := models.GetGrpMemberNumber(req.Gid, o)
	if err != nil {
		resCode = vipcomm.RES_DB_ERROR
		errStr = "get member count faild"
		return
	}
	if memNumber >= kMaxGrpMembers {
		resCode = vipcomm.RES_OVER_LIMIT
		errStr = "more than max member"
		return
	}

	// 无需验证、少人50人无需验证
	if grpInfo.Invite == 0 ||
		(grpInfo.Invite == 2 && memNumber < 50) {
		res.Version, errStr, resCode, err = this.joinOrLeave(uid, req.Gid, 1, false, o)
	} else {
		// 通知群管理验证
		var managers []int64
		managers, err = models.GetGrpManagers(req.Gid, o)
		if err != nil {
			resCode = vipcomm.RES_DB_ERROR
			errStr = "get group managers faild"
			return
		}
		im_utils.AppNotifyByUids(managers, protocol.AppMsg{
			Uri: vipcomm.Uri_JoinGrpCheck,
			Data: protocol.PJoinGrpCheckNotify{
				Uid: uid,
				Gid: req.Gid,
			},
		})
		res.Version = -1
		errStr = "wait for manager check"
		log.Infof("notify manages ok. uid:%v gid:%v managers:%v", uid, req.Gid, managers)
	}
}

// 管理员验证，加群结果
func (this *GroupController) CheckJoinGrp() {
	var (
		req     protocol.JoinCheckGrp
		res     protocol.JoinCheckGrpRes
		err     error
		errStr  = "success"
		resCode = vipcomm.RES_OK
		uid     = this.Uid
	)
	defer func() {
		log.Infof("%v res:%v. uid:%v requid:%v gid:%v agree:%v err:%v",
			errStr, resCode, uid, req.Uid, req.Gid, req.Agree, err)
		this.Data["json"] = protocol.HttpRes{
			Code: resCode,
			Msg:  errStr,
			Data: res,
		}
		this.ServeJSON()
	}()

	// 校验参数
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &req)
	res.JoinCheckGrp = req
	if err != nil || uid == 0 || req.Gid == 0 {
		resCode = vipcomm.RES_INVALID_ARG
		errStr = "invalid args"
		return
	}

	// 管理员不同意，不做任何操作
	if req.Agree == 0 {
		return
	}

	// 校验管理员身份
	o := orm.NewOrm()
	o.Using("app_im")
	memInfo, err := models.GetGrpMemberByUid(uid, req.Gid, o)
	if err != nil || memInfo.Role < int(KGrpManager) {
		resCode = vipcomm.RES_INVALID_OP
		errStr = "need group manager"
		return
	}

	// 校验群的最大人数
	memNumber, err := models.GetGrpMemberNumber(req.Gid, o)
	if err != nil {
		resCode = vipcomm.RES_DB_ERROR
		errStr = "get member count faild"
		return
	}
	if memNumber >= kMaxGrpMembers {
		resCode = vipcomm.RES_OVER_LIMIT
		errStr = "more than max member"
		return
	}

	_, errStr, resCode, err = this.joinOrLeave(req.Uid, req.Gid, 1, true, o)
	return
}

// 退出群
func (this *GroupController) LeaveGrp() {
	var req protocol.PExitGrp
	var res protocol.PExitGrpRes
	uid := this.Uid

	// 最终处理
	var err error
	var errStr = "success"
	var resCode = vipcomm.RES_OK
	defer func() {
		log.Infof("%v res:%v. uid:%v gid:%v ver:%v err:%v",
			errStr, resCode, uid, req.Gid, res.Version, err)
		this.Data["json"] = protocol.HttpRes{
			Code: resCode,
			Msg:  errStr,
			Data: res,
		}
		this.ServeJSON()
	}()

	// 校验参数
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &req)
	res.PExitGrp = req
	if err != nil || uid == 0 || req.Gid == 0 {
		resCode = vipcomm.RES_INVALID_ARG
		errStr = "invalid args"
		return
	}

	o := orm.NewOrm()
	o.Using("app_im")
	res.Version, errStr, resCode, err = this.joinOrLeave(uid, req.Gid, 0, false, o)
}

// 邀请加入群
func (this *GroupController) InviteJoinGrp() {
	var req protocol.PInviteJoinGrp
	res := protocol.PInviteJoinGrpRes{}

	// 最终处理
	var err error
	var errStr = "ok"
	var resCode = vipcomm.RES_OK
	defer func() {
		log.Infof("%v res:%v. uid:%v to:%v gid:%v err:%v",
			errStr, resCode, req.Uid, req.ToUid, req.Gid, err)
		this.Data["json"] = protocol.HttpRes{
			Code: resCode,
			Data: res,
		}
		this.ServeJSON()
	}()

	// 校验参数
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &req)
	req.Uid = this.Uid
	if err != nil || req.Uid == 0 || req.ToUid == 0 || req.Gid == 0 {
		resCode = vipcomm.RES_INVALID_ARG
		errStr = "invalid args"
		return
	}
	res.Uid = req.Uid
	res.Gid = req.Gid
	res.ToUid = req.ToUid

	// 通知对方
	inviteNotify := protocol.PInviteJoinGrpNotify{}
	inviteNotify.Uid = req.Uid
	inviteNotify.Gid = req.Gid
	inviteNotify.ToUid = req.ToUid
	im_utils.AppNotifyPeer(req.ToUid, protocol.AppMsg{
		Uri:  vipcomm.Uri_InviteJoinGrp,
		Data: inviteNotify})
	return
}

// 拉取群成员列表
func (this *GroupController) SyncGrpMembers() {
	var (
		req     protocol.PSyncGrpMembers
		res     protocol.PSyncGrpMembersRes
		err     error
		errStr  = "ok"
		resCode = vipcomm.RES_OK
		uid     = this.Uid
	)
	defer func() {
		log.Infof("%v res:%v. uid:%v gid:%v mem:%v err:%v",
			errStr, resCode, uid, req.Gid, res.Members, err)
		this.Data["json"] = protocol.HttpRes{
			Code: resCode,
			Data: res,
		}
		this.ServeJSON()
	}()

	// 校验参数
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &req)
	if err != nil || req.Gid == 0 {
		resCode = vipcomm.RES_INVALID_ARG
		errStr = "invalid args"
		return
	}
	res.Gid = req.Gid

	o := orm.NewOrm()
	o.Using("app_im")
	res.Members, err = models.GetGrpMember(req.Gid, o)
	if err != nil {
		resCode = vipcomm.RES_DB_ERROR
		errStr = "get member faild"
		return
	}
}

// 全量拉取群列表
func (this *GroupController) FullSyncGrpList() {
	uid := this.Uid
	res := protocol.PFullSyncGrpListRes{
		GInfo: []protocol.PGrpInfo{},
	}

	// 最终处理
	var err error
	var errStr = ""
	var resCode = vipcomm.RES_OK
	defer func() {
		log.Infof("%v res:%v. uid:%v ginfo:%v ver:%v err:%v",
			errStr, resCode, uid, res.GInfo, res.Version, err)
		this.Data["json"] = protocol.HttpRes{
			Code: resCode,
			Data: res,
		}
		this.ServeJSON()
	}()

	// 群列表
	o := orm.NewOrm()
	o.Using("app_im")
	glist, err := models.GetUserGrpList(uid, o)
	log.Info(glist)
	if err != nil {
		resCode = vipcomm.RES_DB_ERROR
		errStr = "get glist faild"
		return
	}

	// 群信息
	ginfos, err := models.BatchGetGrpInfo(glist, o)
	log.Info(ginfos)
	if err != nil {
		resCode = vipcomm.RES_DB_ERROR
		errStr = "get ginfo faild"
		return
	}

	// 版本信息
	res.Version, err = models.GetUserGrpVersion(uid, o)
	if err != nil {
		resCode = vipcomm.RES_DB_ERROR
		errStr = "get version faild"
		return
	}

	for _, ginfo := range ginfos {
		res.GInfo = append(res.GInfo, protocol.PGrpInfo{
			Gid:  ginfo.Id,
			Name: ginfo.Name,
		})
	}
}
