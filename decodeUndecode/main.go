package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"
	"unicode/utf8"
)

// 多语言支持
func multipleLang() {
	// 16进制Unicode码点
	s := "\xe4\xb8\x96\xe7\x95\x8c"
	fmt.Printf("%s\n", s)
	// 16位Unicode码点
	s = "\u4e16\u754c"
	fmt.Printf("%s\n", s)
	// 32位Unicode码点
	s = "\U00004e16\U0000754c"
	fmt.Printf("%s\n", s)
	// 混合语言字符串长度
	s = "Hello, 世界"
	fmt.Printf("len method:%d\n", len(s))
	fmt.Printf("utf8.RuneCountInString:%d\n", utf8.RuneCountInString(s))
	// 使用utf8解码器
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		fmt.Printf("%d\t%c\n", i, r)
		i += size
	}
	// 使用range隐式地解码utf8字符串
	fmt.Println("使用range隐式地解码utf8字符串")
	for i, r := range "Hello, 世界" {
		fmt.Printf("%d\t%q\t%d\n", i, r, r)
	}
}

type User struct {
	Name string `json:"name"` // tag为结构体字段提供元数据
	Age  int    `json:"age"`  // 用于自定义序列化和反序列化行为
	Sex  string `json:"sex"`  //
	// 在与其他系统或格式交互时十分有用
}

func printStructFields(v interface{}) {
	val := reflect.ValueOf(v)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)
		fmt.Printf("%s: %v\n", field.Name, value.Interface())
	}
}

func testfunc(a int32) {
	fmt.Printf("testing1\n")
	if a == 0 {
		fmt.Printf("testing2\n")
		return
	}
	if a == 1 {
		fmt.Printf("testing3\n")
		defer func() {
			fmt.Printf("testing5\n")
		}()
		return
	}
	fmt.Printf("testing4\n")

}

func testfunc2() {
	rc := make(chan error, 10)
	for i := 0; i < 10; i++ {
		count := uint32(i)
		go func(rc chan error, count uint32) {
			time.Sleep(time.Duration(count) * time.Second)
			rc <- fmt.Errorf("error1")
		}(rc, count)
	}
	for i := 0; i < 10; i++ {
		err := <-rc
		fmt.Printf("%v\n", err)
	}
}

func main() {
	user := User{Name: "Alice", Age: 30}
	jsonBytes, _ := json.Marshal(user)
	fmt.Printf("(Serialization)序列化:%v\n", string(jsonBytes))

	var new_user User
	_ = json.Unmarshal(jsonBytes, &new_user)
	fmt.Printf("(Deserialization)反序列化:%v\n", new_user)
	printStructFields(new_user)
	if new_user.Sex == "" {
		fmt.Println("Sex is empty")
	}

	str1 := ""
	if a, err := strconv.ParseUint(str1, 10, 64); err != nil {
		fmt.Println("str1:", a)
	}

	testfunc(4)
	testfunc2()

}
