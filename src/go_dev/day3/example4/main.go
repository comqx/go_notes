package main

import (
	"fmt"
)

// for语句的写法：
// for 初始化语句；条件判断；变量修改{
// }

// 分层打印多个A
func prints(n int) {
	for i := 0; i < n+1; i++ {
		for j := 0; j < i; j++ {
			fmt.Printf("A")
		}
		fmt.Println()
	}
}

//1、 for +条件
func for1(i int) {
	//for后面加条件
	for i > 0 {
		fmt.Println("i>0")
	}
	//for死循环，俩种写法
	for true {
		fmt.Println("i>0")
	}
	for {
		fmt.Print("abcdsdf")
	}
}

//2、for range语句,用来遍历数组，slice，map，chan
func for2() {
	str := "hello world,中国"
	for i, v := range str {
		fmt.Printf("index[%d] val[%c] len[%d]\n", i, v, len([]byte(string(v))))
	}
}

//3、break continue语句
func for3() {
	str := "hello world,中国"
	for i, v := range str {
		if i > 2 {
			continue
		}
		if i > 3 {
			break
		}
		fmt.Printf("index[%d] val[%c] len[%d]\n", i, v, len([]byte(string(v))))
	}
}

//4-1、goto和label语句
func for4() {
label1:
	for i := 0; i <= 5; i++ {
		for j := 0; j <= 5; j++ {
			if j == 4 {
				continue label1 //continue label1，就直接跳转到label1的地方再次执行代码（内层的for循环会终止）
			}
			fmt.Printf("i is: %d,and j is:%d\n", i, j)
		}
	}
}

//4-2、goto和label语句(here)
func for4_2() {
	i := 0
here:
	println(i)
	i++
	if i == 5 {
		return
	}
	goto here //这个只能在同一个函数里面去goto到指定代码，继续执行代码
}

func main() {
	// prints(10)
	// for1(10)
	// for2()
	// for3()
	// for4()
	for4_2()

}
