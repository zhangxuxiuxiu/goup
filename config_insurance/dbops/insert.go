package dbops

import (
	"fmt"
	"reflect"
	"strings"
)

type GenInsert interface {
	GenInsertSql() string
}

var uniqRecord = map[string]any{}

func Insert(e any, tableName string) string {
	if isNil(e) {
		return ""
	}

	t := reflect.TypeOf(e)
	v := reflect.ValueOf(e)

	if gen, ok := v.Interface().(GenInsert); ok {
		return gen.GenInsertSql()
	} else if v.CanAddr() {
		if gen, ok := v.Addr().Interface().(GenInsert); ok {
			return gen.GenInsertSql()
		}
	} else if t.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("only pointer value is supported in dbops.Insert,{error type kind:%v}", t.Kind()))
	}

	if t.Kind() == reflect.Ptr {
		t, v = t.Elem(), v.Elem()
	}
	if t.Kind() != reflect.Struct {
		panic(fmt.Sprintf("only struct is supported in Insert(any,tableName),{error type:%v}", t.Kind()))
	}

	var insertNames []string
	var insertVals []string
	var idVals []string
	var subSqls []string
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		name := ColumnName(f)
		if name == "`-`" {
			continue
		}

		realType := v.Field(i).Type().Kind()
		if basicType(realType) {
			insertNames = append(insertNames, name)
			insertVals = append(insertVals, stringify(v.Field(i).Interface()))
			if hasId(f.Tag) {
				idVals = append(idVals, stringify(v.Field(i).Interface()))
			}
		} else if realType == reflect.Struct || realType == reflect.Interface || realType == reflect.Ptr {
			vi := v.Field(i)
			if realType == reflect.Struct && vi.CanAddr() {
				vi = vi.Addr()
			}
			subSqls = append(subSqls, Insert(vi.Interface(), TableName(f)))
		} else if realType == reflect.Slice {
			for j := 0; j < v.Field(i).Len(); j++ {
				subSqls = append(subSqls, Insert(v.Field(i).Index(j).Addr().Interface(), TableName(f)))
			}
		} else {
			panic(fmt.Sprintf("unsupported field type:%v", realType))
		}
	}

	if len(insertVals) > 0 {
		if len(idVals) > 0 {
			recordId := fmt.Sprintf("%s_%s", tableName, strings.Join(idVals, "_"))
			//		fmt.Printf("current record id:%s\n", recordId)
			if old, exist := uniqRecord[recordId]; !exist {
				uniqRecord[recordId] = e
			} else if !Equal(old, e) {
				panic(fmt.Sprintf("two differrent values with same id=>old{%#v};new{%#v}", old, e))
			} else {
				return ""
			}
		}
		baseSql := fmt.Sprintf("insert into %s (%s) values (%s); ", tableName, strings.Join(insertNames, ","), strings.Join(insertVals, ","))
		subSqls = append(subSqls, baseSql)
	}
	return strings.Join(trimLines(subSqls), "\n")
}

func stringify(v any) string {
	if reflect.TypeOf(v).Kind() == reflect.String {
		return fmt.Sprintf("'%v'", v)
	} else {
		return fmt.Sprintf("%v", v)
	}
}
