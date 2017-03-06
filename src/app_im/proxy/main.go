package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	_ "reflect"
	"vipcomm/mylog"

	log "github.com/cihub/seelog"
	"github.com/googollee/go-socket.io"
)

type PushId struct {
	Id  string `json:"id"`
	Uid string `json:"uid"`
}
type Message struct {
	Msg      string `json:"message"`
	NickName string `json:"nickName"`
	Type     string `json:"type"`
}
type ChatRes struct {
	Topic string `json:"topic"`
	Data  string `json:"data"`
}

func main() {
	// 日志
	//mylog.InitLog("proxy", "", "")
	defer mylog.FlushLog()

	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Errorf("Init socket.io server faild. err:%v", err)
	}
	conns := make(map[string]socketio.Socket)
	server.On("connection", func(so socketio.Socket) {
		log.Infof("on connection")
		so.Join("chatRoom")
		so.On("pushId", func(general interface{}) {
			data := general.(map[string]interface{})
			if id, ok := data["id"].(string); ok {
				conns[id] = so
			}
			pk := PushId{"vxh5v5KuKSy1234", "123456"}
			err := so.Emit("pushId", pk)
			if err != nil {
				log.Errorf("Send pushId to client faild. err:%v", err)
			}
		})

		so.On("chat_message", func(general interface{}) {
			m := general.(map[string]interface{})
			fmt.Println("********************", m["message"], m["nickName"], "******************")
			message := Message{}
			message.Type = "chat_message"

			if msg, ok := m["message"].(string); ok {
				message.Msg = msg
			}
			if nickName, ok := m["nickName"].(string); ok {
				message.NickName = nickName
			}
			data, _ := json.Marshal(message)
			var res ChatRes
			res.Topic = "chatRoom"
			res.Data = base64.StdEncoding.EncodeToString(data)
			//err := so.Emit("push", res)
			server.BroadcastTo("chatRoom", "push", res)
			//			if err != nil {
			//				fmt.Println("###########", err, "##############")
			//			}
			fmt.Println("********************end******************")
		})
		so.On("disconnection", func() {
			log.Infof("Send pushId to client faild. err:%v", err)
		})

	})
	server.On("error", func(so socketio.Socket, err error) {
		fmt.Println("error:", err)
	})
	//	server.On("message", func(msg string) {
	//		fmt.Println("*********aaaaaaaaaa****************\n\n")
	//	})

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Infof("Serving at localhost:5000...")
	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Errorf("ListenAndServe happen error:%v", err)
	}

}
