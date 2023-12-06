package dbops

import (
	"fmt"
	"reflect"
	"strings"
	"xorm.io/xorm/names"
)

func basicType(k reflect.Kind) bool {
	return k <= reflect.Complex128 || k == reflect.String
	//switch k {
	//case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.String, reflect.Bool, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
	//	return true
	//default:
	//	return false
	//}
}

var snake = names.SnakeMapper{}.Obj2Table

func TableName(f reflect.StructField) string {
	if attr := f.Tag.Get("ipc"); attr != "" {
		attrs := strings.Split(attr, ",")
		if attrs[0] != "" {
			return fmt.Sprintf("`%s`", attrs[0])
		}
	}
	return fmt.Sprintf("`%s`", snake(f.Name))
}

var ColumnName = TableName

func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Func, reflect.Chan, reflect.Slice, reflect.Interface, reflect.UnsafePointer:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}

func Update(p, p2 interface{}, tableName string) string {
	// ensure the same entity
	if !Identical(p, p2) {
		sql := ""
		if !isNil(p) {
			sql += Delete(p, tableName)
		}
		if !isNil(p2) {
			if len(sql) > 0 {
				sql += "\n"
			}
			sql += Insert(p2, tableName)
		}
		return sql
	}

	t, v1, v2 := reflect.TypeOf(p), reflect.ValueOf(p), reflect.ValueOf(p2)
	if t.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("only pointer value is supported in dbops.Update,{error type kind:%v}", t.Kind()))
	}
	t, v1, v2 = t.Elem(), v1.Elem(), v2.Elem()
	if t.Kind() != reflect.Struct {
		panic(fmt.Sprintf("only struct is supported in Update(interface{},interface{},tableName),{error type:%v}", t.Kind()))
	}

	var idCols []string
	var idNames []string
	var updateCols []string
	var updateNames []string
	var fieldsSql []string
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Tag.Get("ipc") == "-" {
			continue
		}
		switch f.Type.Kind() {
		case reflect.Struct, reflect.Interface, reflect.Ptr:
			v1 := v1.Field(i)
			v2 := v2.Field(i)
			if f.Type.Kind() == reflect.Struct && v1.CanAddr() {
				v1 = v1.Addr()
				v2 = v2.Addr()
			}
			fieldsSql = append(fieldsSql, Update(v1.Interface(), v2.Interface(), TableName(t.Field(i))))
		case reflect.Slice:
			fieldsSql = append(fieldsSql, updateSlice(v1.Field(i), v2.Field(i), TableName(t.Field(i))))
		default:
			if hasId(f.Tag) {
				idCols = append(idCols, stringify(v2.Field(i).Interface()))
				idNames = append(idNames, ColumnName(f))
				continue
			}
			if basicType(f.Type.Kind()) {
				if !EqualBasic(v1.Field(i).Interface(), v2.Field(i).Interface()) {
					updateCols = append(updateCols, stringify(v2.Field(i).Interface()))
					updateNames = append(updateNames, ColumnName(f))
				}
			} else {
				panic(fmt.Sprintf("unsupported type in Update:%v", f.Type.Kind()))
			}
		}
	}

	if len(updateCols) > 0 {
		if len(idCols) == 0 {
			panic(fmt.Sprintf("no id tag in type:%v", p))
		}

		recordId := fmt.Sprintf("%s_%s", tableName, strings.Join(idCols, "_"))
		if old, exist := uniqRecord[recordId]; !exist {
			uniqRecord[recordId] = v2
		} else if !Equal(old, v2) {
			panic(fmt.Sprintf("two differrent values with same id=>old{%#v};new{%#v}", old, v2))
		} else {
			return ""
		}

		var buffer strings.Builder
		buffer.WriteString("update ")
		buffer.WriteString(tableName)
		buffer.WriteString(" set ")
		for i := 0; i < len(updateCols); i++ {
			if i != 0 {
				buffer.WriteByte(',')
			}
			buffer.WriteString(fmt.Sprintf("%s=%v", updateNames[i], updateCols[i]))
		}
		buffer.WriteString(" where ")
		for i := 0; i < len(idCols); i++ {
			if i != 0 {
				buffer.WriteString(" and ")
			}
			buffer.WriteString(fmt.Sprintf("%s=%v", idNames[i], idCols[i]))
		}
		buffer.WriteByte(';')
		fieldsSql = append(fieldsSql, buffer.String())
	}
	return strings.Join(trimLines(fieldsSql), "\n")
}

func trimLines(lines []string) []string {
	var sentences []string
	for _, line := range lines {
		if len(line) != 0 {
			sentences = append(sentences, line)
		}
	}
	return sentences
}

func find(slice reflect.Value, v interface{}) int {
	for i := 0; i < slice.Len(); i++ {
		up := slice.Index(i).Interface()
		if Identical(up, v) {
			return i
		}
	}
	return -1
}

// TODO customization point for Update
//func UpdateToFn(t reflect.Type) reflect.Type {
//	return reflect.FuncOf([]reflect.Type{t, t}, []reflect.Type{reflect.TypeOf("")}, false)
//}
//
//func CallUpdateTo(updateTo reflect.Value, a, b interface{}) string {
//	return updateTo.Call([]reflect.Value{reflect.ValueOf(a), reflect.ValueOf(b)})[0].String()
//}
//
//// updateTo p&p2 of the same type and has method func (*T)UpdateTo(*T) string
//func updateTo(p, p2 interface{}) string {
//	t := reflect.TypeOf(p)
//	if updateToFn, exist := t.MethodByName("UpdateTo"); exist {
//		if updateToFn.Type == UpdateToFn(t) {
//			return CallUpdateTo(updateToFn.Func, p, p2)
//		}
//	}
//	panic("p&p2 should be of the same type and have method func (*T)UpdateTo(*T) string")
//}

// updateSlice typeOf v1&v2 must both be slice
func updateSlice(v1, v2 reflect.Value, tableName string) string {
	var updates []string
	for i := 0; i < v1.Len(); i++ {
		p := v1.Index(i).Interface()
		i2 := find(v2, p)
		if i2 == -1 {
			updates = append(updates, Delete(v1.Index(i).Addr().Interface(), tableName))
		} else { // update the new in v2
			updates = append(updates, Update(v1.Index(i).Addr().Interface(), v2.Index(i2).Addr().Interface(), tableName))
		}
	}

	// insert the new in v2
	for i := 0; i < v2.Len(); i++ {
		p2 := v2.Index(i).Interface()
		if -1 == find(v1, p2) {
			updates = append(updates, Insert(v2.Index(i).Addr().Interface(), tableName))
		}
	}
	return strings.Join(trimLines(updates), "\n")
}
