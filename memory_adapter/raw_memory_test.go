package memory_adapter

import (
	"fmt"
	"testing"
	"unsafe"
)

type Trivial struct {
	Id   uint64
	Name [64]byte
}

func (obj Trivial) Hello() {
	obj.Id = 100
}

func TestTrivial(*testing.T) {
	triSize := unsafe.Sizeof(Trivial{})
	bytes := make([]byte, triSize*2)

	bytes[0] = '1'
	tri := (*Trivial)(unsafe.Pointer(&bytes[0]))
	fmt.Printf("id:%d\n", tri.Id)

	bytes[triSize] = '2'
	tri2 := (*Trivial)(unsafe.Pointer(&bytes[triSize]))
	fmt.Printf("id:%d\n", tri2.Id)

	tri3 := (*Trivial)(unsafe.Pointer(uintptr(unsafe.Pointer(&bytes[0])) + triSize))
	fmt.Printf("id:%d\n", tri3.Id)

	var obj Trivial
	obj.Hello()
	fmt.Printf("id of trivial:%d\n", obj.Id)
	(&obj).Hello()
	fmt.Printf("id of trivial:%d\n", obj.Id)

}
