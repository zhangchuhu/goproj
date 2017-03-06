package main

import (
	"app_im/im_utils"
	"app_im/models"
	"app_im/protocol"
	"encoding/json"
	_ "fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	log "github.com/cihub/seelog"
)

type ChatMsgController struct {
	beego.Controller
}

// Examples:
//
//   req: POST /im/BindPushId/ {"Uid":"1","PushId":"test"}
//   res: 200                  { "Code": "0", "Msg": "" }

func (this *ChatMsgController) BindPushId() {

	var res protocol.HttpRes
	defer func() {
		this.Data["json"] = res
		this.ServeJSON()
	}()
	req := struct {
		Uid    int64  `json:"uid"`
		PushId string `json:"pushId"`
	}{}

	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &req); err != nil {
		this.Ctx.Output.SetStatus(400)
		res.Code = "-1"
		res.Msg = err.Error()
		return
	}
	log.Infof("uid:%v, pushId:%v", req.Uid, req.PushId)
	o := orm.NewOrm()
	o.Using("app_im") // 指定库名
	if glist, err := models.GetUserGrpList(req.Uid, o); err != nil {
		res.Code = "-1"
		res.Msg = err.Error()
	} else {
		log.Infof("glist=%v", glist)
		for _, groupId := range glist {
			im_utils.DelPushIdByGrp(req.Uid, groupId)
			im_utils.UpdatePushIdToGrp(req.Uid, req.PushId, groupId)
		}
	}
	if err := im_utils.BindPushId(req.Uid, req.PushId); err != nil {
		res.Code = "-1"
		res.Msg = err.Error()
	} else {
		res.Code = "0"
		res.Msg = ""
	}
	return
}

// Examples:
//
//   req: POST /SendChatMsg/ {"Seq":"1","FromUid":10,"ToUid":20,"Msg":"test"}
//   res: 200 empty title    { "Code": "0", "Msg": "" }

func (this *ChatMsgController) SendChatMsg() {
	var res protocol.HttpRes
	defer func() {
		this.Data["json"] = res
		this.ServeJSON()
	}()

	req := struct {
		Seq     string `json:"seq"`
		FromUid int64  `json:"fromUid"`
		ToUid   int64  `json:"toUid"`
		MsgType int32  `json:"msgType"`
		MsgBody string `json:"msgBody"`
	}{}

	data := [...]struct {
		MsgId int64 `json:"msgId"`
	}{{0}}
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &req); err != nil {
		this.Ctx.Output.SetStatus(400)
		res.Code = "-1"
		res.Msg = err.Error()
		return
	}
	log.Infof("seq:%v,fromUid:%v,toUid:%v,msgType:%v,msgBody:%v", req.Seq, req.FromUid, req.ToUid, req.MsgType, req.MsgBody)
	//	msgId, err := redis.Int64(redisf.Do("chat_msg", "INCR", "msg_id"))
	//	if err != nil {
	//		this.Ctx.Output.SetStatus(400)
	//		res.Code = vipcomm.RES_DB_ERROR
	//		res.Msg = err.Error()
	//		return
	//	}
	o := orm.NewOrm()
	o.Using("default") // 指定库名
	msgId, err := models.SaveChatMsg(req.Seq, req.FromUid, req.ToUid, req.MsgType, req.MsgBody, o)
	if err != nil {
		this.Ctx.Output.SetStatus(400)
		res.Code = "-1"
		res.Msg = err.Error()
		return
	}
	data[0].MsgId = msgId
	//下发通知消息：状态栏通知
	//	android := protocol.Android{"YY Live", "你有一条未读取消息", protocol.PChatMsgNotify{1, msgId}}
	//	apn := protocol.Apn{"你有一条未读取消息", 1, "default", protocol.PChatMsgNotify{1, msgId}}
	//	var sysMsg protocol.SysMsg
	//	sysMsg.JAndroid = android
	//	sysMsg.JApn = apn
	//	im_utils.SysNotifyPeer(req.ToUid, sysMsg)
	//下发通知消息:App内透传
	im_utils.AppNotifyPeer(req.ToUid, protocol.AppMsg{"6", protocol.PChatMsgNotify{msgId}})
	res.Code = "0"
	res.Msg = "sucess"
	res.Data = data
	return
}
func (this *ChatMsgController) GetChatMsg() {
	log.Infof("sync chat msg")
	var res protocol.HttpRes
	defer func() {
		this.Data["json"] = res
		this.ServeJSON()
	}()
	req := struct {
		Id  int64 `json:"id"`
		Uid int64 `json:"uid"`
	}{}
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &req); err != nil {
		this.Ctx.Output.SetStatus(400)
		res.Code = "-1"
		res.Msg = err.Error()
		return
	}
	log.Infof("id:%v, uid:%v", req.Id, req.Uid)
	if req.Id == 0 {
		req.Id, _ = im_utils.GetMsgIdByUid(req.Uid)
	}
	o := orm.NewOrm()
	o.Using("app_im") // 指定库名
	data, err := models.GetChatMsg(req.Uid, req.Id, o)
	if err != nil {
		this.Ctx.Output.SetStatus(400)
		res.Code = "-1"
		res.Msg = err.Error()
		return
	}
	//记录已读取的最大msgId
	if len(data) > 0 {
		im_utils.SetMsgId(req.Uid, data[len(data)-1].Id)
	}
	res.Code = "0"
	res.Msg = ""
	res.Data = data
	return

}
