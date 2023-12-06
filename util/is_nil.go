package util

import (
	"reflect"
	"unsafe"
)

// https://codefibershq.com/blog/golang-why-nil-is-not-always-nil
// 位kindDirectIface=1<<5 用来标识iface,eface中data保存的是指针，还是值本身
// https://www.jianshu.com/p/213ef3b3a2b8
// IsNil wrong implementation, see tests
func IsNil(val interface{}) bool {
	//return ((*struct {
	//	_   *int
	//	ptr unsafe.Pointer
	//})(unsafe.Pointer(&val))).ptr == nil
	return (*[2]uintptr)(unsafe.Pointer(&val))[1] == 0
}

func IsNil2(val interface{}) bool {
	return val == nil || reflect.ValueOf(val).IsNil()
}
