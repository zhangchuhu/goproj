package main

import (
	"github.com/astaxie/beego"
	log "github.com/cihub/seelog"
	"vipcomm"
	"vipcomm/mylog"
)

func main() {
	// 日志
	mylog.InitLog()
	defer mylog.FlushLog()

	// 初始化服务
	if err := vipcomm.InitServer(); err != nil {
		log.Errorf("Init server faild. err:%v", err)
		return
	}

	// 查询群身份
	beego.Router("/back_services/get_grp_role/", &BackGrpController{}, "get:GetGrpRole")

	// 查询群信息
	beego.Router("/back_services/get_grp_info/", &BackGrpController{}, "get:GetGrpInfo")

	beego.Run()
}
