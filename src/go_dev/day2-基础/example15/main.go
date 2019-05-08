package main

import "fmt"

func main() {
	var str = "hello world\n"
	var str1=`
	abc
	123
	自行车

	asdf
	`
	var b byte = 'c'
	fmt.Println(str)
	fmt.Println(str1)
	fmt.Println(b)
	fmt.Printf("%c\n",b)
}