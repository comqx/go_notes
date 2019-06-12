package main

import "fmt"

//如果切片大于数组的大小，会自动申请内存，扩展容量
func testslice() {
	var a [5]int = [...]int{1, 2, 3, 4, 5}
	s := a[1:]
	fmt.Printf("s=%p a[1]=%p\n", s, &a[1]) //俩个地址是相同的

	s = append(s, 123)
	s = append(s, 123)
	s = append(s, 123)
	s = append(s, 123)
	s = append(s, 123)
	fmt.Printf("s=%p a[1]=%p\n", s, &a[1]) //因为已经超过数组的大小了，切片从新申请内存
	fmt.Println(s)

	var c = []int{9, 8, 7}
	s = append(s, c...) //把一个数组添加到切片中
	fmt.Println(s)
}

//
func testModifyString() {
	s := "hello world"
	s1 := []rune(s) //按照字符大小格式化

	s1[1] = '0'
	str := string(s1)
	fmt.Println(str)
}
func main() {
	testslice()
	testModifyString()
}
