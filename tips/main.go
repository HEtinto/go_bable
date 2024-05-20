package main

import (
	"fmt"
	"strings"
)

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

// 模拟函数重载-使用变长参数
func concat(sep string, args ...interface{}) string {
	var result string
	for i, v := range args {
		if i != 0 {
			result += sep
		}
		// v.(type)是类型断言的一种特殊形式
		// 它被用在switch语句中,用来判断接口变量v
		// 持有的具体类型.这种语法是go中的类型开关
		// 它允许你针对接口值得实际类型编写不同的
		// 处理逻辑,而不需要显式地写出每个可能类型
		// 的类型断言.
		switch v.(type) {
		case int, int8, int16, int32, int64,
			uint, uint8, uint16, uint32, uint64:
			result += fmt.Sprintf("%d", v)
		case string:
			result += fmt.Sprintf("%s", v)
		case []int:
			ints := v.([]int)
			for i, v := range ints {
				if i != 0 {
					result += sep
				}
				result += fmt.Sprintf("%d", v)
			}
		case []string:
			strs := v.([]string)
			result += strings.Join(strs, sep)
		default:
			fmt.Printf("the argument type [%T] is not supported", v)
			return ""
		}
	}
	return result
}

func tips9() {
	println(concat("-", 1, 2))
	println(concat("-", "hello", "gopher"))
	println(concat("-", "hello", 1, uint32(2),
		[]int{11, 12, 13}, 17,
		[]string{"robot", "ai", "ml"},
		"hacker", 33))
}

// 使用变长参数使函数支持默认参数,需要调用者保证自身按照
// 顺序输入参数
type record struct {
	name    string
	gender  string
	age     uint16
	city    string
	country string
}

func enroll(args ...interface{} /* name, gender, age, city = "Beijing", country = "China" */) (*record, error) {
	if len(args) > 5 || len(args) < 3 {
		return nil, fmt.Errorf("the number of arguments passed is wrong")
	}

	// 主动设置的默认参数,未设置的默认参数为初始值
	r := &record{
		city:    "Beijing", // 默认值：Beijing
		country: "China",   // 默认值：China
	}

	for i, v := range args {
		switch i {
		case 0: // name
			name, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf("name is not passed as string")
			}
			r.name = name
		case 1: // gender
			gender, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf("gender is not passed as string")
			}
			r.gender = gender
		case 2: // age
			age, ok := v.(int)
			if !ok {
				return nil, fmt.Errorf("age is not passed as int")
			}
			r.age = uint16(age)
		case 3: // city
			city, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf("city is not passed as string")
			}
			r.city = city
		case 4: // country
			country, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf("country is not passed as string")
			}
			r.country = country
		default:
			return nil, fmt.Errorf("unknown argument passed")
		}
	}

	return r, nil
}

func tips10() {
	r, _ := enroll("小明", "male", 23)
	fmt.Printf("%+v\n", *r)

	r, _ = enroll("小红", "female", 13, "Hangzhou")
	fmt.Printf("%+v\n", *r)

	r, _ = enroll("Leo Messi", "male", 33, "Barcelona", "Spain")
	fmt.Printf("%+v\n", *r)

	r, err := enroll("小吴", 21, "Suzhou")
	if err != nil {
		fmt.Println(err)
		return
	}
}

// 使用 功能选项模式 编写更加清晰的功能选择代码
type FinishedHouse struct {
	style                  int    // 0: Chinese; 1: American; 2: European
	centralAirConditioning bool   // true或false
	floorMaterial          string // "ground-tile"或"wood"
	wallMaterial           string // "latex"或"paper"或"diatom-mud"
}

type Option func(*FinishedHouse)

func NewFinishedHouse(options ...Option) *FinishedHouse {
	h := &FinishedHouse{
		// default options
		style:                  0,
		centralAirConditioning: true,
		floorMaterial:          "wood",
		wallMaterial:           "paper",
	}

	for _, option := range options {
		option(h) // 此处使用传递进来的功能选项函数去修改结构体的值
		// 参数可读性更好
		// 配置选项高度可扩展
		// 提供使用默认选项最简单的方式
		// 更加安全,调用之后调用者无法再修改option
	}

	return h
}

func WithStyle(style int) Option {
	return func(h *FinishedHouse) {
		h.style = style
	}
}

func WithFloorMaterial(material string) Option {
	return func(h *FinishedHouse) {
		h.floorMaterial = material
	}
}

func WithWallMaterial(material string) Option {
	return func(h *FinishedHouse) {
		h.wallMaterial = material
	}
}

func WithCentralAirConditioning(centralAirConditioning bool) Option {
	return func(h *FinishedHouse) {
		h.centralAirConditioning = centralAirConditioning
	}
}

func tips11() {
	fmt.Printf("%+v\n", NewFinishedHouse()) // 使用默认选项
	fmt.Printf("%+v\n", NewFinishedHouse(WithStyle(1),
		WithFloorMaterial("ground-tile"),
		WithCentralAirConditioning(false)))
}

func main() {
	tips10()
}
