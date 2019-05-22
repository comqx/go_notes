package main

import (
	"fmt"
	"math"
)

//九九乘法表
func jj9() {
	for i := 1; i <= 9; i++ {
		for j := 1; j <= i; j++ {
			product := i * j
			fmt.Printf("%d*%d=%d ", i, j, product)
		}
		fmt.Printf("\n")
	}
}

//一个数如果恰好等于它的因子之和，这个数称为"完数"，例如6=1+2+3，编程找出1000以内的完数
func ws2() {
	for i := 1; i <= 1000; i++ {
		var sum int = 0
		for j := 1; j <= int(math.Sqrt(float64(i))); j++ {
			if i%j == 0 {
				sum += j
			}
		}
		if sum == i {
			fmt.Println(sum)
		}
	}
}

func main() {
	// jj9()
	ws2()
}
