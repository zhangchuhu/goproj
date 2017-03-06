package vipcomm

import (
	//"app_im/protocol"
	"encoding/json"
	"fmt"
	"testing"
	//"time"
	"vipcomm/mylog"
)

type Packet struct {
	Typ int32
	Msg string
}

func Test_Push(t *testing.T) {
	defer mylog.FlushLog()
	for true {
		url := "http://106.38.255.199:11001/api/push"
		pushIds := []string{"vxh5v5KuKSy1234", "vip_app_im_saf"}
		pk := Packet{1, "Hellow"}

		data, _ := json.Marshal(pk)

		err := PushClient(url, pushIds, string(data))
		fmt.Println(err)
		//time.Sleep(time.Second)
	}
}

/*
func TestNotify(t *testing.T) {
	defer mylog.FlushLog()
	pushIds := []string{"vxh5v5KuKSy1234", "vxh5v5KuKSy1235"}
	android := protocol.Android{"title", "message", protocol.PChatMsgNotify{1, 1}}
	apn := protocol.Apn{"message", 1, "default", protocol.PChatMsgNotify{1, 1}}

	var msg protocol.SysMsg
	msg.JAndroid = android
	msg.JApn = apn
	err := SysNotify(pushIds, msg)
	fmt.Println(err)
}
*/
