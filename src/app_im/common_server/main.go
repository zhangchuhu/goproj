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
	beego.Router("/im/bindPushId/", &CommonController{}, "post:BindPushId")
	beego.Router("/im/groupBroast/", &CommonController{}, "post:GroupBroast")
	beego.Router("/im/nearbyGroup/", &CommonController{}, "post:NearbyGroup")
	beego.Router("/im/groupGeo/", &CommonController{}, "post:GroupGeo")
	beego.Run()
}
