package main

import (
	"app_im/protocol"
	log "github.com/cihub/seelog"
	"net/http"
	"strconv"
	"vipcomm"
)

func BindPushId(w http.ResponseWriter, r *http.Request) {
	res := protocol.HttpCommRes{
		Code: vipcomm.RES_OK,
	}
	r.ParseForm()
	uid, err := strconv.ParseInt(r.Form.Get("uid"), 10, 64)
	pushid := r.Form.Get("pushid")
	if err != nil {
		res.Code = vipcomm.RES_INVALID_ARG
		vipcomm.HttpJsonRes(res, w)
		log.Errorf("faild. invalid argument uid:%v", r.Form.Get("uid"))
		return
	}

	vipcomm.HttpJsonRes(res, w)
	log.Infof("ok. uid:%v pushId:%v", uid, pushid)
}
