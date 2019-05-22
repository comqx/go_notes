package main

import "fmt"

// //返回值可以是一个变量
// func add(a int, b int) (c int) {
// 	c = a + b
// 	return
// }

// func calc(a, b int) (sum int, avg int) {
// 	sum = a + b
// 	avg = (a + b) / 2
// 	return
// }
// func return_main() {
// 	add1 := add(10, 20)
// 	calc1, _ := calc(100, 200)
// 	fmt.Println(add1, calc1)
// }

//14、写一个函数add,支持1个或多个int相加，并返回相加结果
func adds(a int, bs ...int) (sum int) {
	sum = a
	for i := 0; i < len(bs); i++ {
		sum = sum + bs[i]
	}
	return
}
func sums_main() {
	sums := adds(1, 123, 123, 43, 412)
	fmt.Println(sums)
}

//15、写一个concat，支持1个或者多个string相拼接，并返回结果
func concat(a string, bs ...string) (concats string) {
	concats = a
	for i := 0; i < len(bs); i++ {
		concats = concats + bs[i]
	}
	return
}
func concat_main() {
	concats := concat("lll", "qqq", "xxx", "www")
	fmt.Println(concats)
}
func main() {
	// return_main()
	sums_main()
	concat_main()
}
