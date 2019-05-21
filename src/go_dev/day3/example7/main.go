package main

import "fmt"

func modify(a int) {
	a = 100
}
func main() {
	a := 8
	fmt.Println(a)
	modify(a)
	fmt.Println(a)
}
