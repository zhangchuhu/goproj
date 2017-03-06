package main

import (
	"app_im/models"
	"encoding/json"

	"app_im/im_utils"
	"app_im/protocol"
	"app_im/token"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	log "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"vipcomm"
	"vipcomm/mylog"
	"vipcomm/redisf"
)

type GroupMsgController struct {
	token.BaseController
}

type sendGroupReq struct {
	GroupId int64  `json:"groupId"`
	MsgId   string `json:"msgId"`
	MsgType int32  `json:"msgType"`
	MsgBody string `json:"msgBody"`
}

type pullGroupMsgReq struct {
	GroupId   int64 `json:"groupId"`
	FromSeqId int64 `json:"fromSeqId"`
	ToSeqId   int64 `json:"toSeqId"`
}

type SeqIdData struct {
	SeqId int64 `json:"seqId"`
}

type sendGroupMsgRes struct {
	Code string     `json:"code"`
	Msg  string     `json:"msg"`
	Data *SeqIdData `json:"data"`
}

type pullGroupMsgRes struct {
	Code string                 `json:"code"`
	Msg  string                 `json:"msg"`
	Data *[]models.GroupMsgInfo `json:"data"`
}

func (this *GroupMsgController) SendGroupMsg() {
	res := sendGroupMsgRes{Code: vipcomm.RES_OK, Msg: "success"}
	req := sendGroupReq{}

	defer func() {
		log.Infof("res.Code=%v, res.Msg=%v", res.Code, res.Msg)
		this.Data["json"] = res
		this.ServeJSON()
	}()

	log.Infof("req.GroupId=%v, uid=%v, msgbody=%v", req.GroupId, this.Uid, req.MsgBody)

	//解析req参数
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &req); err != nil {
		this.Ctx.Output.SetStatus(400)
		res.Code = vipcomm.RES_INVALID_ARG
		res.Msg = err.Error()
		return
	}
	//消息入库，如果msg_id重复则插入失败
	o := orm.NewOrm()
	o.Using("app_im")
	if err := models.InsertGroupMsg(req.GroupId, this.Uid, req.MsgId, req.MsgType, req.MsgBody, o); err != nil {
		this.Ctx.Output.SetStatus(400)
		res.Code = vipcomm.RES_DB_ERROR
		res.Msg = err.Error()
		return
	}
	//从redis获取自增seq_id
	seq_id, err := redis.Int64(redisf.Do("group_seq_id", "INCR", "seqId_"+strconv.FormatInt(req.GroupId, 10)))
	if err != nil {
		this.Ctx.Output.SetStatus(400)
		res.Code = vipcomm.RES_DB_ERROR
		res.Msg = err.Error()
		return
	}

	//把seq_id入库
	if err := models.UpdateGroupMsgSeqId(req.MsgId, req.GroupId, seq_id, o); err != nil {
		this.Ctx.Output.SetStatus(400)
		res.Code = vipcomm.RES_DB_ERROR
		res.Msg = err.Error()
		return
	}

	//返回的seqid
	res.Data = &SeqIdData{SeqId: seq_id}

	//包装下发的结构体
	pack := protocol.GroupMsgData{
		GroupId: req.GroupId,
		Uid:     this.Uid,
		SeqId:   seq_id,
		MsgId:   req.MsgId,
		MsgType: req.MsgType,
		MsgBody: req.MsgBody}

	//获取群的所有用户的push_id,推送状态栏
	if err := im_utils.SysNotifySysByGrp(req.GroupId, protocol.AppMsg{"8", pack}); err != nil {
		this.Ctx.Output.SetStatus(400)
		res.Code = vipcomm.RES_UNKNOW_ERROR
		res.Msg = err.Error()
		return
	}

	//推送至app
	if err := im_utils.AppNotifyByGrp(req.GroupId, protocol.AppMsg{"8", pack}); err != nil {
		this.Ctx.Output.SetStatus(400)
		res.Code = vipcomm.RES_UNKNOW_ERROR
		res.Msg = err.Error()
		return
	}
}

func (this *GroupMsgController) PullGroupMsg() {
	res := pullGroupMsgRes{Code: vipcomm.RES_OK, Msg: "success"}
	req := pullGroupMsgReq{}

	defer func() {
		this.Data["json"] = res
		this.ServeJSON()
	}()
	//解析req参数
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &req); err != nil {
		this.Ctx.Output.SetStatus(400)
		res.Code = vipcomm.RES_INVALID_ARG
		res.Msg = err.Error()
		return
	}
	//根据seqid查询需要拉取的信息
	o := orm.NewOrm()
	o.Using("app_im")
	if data, err := models.BatchGetGroupMsg(req.GroupId, req.ToSeqId, o); err != nil {
		this.Ctx.Output.SetStatus(400)
		res.Code = vipcomm.RES_DB_ERROR
		res.Msg = err.Error()
		res.Data = data
		return
	}
}

func main() {
	// 日志
	mylog.InitLog()
	defer mylog.FlushLog()

	// 初始化服务
	if err := vipcomm.InitServer(); err != nil {
		log.Errorf("Init server faild. err:%v", err)
		return
	}

	beego.Router("/group_chat/sendGroupMsg/", &GroupMsgController{}, "post:SendGroupMsg")
	beego.Router("/group_chat/pullGroupMsg/", &GroupMsgController{}, "post:PullGroupMsg")
	beego.Run()
}
