package main

import (
	"fmt"
	"github.com/astaxie/beego/httplib"
	"net/http"
	"time"
)

func main() {
	defaultTrans := httplib.BeegoHTTPSettings{
		UserAgent:        "beegoServer",
		ConnectTimeout:   10 * time.Second,
		ReadWriteTimeout: 20 * time.Second,
		Gzip:             true,
		DumpBody:         true,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 2,
		},
	}
	httplib.SetDefaultSetting(defaultTrans)

	for {
		req := httplib.Post("http://106.38.255.199:9098/buddy/full_sync/")
		res, err := req.String()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(res)
		}
		time.Sleep(1 * time.Second)
	}

}
