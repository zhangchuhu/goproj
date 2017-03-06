package vipcomm

import (
	"encoding/base64"
	"encoding/json"
	"errors"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/bitly/go-simplejson"
	log "github.com/cihub/seelog"
)

var (
	pushClientUrl = "http://106.38.255.199:11001/api/push"
	sysNotifyUrl  = "http://106.38.255.199:11001/api/notification"
)

func init() {
	if push_address := beego.AppConfig.String("pushAddress"); push_address != "" {
		SetPushAddress(push_address)
	}
}

func SetPushAddress(addr string) {
	if addr != "" {
		pushClientUrl = addr + "/api/push"
		sysNotifyUrl = addr + "/api/notification"
		log.Infof("=== set push addr to %v", addr)
	}
}

// 发送App内通知
func PushClient(url string, pushIds []string, msg string) error {
	if url == "" {
		url = pushClientUrl
	}

	req := httplib.Post(url)
	ids, _ := json.Marshal(pushIds)
	req.Param("pushId", string(ids))
	req.Param("data", base64.StdEncoding.EncodeToString([]byte(msg)))

	res, err := req.String()
	if err != nil {
		return err
	}

	js, err := simplejson.NewJson([]byte(res))
	log.Infof("[PushClientPack] pushIds:%v, res:%v", pushIds, js)
	rescode, _ := js.Get("code").String()
	if rescode != "success" {
		return errors.New("Invalid Rescode " + rescode)
	}
	return nil
}

func PushClientPack(pushIds []string, req interface{}) error {
	log.Infof("[PushClientPack]********** pushIds:%v,url:%v*************", pushIds, pushClientUrl)
	cli := httplib.Post(pushClientUrl)
	ids, _ := json.Marshal(pushIds)
	cli.Param("pushId", string(ids))

	data, _ := json.Marshal(req)
	cli.Param("data", base64.StdEncoding.EncodeToString(data))

	res, err := cli.String()
	if err != nil {
		return err
	}

	js, err := simplejson.NewJson([]byte(res))
	log.Infof("[PushClientPack] pushIds:%v, res:%v", pushIds, js)
	rescode, _ := js.Get("code").String()
	if rescode != "success" {
		return errors.New("Invalid Rescode " + rescode)
	}
	return nil
}

// 发送状态栏通知

func SysNotify(pushIds []string, req interface{}) error {
	httpClient := httplib.Post(sysNotifyUrl)
	ids, _ := json.Marshal(pushIds)
	httpClient.Param("pushId", string(ids))

	data, _ := json.Marshal(req)
	httpClient.Param("notification", string(data))

	res, err := httpClient.String()
	if err != nil {
		return err
	}

	js, err := simplejson.NewJson([]byte(res))
	log.Infof("[SysNotify] pushIds:%v, res:%v", pushIds, js)

	rescode, _ := js.Get("code").String()
	if rescode != "success" {
		return errors.New("Invalid Rescode " + rescode)
	}
	return nil
}
