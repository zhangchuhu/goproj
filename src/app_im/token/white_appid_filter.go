package token

import (
	"github.com/astaxie/beego"
	log "github.com/cihub/seelog"
	"vipcomm"
)

var (
	// appid 密钥
	KAppId2Key map[string]string
)

func init() {
	conf, err := beego.AppConfig.GetSection("white_appid")
	if err != nil {
		log.Errorf("get white_appid config faild")
	}

	KAppId2Key = conf
}

type WhiteAppIdController struct {
	beego.Controller
	AppId string
}

func (this *WhiteAppIdController) Prepare() {
	this.AppId = this.GetString("appid")
	if _, ok := KAppId2Key[this.AppId]; !ok {
		log.Errorf("appid:%v not in white list.", this.AppId)
		this.Data["json"] = map[string]interface{}{
			"code": vipcomm.RES_NOT_IN_WHITE,
			"msg":  "not in white list",
		}
		this.ServeJSON()
		this.StopRun()
	}
}
