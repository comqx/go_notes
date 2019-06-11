package main

import "fmt"

func testslice() {
	var slice []int
	var arr [5]int = [...]int{1, 2, 3, 4, 5}

	//长度和容量是相同的
	slice = arr[2:5]
	fmt.Println(slice) //[3 4 5]
	fmt.Println(len(slice))
	fmt.Println(cap(slice))

	//长度和容量是不相同的
	slice = slice[0:1]
	fmt.Println(slice) //[3]
	fmt.Println(len(slice))
	fmt.Println(cap(slice))
}

func main() {
	testslice()
}
