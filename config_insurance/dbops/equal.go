package dbops

import (
	"math"
	"reflect"
)

func EqualFn(t reflect.Type) reflect.Type {
	return reflect.FuncOf([]reflect.Type{t, t}, []reflect.Type{reflect.TypeOf(true)}, false)
}

func CallEqual(equal reflect.Value, a, b interface{}) bool {
	return equal.Call([]reflect.Value{reflect.ValueOf(a), reflect.ValueOf(b)})[0].Bool()
}

// Equal compare basic type field in struct
func Equal(a, b interface{}) bool {
	ta, tb := reflect.TypeOf(a), reflect.TypeOf(b)
	if ta != tb {
		return false
	}
	if ta.Kind() == reflect.Ptr {
		ta, tb = ta.Elem(), tb.Elem()
	}

	if equalFn, exist := ta.MethodByName("Equal"); exist {
		if equalFn.Type == EqualFn(ta) {
			return CallEqual(equalFn.Func, a, b)
		}
	}

	if ta.Kind() != reflect.Struct {
		panic("only struct is supported in dbops.Equal")
	}

	// only compare fields on first level
	va, vb := reflect.ValueOf(a), reflect.ValueOf(b)
	if va.Kind() == reflect.Ptr {
		va, vb = va.Elem(), vb.Elem()
	}
	//allBasics := true
	for i := 0; i < ta.NumField(); i++ {
		f := ta.Field(i)
		if f.Tag.Get("equal") == "-" {
			continue
		}
		fa, fb := va.Field(i), vb.Field(i)
		if !basicType(f.Type.Kind()) {
			//allBasics = false
		} else if !EqualBasic(fa.Interface(), fb.Interface()) {
			return false
		}
	}

	return true
}

func EqualBasic(a, b interface{}) bool {
	ta, tb := reflect.TypeOf(a), reflect.TypeOf(b)
	if ta != tb {
		panic("a&b should be same type in EqualBasic")
	}

	fa, fb := reflect.ValueOf(a), reflect.ValueOf(b)
	v := false
	switch ta.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v = fa.Int() == fb.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		v = fa.Uint() == fb.Uint()
	case reflect.String:
		v = fa.String() == fb.String()
	case reflect.Bool:
		v = fa.Bool() == fb.Bool()
	case reflect.Float32, reflect.Float64:
		v = math.Abs(fa.Float()-fb.Float()) < 1e-8
	case reflect.Complex64, reflect.Complex128:
		v = fa.Complex() == fb.Complex()
	default:
		panic("only basic type field is supported")
	}
	return v
}
