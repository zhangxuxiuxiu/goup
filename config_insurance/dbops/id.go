package dbops

import (
	"reflect"
	"strings"
)

func IdFn(t reflect.Type) reflect.Type {
	return reflect.FuncOf([]reflect.Type{t, t}, []reflect.Type{reflect.TypeOf(true)}, false)
}

func CallId(equal reflect.Value, a, b any) bool {
	return equal.Call([]reflect.Value{reflect.ValueOf(a), reflect.ValueOf(b)})[0].Bool()
}

func Identical(a, b any) bool {
	t1, t2 := reflect.TypeOf(a), reflect.TypeOf(b)
	if t1 != t2 {
		//panic("a&b should be same type in Identical")
		return false
	}

	if isNil(a) || isNil(b) {
		return false
	}

	v1, v2 := reflect.ValueOf(a), reflect.ValueOf(b)

	if t1.Kind() == reflect.Ptr {
		t1, v1, v2 = t1.Elem(), v1.Elem(), v2.Elem()
	}

	if idFn, exist := t1.MethodByName("Identical"); exist {
		if idFn.Type == IdFn(t1) {
			return CallId(idFn.Func, a, b)
		}
	}

	if t1.Kind() != reflect.Struct {
		panic("only struct type is supported in dbops.Identical")
	}

	//	var hasId = false
	for i := 0; i < t1.NumField(); i++ {
		f := t1.Field(i)
		if hasId(f.Tag) {
			//			hasId = true
			if !EqualBasic(v1.Field(i).Interface(), v2.Field(i).Interface()) {
				return false
			}
		}
	}
	//if !hasId {
	//	panic(fmt.Sprintf("no id marked in type:%v", t1))
	//}
	return true
}

func hasId(tag reflect.StructTag) bool {
	if attr, ok := tag.Lookup("ipc"); ok {
		attrs := strings.Split(attr, ",")
		for _, att := range attrs[1:] {
			if att == "id" {
				return true
			}
		}
	}
	return false
}
