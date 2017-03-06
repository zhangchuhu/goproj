package vipcomm

import (
	"fmt"
	"github.com/astaxie/beego/httplib"
	"net/http"
	"time"
)

func init() {
	defaultTrans := httplib.BeegoHTTPSettings{
		UserAgent:        "vip-server",
		ConnectTimeout:   10 * time.Second,
		ReadWriteTimeout: 20 * time.Second,
		Gzip:             true,
		DumpBody:         true,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 2, // http持久化连接
		},
	}
	httplib.SetDefaultSetting(defaultTrans)
	fmt.Println("http_default_seting")
}
