package main

import (
	"fmt"
	"reflect"
)

/*
使用反射的弊端:
	1.代码不易阅读,不易维护,容易发生线上panic
	2.性能很差,比正常代码慢一到两个数量级
*/

type User struct {
	name string
	age  int
	sex  string
}

// 指针Type转为非指针Type
func case1() {
	typeUser1 := reflect.TypeOf(&User{})
	typeUser2 := reflect.TypeOf(User{})
	if typeUser1.Elem() == typeUser2 {
		fmt.Println("typeUser1.Elem() == typeUser2")
	} else {
		fmt.Println("typeUser1.Elem()!= typeUser2")
	}
}

// 获取结构体成员信息
func case2() {
	typeUser := reflect.TypeOf(User{}) //需要用struct的Type，不能用指针的Type
	fieldNum := typeUser.NumField()    //成员变量的个数
	for i := 0; i < fieldNum; i++ {
		field := typeUser.Field(i)
		fmt.Printf("member [%d] %s offset %d anonymous %t type %s exported %t json tag %s\n",
			i,                     //第几个成员
			field.Name,            //变量名称
			field.Offset,          //相对于结构体首地址的内存偏移量，string类型会占据16个字节
			field.Anonymous,       //是否为匿名成员
			field.Type,            //数据类型，reflect.Type类型
			field.IsExported(),    //包外是否可见（即是否以大写字母开头）
			field.Tag.Get("json")) //获取成员变量后面``里面定义的tag
	}
	fmt.Println()

	//可以通过FieldByName获取Field
	if nameField, ok := typeUser.FieldByName("Name"); ok {
		fmt.Printf("Name is exported %t\n", nameField.IsExported())
	}
	//也可以根据FieldByIndex获取Field
	thirdField := typeUser.FieldByIndex([]int{2}) //参数是个slice，因为有struct嵌套的情况
	fmt.Printf("third field name %s\n", thirdField.Name)
}

// learn web:https://zhuanlan.zhihu.com/p/411313885
func main() {
	// 1.获取类型
	typeI := reflect.TypeOf(1)
	typeS := reflect.TypeOf("hello")
	fmt.Println(typeI)
	fmt.Println(typeS)

	typeUser := reflect.TypeOf(&User{})
	fmt.Println(typeUser)               //*main.User
	fmt.Println(typeUser.Kind())        //ptr
	fmt.Println(typeUser.Elem().Kind()) //struct

	case1()
	case2()
}
