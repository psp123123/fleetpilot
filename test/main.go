package main

import (
	"fmt"
	"reflect"
)

type User struct {
	ID   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

func InspectStruct(s interface{}) {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return
	}

	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		fmt.Println(
			"\nField:", field.Name,
			"\nType:", field.Type,
			"\nTag(db):", field.Tag.Get("db"),
			"\nValue:", value.Interface(),
			"\n-------------------",
		)
	}
}

func main() {
	s := User{101, "zhangsan"}
	InspectStruct(s)
	InspectStruct(&s)
}
