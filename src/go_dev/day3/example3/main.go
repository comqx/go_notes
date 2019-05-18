package main

import "fmt"

func modify(p *int) {
	fmt.Println(p)
	*p = 100000
	return
}

//5、打印一个变量的地址
func main() {
	var a int = 10
	fmt.Println(&a) //打印a变量的地址
	var p *int      //指针变量的声明
	p = &a          //把a变量的地址赋值给p指针
	fmt.Println(p)  //打印出来p的地址
	fmt.Println(*p) //打印出来p变量对应的值
	*p = 100
	fmt.Println(a, *p) //a的值发生变量，因为修改的是地址指定的值

	var b int = 999
	p = &b         //p对应的地址变成了b的地址
	*p = 5         //因此修改P地址对应的值，那么就相等于修改了b
	fmt.Println(a) //a的值不会发生变化
	fmt.Println(b) //b的值变成5

	modify(&a)
	fmt.Println(a)
}
