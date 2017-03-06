package models

import (
	"app_im/protocol"
	"github.com/astaxie/beego/orm"
	log "github.com/cihub/seelog"
	"strconv"
)

// 群基本信息
type DGrpInfo struct {
	Id        int64  `orm:"auto"`
	Name      string `orm:"column(name)"`
	Logo      string `orm:"column(logo)"`
	Location  string `orm:"column(location)"`
	Desc      string `orm"column(desc_)"`
	Privilege int    `orm"column(privi)"`
	Password  string `orm"column(passwd)"`
	Invite    int    `orm:"column(invite)"`
}

func init() {
	orm.RegisterModel(new(DGrpInfo))
}

// 创建群Id，并更新群成员
func NewGrp(req *protocol.PCreateGrp, o orm.Ormer) (gid int64, err error) {
	sql := "insert into tbl_grp_info(privi, invite, name, logo, location, passwd, desc_, create_time, last_modify) values(?, ?, ?, ?, ?, ?, ?, now(), now())"
	res, err := o.Raw(sql, req.Privilege, req.InviteType, req.GrpName, req.Logo, req.Location, req.Password, req.Desc).Exec()
	if err == nil {
		gid, err = res.LastInsertId()
	}
	return
}

// 查询群
func GetGrpInfo(gid int64, o orm.Ormer) (ginfo DGrpInfo, err error) {
	err = o.Raw("select id, name, logo, location, desc_, privi, passwd, invite from tbl_grp_info where id=?", gid).QueryRow(&ginfo)
	return
}

func BatchGetGrpInfo(gids []int64, o orm.Ormer) (ginfos []DGrpInfo, err error) {
	lens := len(gids)
	ginfos = make([]DGrpInfo, 0, lens)
	if lens == 0 {
		return
	}

	sql := []byte("select id, name, logo, location, desc_, privi, passwd, invite from tbl_grp_info where id in (")
	for i, gid := range gids {
		sql = strconv.AppendInt(sql, gid, 10)
		if i != (lens - 1) {
			sql = append(sql, ',')
		}
	}
	sql = append(sql, ')')
	log.Info(string(sql))
	_, err = o.Raw(string(sql)).QueryRows(&ginfos)
	return
}
