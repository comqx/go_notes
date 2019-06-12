package main

import "fmt"

// 切片的内存布局-1
type slice struct { //定义一个结构体
	ptr *[5]int //定义一个指针
	len int
	cap int
}

func make1(s slice, cap int) slice {
	s.ptr = new([5]int)
	s.cap = cap
	s.len = 0
	return s
}

func modify(s slice) {
	s.ptr[1] = 1000
}
func testslice2() {
	var s1 slice
	s1 = make1(s1, 10)
	s1.ptr[0] = 100
	modify(s1)          //修改s.ptr[1]的值
	fmt.Println(s1.ptr) // &[100 1000 0 0 0]
}

// 切片的内存布局-2
func modify1(a []int) {
	a[1] = 1000
}
func testSlice3() {
	var b []int = []int{1, 2, 3, 4} //定义一个数组
	modify1(b)                      //通过函数修改数组b索引1的值
	fmt.Println(b)                  // [1 1000 3 4]，b数组发生变化
}

//
func testslice4() {
	var a = [10]int{1, 2, 3, 4, 5} //定义一个数组
	b := a[1:5]                    //定义一个切片
	fmt.Printf("%p\n", b)          //打印切片b的内存地址
	fmt.Printf("%p\n", &a[1])      //打印数组a的内存地址
}

// 切片初始化练习
func testslice() {
	var slice []int                          //切片定义，没有长度
	var arr [5]int = [...]int{1, 2, 3, 4, 5} //定义一个数组

	//长度和容量是相同的
	slice = arr[2:5]
	fmt.Println(slice) //[3 4 5]
	fmt.Println(len(slice))
	fmt.Println(cap(slice))

	//长度和容量是不相同的
	slice = slice[0:1]
	fmt.Println(slice) //[3]
	fmt.Println(len(slice))
	fmt.Println(cap(slice))
}

func main() {
	testslice()
	testslice2()
	testSlice3()
}
