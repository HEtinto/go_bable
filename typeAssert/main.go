package main

import "fmt"

/*
	 	Golang 中的接口是一种抽象类型,可以存储任何实现了该接口方法的类型实例.
		然而,由于接口本身不包含类型信息,需要通过类型断言来将接口变量转换为实际类型.+----------------+

接口变量内存:
|    iface       |
+----------------+
| data (指针)    | -> 具体类型值的地址
| tab (指针)     | -> 指向itab
+----------------+
+----------------+
|    itab       |
+----------------+
| inter (指针)  | -> 接口类型描述符
| _type (指针)  | -> 具体类型描述符
| fun (数组)    | -> 方法函数指针
+----------------+

[类型断言的最佳实践]
1.避免过度依赖类型断言，频繁使用类型断言可能是设计上的问题。如果发现自己在使用大量的类型断言的时候，需要停下来审视下类型设计是否合理，良好的设计应尽量减少类型断言的使用。
2.安全地使用类型断言，尽可能使用带 ok 的形式进行类型断言，避免程序 panic，使程序更加健壮。
3.当有多个可能的类型需要断言时，可以使用类型分支（type switch），这是一种特殊的类型断言形式，可以更清晰地处理多个类型
*/

// 实现不同类型执行不同操作
func typeCheck(v interface{}) {
	fmt.Println("typeCheck run param:", v)
	switch v.(type) {
	case int:
		fmt.Println("int")
	case string:
		fmt.Println("string")
	default:
		fmt.Println("unknown")
	}
}

func main() {
	var i interface{} = "hello world"
	// 类型断言
	s, ok := i.(string)
	if ok {
		fmt.Println(s)
	} else {
		fmt.Println("类型断言失败")
	}
	typeCheck(i)
}
