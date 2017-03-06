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

type BuddyController struct {
	token.BaseController
}

// 通知加好友
func notifyAddBuddy(uid int64, bid int64, msg string) {
	err := im_utils.AppNotifyPeer(bid, protocol.AppMsg{
		Uri: vipcomm.Uri_AddBuddy,
		Data: protocol.PAddBuddyNotify{
			Uid:   uid,
			Bid:   bid,
			Msg:   msg,
			SeqId: 0,
		},
	})
	if err != nil {
		log.Infof("notify add buddy faild. uid:%v bid:%v err:%v", uid, bid, err)
	} else {
		log.Infof("notify add buddy ok. uid:%v bid:%v", uid, bid)
	}
}

// 更新好友数据，并通知
func setBuddyInfo(uid int64, bid int64, flag int, notify bool, o orm.Ormer) (ver int64, err error) {
	err = models.SetBuddy(uid, bid, flag, o)
	if err != nil {
		log.Errorf("set buddy faild. uid:%v bid:%v flag:%v err:%v", uid, bid, flag, err)
		return
	} else {
		log.Infof("set buddy ok. uid:%v bid:%v flag:%v", uid, bid, flag)
	}

	ver, err = models.GetBuddyVersion(uid, o)
	if err != nil {
		log.Errorf("get buddy version faild. uid:%v bid:%v flag:%v err:%v", uid, bid, flag, err)
		return
	}

	if notify {
		err1 := im_utils.AppNotifyPeer(uid, protocol.AppMsg{
			Uri: vipcomm.Uri_BuddyChange,
			Data: protocol.PBuddyChangeNotify{
				Uid:     uid,
				Bid:     bid,
				Flag:    flag,
				Version: ver,
			},
		})
		if err1 != nil {
			log.Errorf("notify buddy change faild. uid:%v bid:%v flag:%v ver:%v err:%v",
				uid, bid, flag, ver, err1)
		} else {
			log.Infof("notify buddy change ok. uid:%v bid:%v flag:%v ver:%v",
				uid, bid, flag, ver)
		}
	}
	return
}

// 请求添加好友
func (this *BuddyController) AddBuddy() {
	var req protocol.PAddBuddy
	res := protocol.PAddBuddyRes{
		Flag:    1, // 单方好友
		Version: -1,
	}

	// 最终处理
	var err error
	var resCode = vipcomm.RES_OK
	var errStr = "success"
	defer func() {
		log.Infof("%v res:%v. uid:%v bid:%v flag:%v ver:%v err:%v",
			errStr, resCode, req.Uid, req.Bid, res.Flag, res.Version, err)
		this.Data["json"] = &protocol.HttpRes{
			Code: resCode,
			Msg:  errStr,
			Data: &res,
		}
		this.ServeJSON()
	}()

	// 校验参数
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &req)
	req.Uid = this.Uid
	if err != nil || req.Uid == 0 || req.Bid == 0 {
		resCode = vipcomm.RES_INVALID_ARG
		errStr = "invalid args"
		return
	}
	res.Uid = req.Uid
	res.Bid = req.Bid

	// 查看自己是不是对方的好友
	o := orm.NewOrm()
	o.Using("app_im")
	peerInfo, err := models.GetBuddy(req.Bid, req.Uid, o) // 对方好友数据
	if err != nil && err != orm.ErrNoRows {
		resCode = vipcomm.RES_DB_ERROR
		errStr = "get peer buddy faild"
		return
	}

	selfInfo, err := models.GetBuddy(req.Uid, req.Bid, o)
	if err != nil && err != orm.ErrNoRows {
		resCode = vipcomm.RES_DB_ERROR
		errStr = "get self buddy faild"
		return
	}

	if peerInfo.Flag == 0 {
		// 自己不是对方的好友, 通知对方加好友事件
		notifyAddBuddy(req.Uid, req.Bid, req.Msg)
	} else if peerInfo.Flag == 1 {
		// 自己已经是对方的好友
		// 修改B的单方好友状态
		res.Flag = 2
		_, err := setBuddyInfo(req.Bid, req.Uid, 2, true, o)
		if err != nil {
			resCode = vipcomm.RES_DB_ERROR
			errStr = "set peer buddy faild"
			return
		}
	} else if peerInfo.Flag == 2 {
		// 已经是好友关系了，重复请求？
		res.Flag = 2
		resCode = vipcomm.RES_REPEATE
		errStr = "already friend"
	}

	// 修改自己的好友数据
	if selfInfo.Flag != res.Flag {
		res.Version, err = setBuddyInfo(req.Uid, req.Bid, res.Flag, false, o)
		if err != nil {
			resCode = vipcomm.RES_DB_ERROR
			errStr = "set self buddy faild"
		}
	} else {
		res.Version, err = models.GetBuddyVersion(req.Uid, o)
		if err != nil {
			resCode = vipcomm.RES_DB_ERROR
			errStr = "get self version faild"
		}
	}
}

// 同意加好友
func (this *BuddyController) AddBuddyAnswer() {
	var req protocol.PAddBuddyAnswer
	res := protocol.PAddBuddyAnswerRes{
		Version: -1,
	}

	var err error
	var resCode = vipcomm.RES_OK
	var errStr string
	defer func() {
		this.Data["json"] = &protocol.HttpRes{
			Code: resCode,
			Data: &res,
		}
		this.ServeJSON()
		log.Infof("%v res:%v uid:%v bid:%v seq:%v agree:%v ver:%v err:%v",
			errStr, resCode, req.Uid, req.Bid, req.SeqId, req.Agree, res.Version, err)
	}()

	// 校验参数
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &req)
	req.Uid = this.Uid
	if err != nil || req.Uid == 0 || req.Bid == 0 {
		resCode = vipcomm.RES_INVALID_ARG
		errStr = "invalid args"
		return
	}
	res.Uid = req.Uid
	res.Bid = req.Bid
	res.Agree = req.Agree

	// 同意加好友，更新好友数据
	o := orm.NewOrm()
	o.Using("app_im")
	if req.Agree == 1 {
		// 更新请求方
		_, err = setBuddyInfo(req.Uid, req.Bid, 2, true, o)
		if err != nil {
			resCode = vipcomm.RES_DB_ERROR
			errStr = "set peer buddy faild"
			return
		}

		// 更新被请求方
		res.Version, err = setBuddyInfo(req.Bid, req.Uid, 2, false, o)
		if err != nil {
			resCode = vipcomm.RES_DB_ERROR
			errStr = "set self buddy faild"
		}
	} else {
		res.Version, err = models.GetBuddyVersion(req.Bid, o)
	}
}

// 删除好友
func (this *BuddyController) DelBuddy() {
	var req protocol.PDelBuddy
	res := protocol.PDelBuddyRes{
		Version: -1,
	}

	var err error
	var resCode = vipcomm.RES_OK
	var errStr string
	defer func() {
		this.Data["json"] = &protocol.HttpRes{
			Code: resCode,
			Data: &res,
		}
		this.ServeJSON()
		log.Infof("%v res:%v uid:%v bid:%v ver:%v err:%v",
			errStr, resCode, req.Uid, req.Bid, res.Version, err)
	}()

	// 校验参数
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &req)
	req.Uid = this.Uid
	if err != nil || req.Uid == 0 || req.Bid == 0 {
		resCode = vipcomm.RES_INVALID_ARG
		errStr = "invalid args"
		return
	}
	res.Uid = req.Uid
	res.Bid = req.Bid

	o := orm.NewOrm()
	o.Using("app_im")
	binfo, err := models.GetBuddy(req.Bid, req.Uid, o) // 对方好友数据
	if err != nil {
		resCode = vipcomm.RES_DB_ERROR
		errStr = "get peer buddy faild"
		return
	}

	// 如果双方互为好友关系，则自己成为对方的单方好友
	if binfo.Flag == 2 {
		_, err = setBuddyInfo(req.Bid, req.Uid, 1, true, o)
		if err != nil {
			resCode = vipcomm.RES_DB_ERROR
			errStr = "set peer buddy faild"
			return
		}
	}

	// 删除自己的好友数据
	res.Version, err = setBuddyInfo(req.Uid, req.Bid, 0, false, o)
	if err != nil {
		resCode = vipcomm.RES_DB_ERROR
		errStr = "set self buddy faild"
		return
	}
}

// 全量同步好友列表
func (this *BuddyController) FullSyncBuddy() {
	var req protocol.PFullSyncBuddy
	res := protocol.PFullSyncBuddyRes{
		Version: -1,
		Bids:    []int64{},
		SBids:   []int64{},
		Black:   []int64{},
		BeBlack: []int64{},
	}

	var err error
	var resCode = vipcomm.RES_OK
	var errStr string
	defer func() {
		this.Data["json"] = &protocol.HttpRes{
			Code: resCode,
			Data: &res,
		}
		this.ServeJSON()
		log.Infof("%v res:%v uid:%v ver:%v bid:%v sbids:%v err:%v",
			errStr, resCode, req.Uid, res.Version, res.Bids, res.SBids, err)
	}()

	if err = json.Unmarshal(this.Ctx.Input.RequestBody, &req); err != nil {
		resCode = vipcomm.RES_INVALID_ARG
		errStr = "invalid args"
		return
	}
	req.Uid = this.Uid

	// 好友列表
	o := orm.NewOrm()
	o.Using("app_im")
	res.Bids, res.SBids, err = models.GetAllBuddys(req.Uid, o)
	if err != nil {
		resCode = vipcomm.RES_DB_ERROR
		errStr = "get buddys faild"
		return
	}

	res.Version, err = models.GetBuddyVersion(req.Uid, o)
	if err != nil {
		resCode = vipcomm.RES_DB_ERROR
		errStr = "get buddys version faild"
		return
	}

	// 黑名单
	// to-do

	// 被黑名单
	// to-do
}
