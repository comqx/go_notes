package main

import "fmt"

// 1、函数是可以赋值给一个变量的
func add(a, b int) int {
	return a + b
}
func func1_main() {
	c := add
	fmt.Println(c)   //这个变量打印出来的值就是一个地址
	sum := c(10, 20) //可以直接为c这个变量添加参数进行获取函数返回的结果
	fmt.Println(sum)
}

// 2、自定义函数类型
//1-// type op_func func(int, int) int //自定义一个函数类型op_func，传入值为int,int，返回值为int

func add2(a, b int) int {
	return a - b
}

//1-// func operator(op op_func, a, b int) int {
func operator(op func(int, int) int, a, b int) int { //也可以使用把函数当做类型传入
	return op(a, b)
}
func func2_main() {
	c := add2
	sum2 := operator(c, 100, 2000) //传入函数，和俩个值
	fmt.Println(sum2)
}

func main() {
	func1_main()
	func2_main()
}
