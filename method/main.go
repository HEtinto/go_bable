package main

import "fmt"

type Point struct {
	x, y int
}

// 定义嵌套Point的新类型,该类型将包含Point的方法
type Point1 struct {
	Point
	x1, y1 int
}

// 为Point类型实现一个print方法
func (p Point) print1() {
	fmt.Println("Point methd call:", p.x, p.y)
}

// 为Point类型实现一个指针方法
func (p *Point) print() {
	if p == nil {
		fmt.Println("nil point")
		return
	}
	fmt.Println("(*Point) methd call:", (*p).x, (*p).y)
}

// Bit数组实现集合
// An IntSet is a set of small non-negative integers.
// Its zero value represents the empty set.
type IntSet struct {
	words []uint64
}

// Has reports whether the set contains the non-negative value x.
func (s *IntSet) Has(x int) bool {
	word, bit := x/64, uint(x%64)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

// Add adds the non-negative value x to the set.
func (s *IntSet) Add(x int) {
	word, bit := x/64, uint(x%64)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

// UnionWith sets s to the union of s and t.
func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

func main() {
	p := Point{1, 2}
	p.print1()
	// 指针方法
	p1 := &Point{3, 4}
	p1.print()
	// 隐式地调用指针方法
	p.print()
	p1 = nil
	// 如果方法中定义了对nil的处理,那么这个指针方法可以传递nil参数
	p1.print()
	// 测试嵌套类型
	p2 := Point1{Point{1, 2}, 3, 4}
	p2.print1()
	p2.print()
	// 使用匿名结构体实现嵌套类型Point1
	p3 := struct {
		Point
		x1, y1 int
	}{
		Point: Point{1, 2},
		x1:    3,
		y1:    4,
	}
	fmt.Println("匿名结构体操作")
	p3.print1()
	p3.print()
}
