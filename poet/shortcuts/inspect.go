package shortcuts

import (
	"fmt"
	"reflect"
)

func Inspect(abc interface{}) {

	// 获取反射类型对象
	t := reflect.TypeOf(abc)

	// 遍历所有字段
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fmt.Printf("Field: %s %v\n", f.Name, f.Type)
	}

	// 遍历所有方法
	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		fmt.Printf("Method: %s %v\n", m.Name, m.Type)
	}

}
