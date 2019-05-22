package main

import "fmt"

//无论是值传递，还是引用传递，传递给函数的都是变量的副本。
func modify(a int) {
	a = 100
}
func main() {
	a := 8
	fmt.Println(a)
	modify(a)
	fmt.Println(a)
}
