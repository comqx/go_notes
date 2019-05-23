package main

import (
	"fmt"
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
		var sum int
		for j := 1; j <= (i / 2); j++ {
			if i%j == 0 {
				sum += j
			}
		}
		if sum == i {
			fmt.Println(sum)
		}

	}
}

//输入一个字符串，判断是否是回文，回文字符串是指从左到右读和从右到左读完全相同的字符串
func strs() {
	var str_a string
	for {
		fmt.Scanf("%s", &str_a)
		var str_rev string
		for i := len(str_a) - 1; i >= 0; i-- {
			str_rev += fmt.Sprintf("%c", str_a[i])
		}
		if str_a == str_rev {
			fmt.Printf("string true---%s\n", str_rev)
		} else {
			fmt.Printf("string false---%s\n", str_rev)
			break
		}
	}
}
//输入一行字符，分别统计出其中英文字母，空格，数字和其他字符的个数
def count4(){
	var str_2 string
	for {
		fmt.Scanf("%s",&str_2)
		for i := 0;i<=len(str_2);i++{
			if fmt.Sprintf("%c",str_2[i])
		}
	}
}

//计算俩个大数相加的和，这俩个大数会超过int64的表示范围

func main() {
	// jj9()
	// ws2()
	strs()
}
