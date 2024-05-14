package main

import (
	"fmt"
	"sort"
)

type ByteCounter int

// ByteCounter可以作为io.Writer所要求的接口的实现传递给Fprinf函数
func (c *ByteCounter) Write(p []byte) (int, error) {
	*c += ByteCounter(len(p)) // convert int to ByteCounter
	return len(p), nil
}

// type CountWriter struct {
// 	w     Writer
// 	count int64
// }

// /*
// 传入一个io.Writer接口类型，
// 返回一个把原来的Writer封装在
// 里面的新的Writer类型和一个表
// 示新的写入字节数的int64类型指针。
// */
// func CountingWriter(w io.Writer) (io.Writer, *int64) {
// 	cw := &CountWriter{w: w}
// 	return cw, &cw.count
// }

// 排序相关
type Point struct {
	x, y int
}

type sortPoint struct {
	p    []*Point
	less func(x, y *Point) bool
}

func (sp sortPoint) Len() int           { return len(sp.p) }
func (sp sortPoint) Less(i, j int) bool { return sp.p[i].x < sp.p[j].x }
func (sp sortPoint) Swap(i, j int)      { sp.p[i], sp.p[j] = sp.p[j], sp.p[i] }

func main() {
	c := ByteCounter(10)
	// Fprintf的第一个参数为io.Writer,
	// 由于ByteCounter实现了io.Writer接口,所以可以传递给Fprintf函数
	fmt.Fprintf(&c, "hello, world")
	fmt.Println(c)
	var p = []*Point{
		{3, 2},
		{1, 3},
		{5, 6},
		{2, 4},
		{9, 10},
	}
	sort.Sort(sortPoint{p, func(x, y *Point) bool {
		return x.y < y.y
	}})
	for _, v := range p {
		fmt.Println(v.x, v.y)
	}
}
