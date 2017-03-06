package main

import (
	"vipcomm"
	"vipcomm/mylog"
	_ "vipcomm/mysql"

	"github.com/astaxie/beego"
	log "github.com/cihub/seelog"
)

func main() {

	// 日志
	//mylog.InitLog()
	defer mylog.FlushLog()

	// 初始化服务
	if err := vipcomm.InitServer(); err != nil {
		log.Errorf("Init server faild. err:%v", err)
		return
	}
	beego.Router("/im/bindPushId/", &ChatMsgController{}, "post:BindPushId")
	beego.Router("/im/sendChatMsg/", &ChatMsgController{}, "post:SendChatMsg")
	beego.Router("/im/getChatMsg/", &ChatMsgController{}, "post:GetChatMsg")
	beego.Run()
}
