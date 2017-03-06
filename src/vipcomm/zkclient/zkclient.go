package zkclient

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/utils"
	log "github.com/cihub/seelog"
	"github.com/samuel/go-zookeeper/zk"
	"strings"
	"time"
)

func init() {
	defer log.Flush()
}

/*
* 通过在/counters下创建节点，得到version的值，用来做轮训
 */
func ZkGetIndex(conn *zk.Conn, path string) (idx int32, err error) {
	nodeName := strings.Replace(path, "/", "_", -1)
	nodeName = strings.TrimLeft(nodeName, "_")

	//fmt.Println(nodeName)
	newPath := "/counters/" + nodeName
	var retry = 0
	var st *zk.Stat
	var exists bool
	for ; retry < 10; retry++ {
		exists, st, err = conn.Exists(newPath)
		if exists {
			// 节点已经存在
			st, err = conn.Set(newPath, []byte(""), st.Version)
			if err == nil {
				idx = st.Version
				break
			}
		} else if err == nil {
			_, err = conn.Create(newPath, []byte(""), 0, zk.WorldACL(zk.PermAll))
		}
	}

	//fmt.Println(retry, idx, err)
	return
}

func ZkClient(hosts string) (conn *zk.Conn, err error) {
	if hosts == "" {
		confPath := "/data/services/common-conf/vip_zookeeper.conf"
		if !utils.FileExists(confPath) {
			hosts = beego.AppConfig.String("zk_addr")
			log.Infof("init zkclient from conf/app.conf hosts:%v", hosts)
		} else {
			zkconf, err := config.NewConfig("ini", confPath)
			if err != nil {
				log.Errorf("parse config faild. file:%v", confPath)
				return conn, err
			}
			hosts = zkconf.String("hosts")
			log.Infof("init zkclient from %v hosts:%v", confPath, hosts)
		}
	}

	if hosts == "" {
		err = errors.New("hosts is invalid")
		return
	}

	conn, _, err = zk.Connect(strings.Split(hosts, ","), time.Second*10)
	if err != nil {
		log.Errorf("connect to zookeeper faild. err:%v", err)
		return
	}

	err = conn.AddAuth("digest", []byte("tom:tom"))
	if err != nil {
		log.Errorf("auth zookeeper faild. err:%v", err)
		return
	}
	return
}
