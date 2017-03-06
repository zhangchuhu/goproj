package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"strconv"
)

// 群成员
type DGrpMember struct {
	Id   int64  `orm:"auto"`
	Gid  int64  `orm:"column(gid_)"`
	Uid  int64  `orm:"column(uid_)"`
	Role int    `orm:"column(role_)"` // 角色
	Flag int    `orm:"column(flag_)"` // 0 退出群 1 在群里
	Nick string `orm:"column(nick_)"` // 群昵称
}

func init() {
	orm.RegisterModel(new(DGrpMember))
}

func genGmName(gid int64) string {
	return fmt.Sprintf("tbl_grp_member_%d", gid%1)
}

/*
* 设置群成员
* @param op    1 加入群， 0 退出群
* @param role  0 普通成员、1 群管理、 2 群主
 */
func SetGrpMember(uid, gid int64, op, role int, o orm.Ormer) (err error) {
	sql := "insert into " + genGmName(gid) + "(gid_, uid_, role_, flag_, last_modify) values(?, ?, ?, ?, now()) on duplicate key update role_=?, flag_=?, last_modify=now()"
	_, err = o.Raw(sql, gid, uid, role, op, role, op).Exec()
	return
}

/*
* 加入/退出群
 */
func UpdateGrpMemberFlag(uid, gid int64, flag int, o orm.Ormer) (err error) {
	_, err = o.Raw("update "+genGmName(gid)+" set flag_=? where gid_=? and uid_=?", flag, gid, uid).Exec()
	return
}

/*
* 根据uid获取群昵称、角色等信息
 */
func GetGrpMemberByUid(uid, gid int64, o orm.Ormer) (info DGrpMember, err error) {
	err = o.Raw("select uid_, gid_, role_, flag_, nick_ from "+genGmName(gid)+" where gid_=? and uid_=?", gid, uid).QueryRow(&info)
	return
}

/*
* 获取群管理（包括群主)
 */
func GetGrpManagers(gid int64, o orm.Ormer) (ulist []int64, err error) {
	var res orm.ParamsList
	n, err := o.Raw("select uid_ from "+genGmName(gid)+" where gid_ = ? and  role_ > 0 and flag_ = 1", gid).ValuesFlat(&res)
	if err == nil && n > 0 {
		for _, i := range res {
			uid, _ := strconv.ParseInt(i.(string), 10, 64)
			ulist = append(ulist, uid)
		}
	}
	return
}

/*
* 获取群成员
 */
func GetGrpMember(gid int64, o orm.Ormer) (ulist []int64, err error) {
	var res orm.ParamsList
	n, err := o.Raw("select uid_ from "+genGmName(gid)+" where gid_ = ? and flag_ = 1", gid).ValuesFlat(&res)
	if err == nil && n > 0 {
		for _, i := range res {
			uid, _ := strconv.ParseInt(i.(string), 10, 64)
			ulist = append(ulist, uid)
		}
	}
	return
}

/*
* 获取群人数
 */
func GetGrpMemberNumber(gid int64, o orm.Ormer) (num int, err error) {
	var res orm.ParamsList
	n, err := o.Raw("select count(1) from "+genGmName(gid)+" where gid_ = ? and flag_ = 1", gid).ValuesFlat(&res)
	if err == nil && n > 0 {
		num, _ = strconv.Atoi(res[0].(string))
	}
	return

}
