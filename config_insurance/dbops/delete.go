package dbops

import (
	"fmt"
	"reflect"
	"strings"
)

type GenDelete interface {
	GenDeleteSql() string
}

func Delete(e any, tableName string) string {
	if isNil(e) {
		return ""
	}

	t := reflect.TypeOf(e)
	v := reflect.ValueOf(e)
	if gen, ok := v.Interface().(GenDelete); ok {
		return gen.GenDeleteSql()
	} else if v.CanAddr() {
		if gen, ok := v.Addr().Interface().(GenDelete); ok {
			return gen.GenDeleteSql()
		}
	} else if t.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("only pointer value is supported in dbops.Delete,{error type kind:%v}", t.Kind()))
	}

	if t.Kind() == reflect.Ptr {
		t, v = t.Elem(), v.Elem()
	}
	if t.Kind() != reflect.Struct {
		panic(fmt.Sprintf("only struct is supported in Delete(any,tableName),{error type:%v}", t.Kind()))
	}

	var idCols []string
	var idNames []string
	var subSqls []string
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if hasId(f.Tag) {
			if basicType(f.Type.Kind()) {
				idCols = append(idCols, stringify(v.Field(i).Interface()))
				idNames = append(idNames, ColumnName(f))
			} else {
				panic(fmt.Sprintf("only basic type supported with 'id' tag"))
			}
		}

		if f.Type.Kind() == reflect.Struct || f.Type.Kind() == reflect.Interface || f.Type.Kind() == reflect.Ptr {
			vi := v.Field(i)
			if f.Type.Kind() == reflect.Struct && vi.CanAddr() {
				vi = vi.Addr()
			}
			subSqls = append(subSqls, Delete(vi.Interface(), TableName(f)))
		} else if f.Type.Kind() == reflect.Slice {
			for j := 0; j < v.Field(i).Len(); j++ {
				subSqls = append(subSqls, Delete(v.Field(i).Index(j).Addr().Interface(), TableName(f)))
			}
		} else if !basicType(f.Type.Kind()) {
			panic(fmt.Sprintf("unsupported field type in Delete:%v", f.Type))
		}
	}

	if len(idCols) > 0 {
		recordId := fmt.Sprintf("%s_%s", tableName, strings.Join(idCols, "_"))
		if old, exist := uniqRecord[recordId]; !exist {
			var buffer strings.Builder
			buffer.WriteString("update ")
			buffer.WriteString(tableName)
			buffer.WriteString("  set is_deleted=1 where ")
			for j := 0; j < len(idCols); j++ {
				if j != 0 {
					buffer.WriteString(" and ")
				}
				buffer.WriteString(fmt.Sprintf("%s=%v", idNames[j], idCols[j]))
			}
			buffer.WriteByte(';')
			subSqls = append(subSqls, buffer.String())

			uniqRecord[recordId] = e
		} else if !Equal(old, e) {
			panic(fmt.Sprintf("two differrent values with same id=>old{%#v};new{%#v}", old, e))
		} else {
			return ""
		}
	}
	return strings.Join(trimLines(subSqls), "\n")
}
