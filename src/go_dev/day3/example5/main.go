package main

import (
	"fmt"
	"math/rand"
	"strconv"
)

// // 1、if/else分支判断
// func if1() {
// 	//1、if
// 	if 条件1 {

// 	}
// 	//2、if/else
// 	if 条件2 {

// 	} else {

// 	}
// 	//3、多次判断
// 	if 条件3 {

// 	} else if 条件 {

// 	} else if 条件 {

// 	}
// }

//2、判断是否报错

func if2() {
	var str string
	fmt.Scanf("%s", &str)
	number, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println("convert faild,err:", err)
		return
	}
	fmt.Println(number)
}

//3、if分支判断或者使用switch case语句 (switch 后面有条件)
func switch3() {
	var a int
	fmt.Scanf("%s", &a)
	fmt.Printf("%T", a)
	switch a {
	case 0:
		fmt.Println("a is eq 0")
		fmt.Println("yes", 0)
		fallthrough //这个分支执行完成后，继续往下面的分支走
	case 10, 11:
		fmt.Println("a is eq 10")
	default:
		fmt.Println("a is eq default")
	}
}

//4、if分支判断或者使用switch case语句 (switch 后面无条件,case后面是条件)
func switch4() {
	var a int
	fmt.Scanf("%s", &a)
	switch {
	case a > 0 && a < 10:
		fmt.Println("a is eq 0")
		fmt.Println("yes", a)
	case a > 10:
		fmt.Println("a is eq 10")
	default:
		fmt.Println("a is eq default")
	}
}

//5、使用switch case语句 (switch 后面有条件,case后面是条件)
func switch5() {
	var a int
	fmt.Scanf("%s", &a)
	switch a := 100; {
	case a > 0 && a < 10:
		fmt.Println("a is eq 0")
		fmt.Println("yes", a)
	case a > 10:
		fmt.Println("a is eq 10")
	default:
		fmt.Println("a is eq default")
	}
}

//6、猜数字练习，随机生成一个0到100的数，用户在终端输入数字，如果等于n，提示猜对了，如果不相同，提示大于还是小于
func nums() {
	var n int
	n = rand.Intn(100) //100以内的随机数
	for {
		var input int
		fmt.Scanf("%d\n", &input)
		flag := false
		switch {
		case input == n:
			fmt.Println("you are right")
			flag = true
		case input > n:
			fmt.Println("bigger")
		default:
			fmt.Println("less")
		}
		if flag {
			break
		}
	}
}

func main() {
	// if2()
	// switch3()
	// switch4()
	// switch5()
	nums()
}
