package main

import (
	"app_im/models"
	"app_im/protocol"
	"app_im/token"
	"github.com/astaxie/beego/orm"
	log "github.com/cihub/seelog"
	"vipcomm"
)

type BackGrpController struct {
	token.WhiteAppIdController
}

// 查询群身份
func (this *BackGrpController) GetGrpRole() {
	var (
		res     protocol.PGrpMemberInfo
		err     error
		errStr  = "success"
		resCode = vipcomm.RES_OK
		appId   = this.GetString("appid")
		uid, _  = this.GetInt64("uid")
		gid, _  = this.GetInt64("gid")
	)
	defer func() {
		log.Infof("%v res:%v. appId:%v uid:%v gid:%v role:%v err:%v",
			errStr, resCode, appId, uid, gid, res.Role, err)
		this.Data["json"] = protocol.HttpRes{
			Code: resCode,
			Data: res,
			Msg:  errStr,
		}
		this.ServeJSON()
	}()

	// read db
	o := orm.NewOrm()
	o.Using("app_im")
	info, err := models.GetGrpMemberByUid(uid, gid, o)
	if err != nil {
		resCode = vipcomm.RES_DB_ERROR
		errStr = "get user info faild"
		return
	}

	res.Role = info.Role
	return
}

// 查询群信息
func (this *BackGrpController) GetGrpInfo() {
	var (
		res     protocol.PGrpInfo
		err     error
		errStr  = "success"
		resCode = vipcomm.RES_OK
		appId   = this.GetString("appid")
		gid, _  = this.GetInt64("gid")
	)
	defer func() {
		log.Infof("%v res:%v. appId:%v gid:%v name:%v logo:%v err:%v",
			errStr, resCode, appId, gid, res.Name, res.Logo, err)
		this.Data["json"] = protocol.HttpRes{
			Code: resCode,
			Data: res,
			Msg:  errStr,
		}
		this.ServeJSON()
	}()

	// read db
	o := orm.NewOrm()
	o.Using("app_im")
	info, err := models.GetGrpInfo(gid, o)
	if err != nil {
		resCode = vipcomm.RES_DB_ERROR
		errStr = "get user info faild"
		return
	}

	res.Gid = info.Id
	res.Logo = info.Logo
	res.Name = info.Name
	res.Location = info.Location
	res.Desc = info.Desc
	res.Privilege = info.Privilege
	res.InviteType = info.Invite
	if info.Password != "" {
		res.Password = 1
	}
	return
}
