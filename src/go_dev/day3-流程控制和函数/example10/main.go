package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

//1. 九九乘法表
func jj9() {
	for i := 1; i <= 9; i++ {
		for j := 1; j <= i; j++ {
			product := i * j
			fmt.Printf("%d*%d=%d\t", i, j, product)
		}
		fmt.Printf("\n")
	}
}

//2. 一个数如果恰好等于它的因子之和，这个数称为"完数"，例如6=1+2+3，编程找出1000以内的完数
func ws2() {
	for i := 1; i <= 100000; i++ {
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

//3. 输入一个字符串，判断是否是回文，回文字符串是指从左到右读和从右到左读完全相同的字符串
// > 新知识点：[]rune(str_a)
func strs() {
	var str_a string
	for {
		fmt.Scanf("%s", &str_a)
		t := []rune(str_a) //rune表示一个字符,比如一个汉字字符占用3个字节
		length := len(t)
		for i, _ := range t {
			if i == length/2 {
				break
			}
			last := length - i - 1
			if t[i] != t[last] {
				fmt.Println("no")
			}
		}
		fmt.Println("yes")
		//第二种方式，不适用与中文
		// var str_rev string
		// for i := len(str_a) - 1; i >= 0; i-- {
		// 	str_rev += fmt.Sprintf("%c", str_a[i])
		// }
		// if str_a == str_rev {
		// 	fmt.Printf("string true---%s\n", str_rev)
		// } else {
		// 	fmt.Printf("string false---%s\n", str_rev)
		// 	break
		// }
	}
}

//4. 输入一行字符，分别统计出其中英文字母，空格，数字和其他字符的个数
// >新函数,[]rune(),bufio.NewReader(os.Stdin),reader.ReadLine()
func count4(str string) (worldCount, spaceCount, numberCount, otherCount int) {
	t := []rune(str) //表示一个字符
	for _, v := range t {
		switch {
		case v >= 'a' && v <= 'z':
			fallthrough
		case v >= 'A' && v <= 'Z':
			worldCount++
		case v == ' ':
			spaceCount++
		case v >= '0' && v <= '9':
			numberCount++
		default:
			otherCount++
		}
	}
	return
}
func read_line() {
	reader := bufio.NewReader(os.Stdin) //终端输入
	result, _, err := reader.ReadLine() //读取行
	if err != nil {
		fmt.Println("read form console err:", err)
		return
	}
	wc, sc, nc, oc := count4(string(result))
	fmt.Printf("worldcount:%d\nspacecount:%d\nnumvercount:%d\nothercount:%d\n", wc, sc, nc, oc)

}

//5. 计算俩个大数相加的和，这俩个大数会超过int64的表示范围
func multi(str1, str2 string) (result string) {
	if len(str1) == 0 && len(str2) == 0 {
		result = "0"
		return
	}
	var index1 = len(str1) - 1
	var index2 = len(str2) - 1
	var left int
	for index1 >= 0 && index2 >= 0 {
		c1 := str1[index1] - '0'
		c2 := str2[index2] - '0'
		sum := int(c1) + int(c2) + left
		if sum >= 10 {
			left = 1
		} else {
			left = 0
		}
		c3 := (sum % 10) + '0'
		result = fmt.Sprintf("%c%s", c3, result)
		index1--
		index2--
	}
	for index1 >= 0 {
		c1 := str1[index1] - '0'
		sum := int(c1) + left
		if sum >= 10 {
			left = 1
		} else {
			left = 0
		}
		c3 := (sum % 10) + '0'
		result = fmt.Sprintf("%c%s", c3, result)
		index1--
	}
	for index2 >= 0 {
		c1 := str2[index2] - '0'
		sum := int(c1) + left
		if sum >= 10 {
			left = 1
		} else {
			left = 0
		}
		c3 := (sum % 10) + '0'
		result = fmt.Sprintf("%c%s", c3, result)
		index2--
	}
	if left == 1 {
		result = fmt.Sprintf("1%s", result)
	}
	return
}

func sum_main() {
	reader := bufio.NewReader(os.Stdin)
	result, _, err := reader.ReadLine()
	if err != nil {
		fmt.Println("read from console err:", err)
	}
	strSlice := strings.Split(string(result), "+")
	if len(strSlice) != 2 {
		fmt.Println("please input a+b")
	}

	strnumber1 := strings.TrimSpace(strSlice[0])
	strnumber2 := strings.TrimSpace(strSlice[1])
	fmt.Println(multi(strnumber1, strnumber2))

}

func main() {
	// jj9()
	// ws2()
	// strs()
	// read_line()
	sum_main()
}
