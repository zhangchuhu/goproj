package vipcomm

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	RES_OK                   = "0"
	RES_DB_ERROR             = "-1"
	RES_INVALID_ARG          = "-2"
	RES_UNKNOW_ERROR         = "-3"
	RES_TOKEN_VALIDATE_ERROR = "-4"
	RES_NOT_IN_WHITE         = "-5"

	RES_NOT_FOUND        = "-101" // 未找到群Id或其他
	RES_REPEATE          = "-102" // 重复操作
	RES_INVALID_PASSWORD = "-103" // 密码错误
	RES_OVER_LIMIT       = "-104" // 人数超过限制
	RES_INVALID_OP       = "-105" // 非法操作
)

func HttpJsonRes(res interface{}, w http.ResponseWriter) (n int, err error) {
	js, _ := json.Marshal(res)
	return fmt.Fprint(w, string(js))
}
