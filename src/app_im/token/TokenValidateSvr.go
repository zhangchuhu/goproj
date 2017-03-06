package token

// #cgo LDFLAGS: -L${SRCDIR} -ltoken_validate
// #cgo CPPFLAGS: -I${SRCDIR}/ycloud_token/include
// #include "TokenValidateSvr.h"
// #include <stdlib.h>
import "C"
import "errors"
import "unsafe"

var TokenService GoTokenValidateSvr

func init() {
	TokenService = New()
}

func Validate(token string) (int64, error) {
	uid, err := TokenService.Validate(token)
	return uid, err
}

type GoTokenValidateSvr struct {
	ptr C.TokenValidateSvrPtr
}

func New() GoTokenValidateSvr {
	var ret GoTokenValidateSvr
	ret.ptr = C.Init()
	return ret
}

func (f GoTokenValidateSvr) Free() {
	C.Free(f.ptr)
}

func (f GoTokenValidateSvr) Validate(token string) (int64, error) {
	cs := C.CString(token)
	defer C.free(unsafe.Pointer(cs))
	res := C.Validate(f.ptr, cs, C.int(len(token)))
	uid := int64(res)
	if uid == -1 {
		return 0, errors.New("yctoken mvartidate error")
	} else if uid == -2 {
		return 0, errors.New("yctoken expired")
	}
	return uid, nil
}

/*
func main() {
	token := "sdfsdfsdfas"
	foo := New()
	foo.Bar(token)
	foo.Free()
}
*/
