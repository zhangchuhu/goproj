package main

import (
	"app_im/im_utils"
	"app_im/models"
	"app_im/protocol"
	"encoding/json"
	_ "fmt"

	"vipcomm/redisf"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	log "github.com/cihub/seelog"
)

type CommonController struct {
	beego.Controller
}

// Examples:
//
//   req: POST /im/BindPushId/ {"Uid":"1","PushId":"test"}
//   res: 200                  { "Code": "0", "Msg": "" }

func (this *CommonController) BindPushId() {

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

func (this *CommonController) GroupBroast() {
	var res protocol.HttpRes
	defer func() {
		this.Data["json"] = res
		this.ServeJSON()
	}()

	req := struct {
		Gid     int64       `json:"gid"`
		MsgType int32       `json:"msgType"`
		Extra   interface{} `json:"extra"`
	}{}

	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &req); err != nil {
		this.Ctx.Output.SetStatus(400)
		res.Code = "-1"
		res.Msg = err.Error()
		return
	}
	im_utils.AppNotifyByGrp(req.Gid, req)
	res.Code = "0"
	res.Msg = "sucess"
	return
}

func (this *CommonController) NearbyGroup() {
	var res protocol.HttpRes
	defer func() {
		this.Data["json"] = res
		this.ServeJSON()
	}()

	req := struct {
		Lng    float64 `json:"lng"`
		Lat    float64 `json:"lat"`
		Radius int32   `json:"radius"`
	}{}

	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &req); err != nil {
		this.Ctx.Output.SetStatus(400)
		res.Code = "-1"
		res.Msg = err.Error()
		return
	}
	res.Code = "0"
	res.Msg = "sucess"
	return
}

func (this *CommonController) GroupGeo() {
	var res protocol.HttpRes
	defer func() {
		this.Data["json"] = res
		this.ServeJSON()
	}()

	req := struct {
		Lng float64 `json:"lng"`
		Lat float64 `json:"lat"`
		Gid int32   `json:"gid"`
	}{}

	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &req); err != nil {
		this.Ctx.Output.SetStatus(400)
		res.Code = "-1"
		res.Msg = err.Error()
		return
	}
	key := "grpgeo"
	_, err := redisf.ClusterDo("nearby", "GEOADD", key, req.Lng, req.Lat, req.Gid)
	if err != nil {
		log.Infof("-set %s: %s\n", key, err.Error())
	} else {
		log.Infof("success\n")
	}
	res.Code = "0"
	res.Msg = "sucess"
	return
}
