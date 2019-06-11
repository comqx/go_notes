package main

import "fmt"

func fab(n int) {
	var a []uint64
	a = make([]uint64, n)

	a[0] = 1
	a[1] = 1
	for i := 2; i < n; i++ {
		a[i] = a[i-1] + a[i-2]
	}
	for _, v := range a {
		fmt.Println(v)
	}
}

func init_sz() {
	var age0 [5]int = [5]int{1, 2, 3}

	var age1 = [5]int{1, 2, 3, 4, 5}
	var age2 = [...]int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	var str = [5]string{3: "hello world", 4: "tom"}
	var age3 = [...]int{2: 100, 5: 200}
	fmt.Println(age0, age1, age2, str, age3)
}

func dw() {
	var f [2][5]int = [...][5]int{{1, 2, 3, 4, 5}, {6, 7, 8, 9, 10}} //2行5列

	for row, v := range f { //遍历行，v其实也就是一个数组
		for col, v1 := range v { //遍历列
			fmt.Printf("{%d,%d}=%d", row, col, v1) //坐标的值
		}
		fmt.Println()
	}
}
{0,0}=1{0,1}=2{0,2}=3{0,3}=4{0,4}=5
{1,0}=6{1,1}=7{1,2}=8{1,3}=9{1,4}=10
func main() {
	fab(10)
	init_sz()
	dw()
}
