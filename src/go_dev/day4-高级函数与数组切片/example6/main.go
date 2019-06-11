package main

import (
	"fmt"
)

// 访问越界问题
func test1() {
	var a [10]int
	// j := 10 // 不能超出数组的范围，会报错panic
	j := 9
	a[0] = 10
	a[j] = 100
	fmt.Println(a)

	for i := 0; i < len(a); i++ {
		fmt.Println(a[i])
	}
	for index, val := range a {
		fmt.Printf("a[%d]=%d\n", index, val)
	}
}

//数组是值类型，改变副本的值，不会改变本身的值
func test2() {
	var a [10]int
	b := a
	b[0] = 101
	fmt.Println(a) //[0 0 0 0 0 0 0 0 0 0] a的值没有发生变化
}
func test3(arr [5]int) {
	arr[0] = 1000
}

// 如果需要改变原来的数组的值需要传入地址进去
func test4(arr *[5]int) {
	arr[0] = 1000
}

func main() {
	test1()
	test2()

	var a [5]int
	test3(a)
	fmt.Println(a) //[0 0 0 0 0]
	test4(&a)      //&a传入值类型的地址
	fmt.Println(a) //[1000 0 0 0 0]
}
