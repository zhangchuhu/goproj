package token

import (
	"github.com/astaxie/beego"
	log "github.com/cihub/seelog"
	"vipcomm"
)

type BaseController struct {
	beego.Controller
	Uid int64
}

func (this *BaseController) Prepare() {
	tok := this.Ctx.Input.Header("token")
	log.Debugf("token:%v", tok)
	uid, err := Validate(tok)
	if err != nil {
		log.Errorf("parse token faild. token:%v", tok)
		this.Data["json"] = map[string]interface{}{"code": vipcomm.RES_TOKEN_VALIDATE_ERROR, "msg": "token parse error"}
		this.ServeJSON()
		this.StopRun()
	}
	this.Uid = uid
}
