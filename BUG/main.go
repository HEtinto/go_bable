/*
此文件用于收集go语言常见使用陷阱,在学习的过程中掌握正确使用go语言的方式
*/
package main

import (
	"errors"
	"fmt"
	"time"
)

// for range 重用变量导致的bug
func bug1() {
	m := [...]int{1, 2, 3, 4}
	for i, v := range m {
		go func() {
			time.Sleep(time.Second * 3)
			// 此处使用的是外层定义的i, v, 由于打印前有延时, 最终i, v的值是最后一次循环的值
			// 如果立即打印而不加延时, 那么i, v的值是不确定的, 由go routine调度器和循环决定
			// 因为i, v的值始终由外层i, v变量决定
			fmt.Println(i, v)
		}()
	}
	time.Sleep(time.Second * 10)
}

// 修正bug1
func fix1() {
	m := [...]int{1, 2, 3, 4}
	for i, v := range m {
		go func(i, v int) {
			time.Sleep(time.Second * 3)
			fmt.Println(i, v)
		}(i, v) // 传递参数为值拷贝传递, 所以传递的参数与原变量解除绑定
	}
	time.Sleep(time.Second * 10)
}

// for range 参与迭代的是原变量的副本
func bug2() {
	a := [5]int{1, 2, 3, 4, 5}
	r := [5]int{}
	fmt.Printf("original a:%+v\n", a)

	for i, v := range a {
		if i == 0 {
			a[1] = 20
			a[2] = 30
		}
		// 我们期望r[i]将保存修改后的a数组的值
		// 但是由于i, v遍历的是a的拷贝,即原始a数组
		// 所以r将保存原始a数组的值
		r[i] = v
	}
	/* for i, v := range a 等价于(其中a1是a的拷贝)
	for i, v := range a1 {
		if i == 0 {
			a[1] == 20
			a[2] == 30
		}
		r[i] = v
	}
	*/
	fmt.Printf("modified a:%+v\n", a) // modified a:[1 20 30 4 5]
	fmt.Printf("copy r:%+v\n", r)     // copy r:[1 2 3 4 5]
}

// 使用指针作为for range的迭代表达式
/*注意
数组切片不同,数组是拷贝整个数组,与原数组属于完全不同的内存区域,
而切片本身的设计只是引用数组,拷贝引用也还是指向原数组，所以切片并不会产生bug2的情形.
*/
func fix2() {
	a := [5]int{1, 2, 3, 4, 5}
	r := [5]int{}
	fmt.Printf("original a:%+v\n", a)

	// 拷贝值改成拷贝指针,指针此时仍然指向a
	for i, v := range &a {
		if i == 0 {
			a[1] = 20
			a[2] = 30
		}
		r[i] = v
	}
	fmt.Printf("modified a:%+v\n", a) // modified a:[1 20 30 4 5]
	fmt.Printf("copy r:%+v\n", r)     // copy r:[1 20 30 4 5]
}

// 虽然切片使用指向底层数组的引用来修改数据,但是保存引用的切片结构体也是值拷贝
// 这将导致一些问题
// 切片的结构体类似于: struct {array *p, len int, cap int}, 其中array是底层数组的指针
// len是切片的长度, cap是容量, 在拷贝时len和cap是值拷贝,对他们的修改都不会反映到原始切片
func bug3() {
	a := []int{1, 2, 3, 4, 5}
	r := []int{}
	fmt.Printf("original a:%+v\n", a)
	for i, v := range a { // 此处a是副本切片结构体
		if i == 0 {
			a = append(a, 6, 7) // 此处修改了原始a,而非副本a,原始切片将被修改
			// 由于副本a的len没有被改变, 所以迭代次数仍为5, 虽然此时原始a的len为7
		}
		r = append(r, v)
	}
	fmt.Printf("modified a:%+v\n", a) // modified a:[1 2 3 4 5 6 7]
	fmt.Printf("copy r:%+v\n", r)     // copy r:[1 2 3 4 5]
}

// 使用临时数据保存新增元素,在最后追加
func fix3() {
	a := []int{1, 2, 3, 4, 5}
	r := []int{}
	fmt.Printf("original a:%+v\n", a)
	tmp := []int{}
	for i, v := range a { // 此处a是副本切片结构体
		if i == 0 {
			a = append(a, 6, 7)
			tmp = append(tmp, 5, 7)
		}
		r = append(r, v)
	}
	r = append(r, tmp...)
	fmt.Printf("modified a:%+v\n", a) // modified a:[1 2 3 4 5 6 7]
	fmt.Printf("copy r:%+v\n", r)     // copy r:[1 2 3 4 5]
}

// for range 遍历string
// string在go内部表示为 struct {*byte, len}
// 不过for range对于string来说,
// 每次循环的单位是一个rune,而不是一个byte,返回的第一个值为迭代字符码点的第一字节的位置
func bug4() {
	s := "中国人" // 每个字用三个字节表示
	for i, v := range s {
		fmt.Printf("%d %s 0x%x\n", i, string(v), v)
		// 0 中 0x4e2d
		// 3 国 0x56fd
		// 6 人 0x4eba
	}
	// 如果range表达式的string中含有非法utf8字节序列
	// 那么v将返回0xfffd,并在下一次迭代中仅仅前进一字节
	// byte sequence of s: 0xe4 0xb8 0xad 0xe5 0x9b 0xbd 0xe4 0xba 0xba
	var sl = []byte{0xe4, 0xb8, 0xad, 0xe5, 0x9b, 0xbd, 0xe4, 0xba, 0xba} // "中国人"
	// 0xe4表示开头的字节序是一个三字节的utf8字符
	// 0xb8和0xad是该三字节字符的后续字节
	// 0xe5表示开肉的字节序是一个三字节的utf8字符...
	for _, v := range sl {
		fmt.Printf("0x%x\t", v)
	}
	fmt.Println("\n构造非法utf8字节序列")

	// 故意构造非法UTF8字节序列
	// 在utf8中,0xd0到0xdf范围的字节只能出现在两字节序列
	// 的第二个字节位置,而不是三字节序列的第三字节位置
	// 这将违反utf8编码规则
	sl[3] = 0xd0
	sl[4] = 0xd6
	sl[5] = 0xb9

	for i, v := range string(sl) {
		fmt.Printf("%d 0x%x, s:%v\n", i, v, string(v))
		// 0 0x4e2d, s:中
		// 3 0xfffd, s:�    程序遇到错误的utf8编码
		// 4 0x5b9, s:ֹ       该字符是一个合法两字节希伯来语字符,utf8可以解析
		// 6 0x4eba, s:人    程序回归正常,获取到了正确的字符
	}
	fmt.Printf("len s1:%v\n", len(sl))

}

/*
map在for range循环中是一个副本,不过是对map描述结构的副本
这个副本指向的是同一个hmap描述结构,所以修改会反映到源map上
同时由于map的遍历顺序是随机的,如果在遍历map的过程中对map进行
修改,那么修改产生的结果是无法预知的
*/
func bug5() {
	m := map[string]int{
		"tony": 21,
		"tom":  22,
		"jim":  23,
	}
	counter := 0
	for k, v := range m {
		if counter == 0 {
			delete(m, "tony") // 当删除map中的"tony"时,k的值可能为"tony",但"tony"会被从map中删除
		}
		counter++
		fmt.Println(k, v)
	}
	fmt.Println("counter is", counter)
	fmt.Printf("map: %+v\n", m)
	// 可能的结果1
	// 	tony 21
	// tom 22
	// jim 23
	// counter is 3
	// map: map[jim:23 tom:22]
	// 可能的结果2
	// jim 23
	// tom 22
	// counter is 2
	// map: map[jim:23 tom:22]
}

// 在for range表达式中使用chan时,chan会阻塞
// 在表达式上,直到有数据读取,或者chan关闭
// 所以在使用时需要注意关闭chan,否则程序将一直阻塞
// 同时需要注意 chan并非0值可用的,这表示你向一个var c chan int发送数据将会导致程序阻塞
// 这和map不支持零值可用是一样的
func bug6() {
	ch := make(chan int)
	go func() {
		time.Sleep(3 * time.Second)
		ch <- 1
		ch <- 2
		ch <- 3
		close(ch) // 关闭chan,防止阻塞
	}()
	for v := range ch {
		fmt.Println(v)
	}
	// 不支持零值可用

	// 运行下面的代码将导致死锁
	// var c chan int
	// go func() {
	// 	c <- 100
	// }()
	// for v := range c {
	// 	fmt.Println(v)
	// }
}

// break调出select但未调出循环
// chapter3/sources/control_structure_idiom_6.go
func bug7() {
	exit := make(chan interface{})

	go func() {
		for {
			select {
			case <-time.After(time.Second):
				fmt.Println("tick")
			case <-exit:
				fmt.Println("exiting...")
				break // break没有跳出for循环,只是跳出了select
				// 在主进程的挂起结束后,主进程结束时,select将会退出
			}
		}
		fmt.Println("exit!")
	}()

	time.Sleep(3 * time.Second)
	exit <- struct{}{}

	// wait child goroutine exit
	time.Sleep(3 * time.Second)
}

// 使用标签来跳出
func fix7() {
	exit := make(chan interface{})

	go func() {
	loop:
		for {
			select {
			case <-time.After(time.Second):
				fmt.Println("tick")
			case <-exit:
				fmt.Println("exiting...")
				break loop // break没有调出for循环,只是调出了select
			}
		}
		fmt.Println("exit!")
	}()

	time.Sleep(3 * time.Second)
	exit <- struct{}{}

	// wait child goroutine exit
	time.Sleep(3 * time.Second)
}

// 使用case 表达式列表fallthrough
func bug8() {
	a := int(3)
	switch a {
	case 1, 2, 3:
		fmt.Println("a <= 3")
	case 4:
		fmt.Println("a == 4")
	default:
		fmt.Println("a > 4")
	}
}

// defer求值的时机,defer在执行是进行求值
func bug9() {
	a := []int{1, 2, 3, 4}
	defer func(a []int) {
		fmt.Println("deferred function a:", a)
	}(a) // 此时a为[1 2 3 4], deferred函数将输出[1 2 3 4]
	a = append(a, 5, 6, 7)
	fmt.Println("The value of a slice: ", a)
}

// for range 循环变量重用问题
// 这里涉及的是go在调用方法是进行的隐式自动转换
// 在对非指针类型的实例调用指针方法时发生隐式转换
// 导致了循环变量重用

type field struct {
	name string
}

func (p *field) print() {
	fmt.Println(p.name)
}

func bug10() {
	data1 := []*field{{"one"}, {"two"}, {"three"}}
	for _, v := range data1 {
		go v.print() // 此处传入的是data1各个元素的地址
		// 等价形式：go (*field).print(v) // v是各个元素的地址
	}

	time.Sleep(3 * time.Second)
	data2 := []field{{"four"}, {"five"}, {"six"}}
	for _, v := range data2 {
		go v.print() // 此处传入的是同一个变量v的地址,即复用了v
		// 由于传递的是同一个v的拷贝,所以最终打印出来的值依赖于go routine的调度
		// 等价形式：go (*field).print(&v) // 同一个v
	}

	time.Sleep(3 * time.Second)
}

// 接口类型等值判断条件: 接口类型中类型相等时才相等
// 在下面的例子中,returnsError函数中的p的接口类型是MyError
// 和nil的不同,nil的类型为0,data字段也为0
// 接口内部实现
// 需要注意的是,两个非空接口或两个空接口进行判断时需要判断类型和指向的数据都相等才相等
// 但是非空接口和空接口进行判断时,只需要判断类型相等即可
/*
// $GOROOT/src/runtime/runtime2.go
type iface struct {
    tab  *itab
    data unsafe.Pointer
}

type eface struct {
    _type *_type
    data  unsafe.Pointer
}
*/
type MyError struct {
	error
}

var ErrBad = MyError{
	error: errors.New("bad error"),
}

func bad() bool {
	return false
}

func returnsError() error {
	var p *MyError = nil
	if bad() {
		p = &ErrBad
	}
	return p
}

func bug11() {
	e := returnsError()
	// e的接口类型中的类型为MyError,与空接口类型nil的类型0x0不相符
	// 所以e不等于nil
	if e != nil {
		fmt.Printf("error: %+v\n", e)
		return
	}
	fmt.Println("ok")
}

func main() {
	bug11()
}
