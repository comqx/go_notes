package main

import "fmt"

func test() {
	s1 := new([]int) //切片

	fmt.Println(s1) //返回指针  &[]

	s2 := make([]int, 10) //长度为10的切片
	fmt.Println(s2)

	(*s1)[0] = 100
	s2[0] = 100
	return
}

func main() {
	test()
}
