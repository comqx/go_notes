package main

import "fmt"

func test() {
	s1 := new([]int) //slice,切片

	fmt.Println(s1) //返回指针    &[]

	s2 := make([]int, 10) //长度为10的切片    [0 0 0 0 0 0 0 0 0 0]
	fmt.Println(s2)

	*s1 = make([]int, 5) //初始化slice，s1,这里分配5个容量的空间，默认填充0
	(*s1)[0] = 100       //给s1赋值索引0添加值 返回：&[100 0 0 0 0]
	s2[0] = 100          //给s2中索引0添加值，返回：[100 0 0 0 0 0 0 0 0 0]
	fmt.Println(s1, s2)
	return
}

func main() {
	test()
}
