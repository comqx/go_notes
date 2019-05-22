package main

import (
	"fmt"
)

func defer1() {
	var i int = 0
	defer fmt.Println(i)        //函数返回的时候才执行,因此defer的执行结果是0，而不是10，第二执行
	defer fmt.Println("second") //执行顺序，函数返回结果后，第一执行，遵守先进后出的原则
	i = 10
	fmt.Println(i)

}
func main() {
	defer1()
}
