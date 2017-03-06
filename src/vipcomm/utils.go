package vipcomm

import (
	"reflect"
)

// 检查接口是否nil
func NilInterface(f interface{}) bool {
	return f == nil || f == reflect.Zero(reflect.TypeOf(f)).Interface()
}

// 切片去重
func UniqueInt64(s []int64) {
}
