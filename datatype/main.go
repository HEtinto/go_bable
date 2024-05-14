package main

import (
	"fmt"
)

func reverse(intptr *[5]int) {
	for i := int(len(intptr) / 2); i >= 0; i-- {
		(*intptr)[i], (*intptr)[len(intptr)-i-1] = (*intptr)[len(intptr)-i-1], (*intptr)[i]
	}
	// fmt.Println(*intptr)
}

func main() {
	// 数组 具有固定大小
	r := [...]int{99: 1}
	fmt.Println(r)
	a1 := [...]int{1, 2, 3}
	a2 := [...]int{3, 2, 1}
	fmt.Println(a1 == a2)
	// 重写reverse函数
	array1 := [...]int{1, 2, 3, 4, 5}
	reverse(&array1)
	fmt.Println("执行反转后的数组:", array1)
	// map
	// 第一种创建方式
	// 其中make函数参数:参数1:数据类型(slice,map,channel),参数2:长度,参数3:底层数据容量
	ages := make(map[string]int)
	ages["alice"] = 31
	ages["charlie"] = 34
	fmt.Println(ages)
	// 禁止对map中的元素进行取地址,因为map中的元素会动态增长
	// 第二种创建方式
	ages1 := map[string]int{
		"alice":   31,
		"charlie": 34,
	}
	fmt.Println(ages1)
	// 删除
	delete(ages1, "charlie")
	fmt.Println(ages1)
	// 添加/修改
	ages1["yujianming"] = 100
	fmt.Println(ages1)
	// 结构体和匿名类型
	type Point struct {
		x, y int
	}
	type Circle struct {
		Point
		radius int
	}
	var c Circle
	c.x = 1
	c.y = 2
	c.radius = 3
	fmt.Println(c)
	c = Circle{
		Point:  Point{x: 4, y: 4},
		radius: 6,
	}
	fmt.Printf("%#v\n", c)
}
