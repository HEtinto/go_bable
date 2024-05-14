package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

func add(x int, y int) int   { return x + y }
func sub(x, y int) (z int)   { z = x - y; return }
func first(x int, _ int) int { return x }
func zero(int, int) int      { return 0 }

// 多返回值函数 bare return
func swap(x, y string) (x1 string, y1 string) {
	x1 = y
	y1 = x
	return
}

// 多返回值函数 normal return
func swap1(x, y string) (string, string) {
	x1 := y
	y1 := x
	return x1, y1
}

// 错误处理之重试
func func1() {
	const timeout = 1 * time.Minute
	deadline := time.Now().Add(timeout)
	var a int = 1
	var err *int = &a
	for tries := 0; time.Now().Before(deadline); tries++ {
		fmt.Println("do something...")
		// 增长重试时间
		if err == nil {
			fmt.Println("直接返回")
			return
		} else {
			fmt.Println("重试...")
			err = nil
		}
		time.Sleep(time.Second << uint(tries))
	}
}

// 从标准io读取数据,并检查文件结束符
func funcReadStdio() error {
	in := bufio.NewReader(os.Stdin)
	for {
		r, _, err := in.ReadRune()
		if err == io.EOF {
			break // finish reading
		}
		if err != nil {
			return fmt.Errorf("read failed:%v", err)
		}
		fmt.Printf("%c", r)
	}
	return nil
}

func print1(s string) {
	fmt.Println("print1:", s)
}

// 函数参数
func funcUseFuncParam(s string, print func(s string)) {
	print(s)
}

// squares返回一个匿名函数。
// 该匿名函数每次被调用时都会返回下一个数的平方。
func squares() func() int {
	var x int
	return func() int {
		x++
		return x * x
	}
}

// 可变参数求和函数
func sum(vals ...int) int {
	total := 0
	for _, val := range vals {
		total += val
	}
	return total
}

// 可变参数max函数
func max(vals ...int) int {
	max := vals[0]
	for _, val := range vals {
		if max < val {
			max = val
		}
	}
	return max
}

// defer function
func testdefer(num int) {
	if num == 0 {
		defer fmt.Println("num == 0")
	}
	if num == 1 {
		defer fmt.Println("num == 1")
	}
	if num == 2 {
		defer fmt.Println("num == 2")
	}
}

// recover从panic异常中恢复
func testrecover() {
	defer func() {
		// 在触发了panic异常后，会执行defer函数中的代码
		// recover调用会使得程序不会崩溃
		// 如果调用recover前没有触发panic，recover会返回nil
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	panic("test recover")
}

// 使用panic和recover编写一个不包含return语句但能返回一个非零值的函数
func testpanic1() (ret int) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			ret = 1024
		}
	}()
	panic(2)
}

func main() {
	fmt.Printf("%T\n", add)   // "func(int, int) int"
	fmt.Printf("%T\n", sub)   // "func(int, int) int"
	fmt.Printf("%T\n", first) // "func(int, int) int"
	fmt.Printf("%T\n", zero)  // "func(int, int) int"
	// 实参通过值的方式传递，因此函数的形参是实参的拷贝。
	// 对形参进行修改不会影响实参。但是，如果实参包括引用类型，如指针，
	// slice(切片)、map、function、channel等类型，
	// 实参可能会由于函数的间接引用被修改。
	fmt.Println(swap("1", "2"))
	fmt.Println(swap1("1", "2"))
	func1()
	// funcReadStdio()
	funcUseFuncParam("hello world", print1)
	// 函数值,也称闭包
	// 从下面的输出可以看到函数值记录了状态,匿名函数和squares之间存在变量的引用
	fmt.Println("闭包测试")
	f := squares()
	fmt.Println(f()) // "1"
	fmt.Println(f()) // "4"
	fmt.Println(f()) // "9"
	// 陷阱 捕获迭代变量 在使用range时也会有该问题
	fmt.Println("函数值陷阱")
	var printfunc []func()
	for i := 0; i < 3; i++ {
		printfunc = append(printfunc, func() {
			fmt.Println(strconv.Itoa(i))
		})
	}
	for _, f := range printfunc {
		// 目标是打印 0 1 2 但是由于i作为引用被函数值使用,所以i是同一个
		// 在最后一次迭代中 i的值为3
		f()
	}
	// 修改
	fmt.Println("创建新的变量来解决函数值/闭包引用问题")
	var printfunc1 []func()
	for i := 0; i < 3; i++ {
		j := i // 创建一个新的变量
		printfunc1 = append(printfunc1, func() {
			fmt.Println(strconv.Itoa(j))
		})
	}
	for _, f := range printfunc1 {
		// 目标是打印 0 1 2 但是由于i作为引用被函数值使用,所以i是同一个
		// 在最后一次迭代中 i的值为3
		f()
	}

	fmt.Println("可变参数函数传递普通参数:", sum(1, 2, 3))
	fmt.Println("可变参数函数传递切片:", sum([]int{1, 2, 3}...))
	fmt.Println("最大值:-1, 0, 3:", max(-1, 0, 3))
	// 测试defer函数
	testdefer(0)
	testdefer(1)
	testdefer(2)
	// recover and panic
	fmt.Println(testpanic1())
}
