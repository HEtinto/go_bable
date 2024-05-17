package main

import "fmt"

// 1.函数类型转换
type BinaryAdder interface {
	Add(int, int) int
}

// 函数类型,用于将自定义的函数转为实现了BinaryAdder Add的方法
type MyAdderFunc func(int, int) int

// MyAdderFunc 实现了BinaryAdder Add方法
func (f MyAdderFunc) Add(x, y int) int {
	return f(x, y)
}

func MyAdd(x, y int) int {
	return x + y
}

// 使用
func tips1() {
	i := MyAdderFunc(MyAdd)
	// i.Add实际调用的是MyAdd函数
	fmt.Println(i.Add(1, 2))
}

// 2.柯里化函数:将接收多个参数的函数转为接收一个单一函数
// 原始函数
func times(x, y int) int {
	return x * y
}

// 柯里化
func partialTimes(x int) func(int) int {
	return func(y int) int {
		return times(x, y)
	}
}
func tips2() {
	// 获取一个函数
	f := partialTimes(2)
	// 调用函数
	fmt.Println(f(3)) // 2*3
}

// 函子:函子是容器类型,该容器类型要实现
// 一个接收函数类型参数的方法,并在容器
// 的每个元素上应用那个函数
type IntSliceFunctor interface {
	Fmap(fn func(int) int) IntSliceFunctor
}

type intSliceFunctorImpl struct {
	ints []int
}

func (isf intSliceFunctorImpl) Fmap(fn func(int) int) IntSliceFunctor {
	newInts := make([]int, len(isf.ints))
	for i, elt := range isf.ints {
		retInt := fn(elt)
		newInts[i] = retInt
	}
	return intSliceFunctorImpl{ints: newInts}
}

func NewIntSliceFunctor(slice []int) IntSliceFunctor {
	return intSliceFunctorImpl{ints: slice}
}

func tips3() {
	// 原切片
	intSlice := []int{1, 2, 3, 4}
	fmt.Printf("init a functor from int slice: %#v\n", intSlice)
	f := NewIntSliceFunctor(intSlice)
	fmt.Printf("original functor: %+v\n", f)

	mapperFunc1 := func(i int) int {
		return i + 10
	}

	mapped1 := f.Fmap(mapperFunc1)
	fmt.Printf("mapped functor1: %+v\n", mapped1)

	mapperFunc2 := func(i int) int {
		return i * 3
	}
	mapped2 := mapped1.Fmap(mapperFunc2)
	fmt.Printf("mapped functor2: %+v\n", mapped2)
	fmt.Printf("original functor: %+v\n", f) // 原函子没有改变
	fmt.Printf("composite functor: %+v\n", f.Fmap(mapperFunc1).Fmap(mapperFunc2))
}

// 4.延续传递式
// 原始函数
func Max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

// 修改后
func Max1(a, b int, f func(c int)) {
	if a > b {
		f(a)
	} else {
		f(b)
	}
}

func tips4() {
	// 原始函数使用
	fmt.Printf("Max(1, 2): %d\n", Max(1, 2))
	// 修改后函数的使用
	Max1(1, 2, func(c int) {
		fmt.Printf("Max1(1, 2): %d\n", c)
	})
}

// 5.使用defer拦截panic
func bar() {
	fmt.Println("raise a panic")
	panic(-1)
}

func foo() {
	defer func() {
		// recover函数从返回panic,如果程序之前没有触发panic,那么将返回nil
		if err := recover(); err != nil {
			fmt.Println("recover from panic")
		}
	}()
	bar()
}

func tips5() {
	foo()
	fmt.Println("main run normally.")
}

// 6.使用deferred函数修改函数的具名返回值
func foo1(i int) (j int) {
	defer func() {
		// deferred函数在foo1函数跳转到main函数执行前将j加1
		j++
	}()
	return i
}

func tips6() {
	fmt.Println(foo1(1))
}

// 7.使用deferred函数打印日志
func tips7() {
	defer func() {
		fmt.Println("Print some log info.")
	}()
}

// 8.还原初始值
func tips8() {
	oldValue := 100
	func() {
		oldValueBak := oldValue
		fmt.Println("old value: ", oldValue)
		defer func() {
			// 从备份中恢复旧值
			oldValue = oldValueBak
		}()
		oldValue = 200
		fmt.Println("change old value to:", oldValue)
	}()
	fmt.Println("after function run, old value:", oldValue)
}

func main() {
	tips8()
}
