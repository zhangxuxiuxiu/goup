package util

import (
	"errors"
	"reflect"
	"testing"
)

type EgError struct {
	Msg error
}

func (e *EgError) Error() string {
	return "EgError:" + e.Msg.Error()
}

func TestIsNil(t *testing.T) {
	t.Logf("interface{}(nil) == nil?%v", interface{}(nil) == nil)
	t.Logf("IsNil(nil)=%v", IsNil(nil))

	t.Logf("interface{}(*int(nil)) == nil?%v", interface{}((*int)(nil)) == nil)
	t.Logf("IsNil((*int)(nil))=%v", IsNil((*int)(nil)))

	a := 2
	t.Logf("interface{}(*int(&a)) == nil?%v", interface{}(&a) == nil)
	t.Logf("IsNil(&a)=%v", IsNil(&a))

	//egError := (error)((*EgError)(nil))
	t.Logf("(error)((*EgError)(nil))==nil?%v", (error)((*EgError)(nil)) == nil)
	t.Logf("IsNil((error)((*EgError)(nil)))?%v", IsNil((error)((*EgError)(nil))))

	t.Logf("(error)(&EgError{\"test\"})==nil?%v", (error)(&EgError{errors.New("test")}) == nil)
	t.Logf("IsNil((error)(&EgError{\"test\"}))?%v", IsNil((error)(&EgError{errors.New("test")})))

	var e interface{}
	var f interface{}
	f = e
	t.Logf("IsNil(f{interface{}})?%v", IsNil(f))
	t.Logf("IsNil(error(nil)))?%v", IsNil(error(nil)))
	t.Logf("IsNil(error((*EgError)(nil)))?%v", IsNil(error((*EgError)(nil))))
	t.Logf("IsNil((*EgError)(nil))?%v", IsNil((*EgError)(nil)))
	var ege interface{} = (*EgError)(nil)
	t.Logf("IsNil((*EgError)(nil))?%v", IsNil(ege))

	//TODO IsNil fails here
	t.Logf("IsNil(([]int)(nil))?%v", IsNil(([]int)(nil)))
	t.Logf("IsNil(reflect.ValueOf(\"\").Interface())?%v", IsNil(reflect.ValueOf("").Interface()))
	t.Logf("IsNil((chan int)(nil))?%v", IsNil((chan int)(nil)))
	t.Logf("IsNil(map[int]int(nil))?%v", IsNil(map[int]int(nil)))
	t.Logf("IsNil(reflect.ValueOf(EgError{}).Field(0).Interface())?%v", IsNil(reflect.ValueOf(EgError{}).Field(0).Interface()))
	se := reflect.ValueOf(&EgError{}).Elem()
	ee := se.Field(0)
	t.Logf("kind of ee:%v", ee.Kind())
	t.Logf("IsNil(reflect.ValueOf(&EgError{}).Elem().Interface())?%v", IsNil(se.Interface()))
	t.Logf("IsNil(reflect.ValueOf(&EgError{}).Elem().Field(0).Interface())?%v", IsNil(ee.Interface()))

}

func TestIsNil2(t *testing.T) {
	var a interface{}
	t.Logf("IsNil2(a)?%v", IsNil2(a))
	t.Logf("IsNil2(([]int)(nil))?%v", IsNil2(([]int)(nil)))

}
