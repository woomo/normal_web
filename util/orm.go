package util

import (
	"reflect"
	"strings"
)

func GetGormFields(stc any) []string {
	typ := reflect.TypeOf(stc)
	if typ.Kind() == reflect.Ptr { //如果传的是指针类型，先解析指针
		typ = typ.Elem()
	}

	if typ.Kind() == reflect.Struct {
		columns := make([]string, 0, typ.NumField())
		for i := 0; i < typ.NumField(); i++ {
			fieldType := typ.Field(i)
			if fieldType.IsExported() {
				if fieldType.Tag.Get("gorm") == "-" {
					continue
				}
				name := Camel2Snake(fieldType.Name)
				if len(fieldType.Tag.Get("gorm")) > 0 {
					context := fieldType.Tag.Get("gorm")
					if strings.HasPrefix(context, "column:") {
						context = context[7:]
						pos := strings.Index(context, ";")
						if pos > 0 {
							name = context[:pos]
						} else {
							name = context
						}
					}
				}
				columns = append(columns, name)
			}
		}
		return columns
	} else {
		return nil
	}
}
