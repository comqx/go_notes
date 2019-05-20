package mian

import (
	"fmt"
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

//2、

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
func main() {
	if2()
}
