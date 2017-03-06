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

	// 添加好友
	beego.Router("/buddy/addbuddy/", &BuddyController{}, "post:AddBuddy")

	// 同意/拒绝添加好友
	beego.Router("/buddy/addbuddy_answer/", &BuddyController{}, "post:AddBuddyAnswer")

	// 删除好友
	beego.Router("/buddy/del_buddy/", &BuddyController{}, "post:DelBuddy")

	// 同步好友列表
	beego.Router("/buddy/full_sync/", &BuddyController{}, "post:FullSyncBuddy")

	// 创建群
	beego.Router("/group/create_grp/", &GroupController{}, "post:CreateGrp")

	// 加入群
	beego.Router("/group/join_grp/", &GroupController{}, "post:JoinGrp")

	// 退出群
	beego.Router("/group/exit_grp/", &GroupController{}, "post:LeaveGrp")

	// 群邀请
	beego.Router("/group/invite/", &GroupController{}, "post:InviteJoinGrp")

	// 同步群成员列表
	beego.Router("/group/sync_member/", &GroupController{}, "post:SyncGrpMembers")

	// 全量同步群列表
	beego.Router("/group/full_sync/", &GroupController{}, "post:FullSyncGrpList")

	// 管理员验证入群
	beego.Router("/group/check_join_grp/", &GroupController{}, "post:CheckJoinGrp")
	beego.Run()
}
