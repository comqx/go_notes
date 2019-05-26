package main

import (
	"errors"
	"fmt"
	"time"
)

//定义一个报错函数
func initConfig() (err error) {
	return errors.New("init config failed ")
}

// recover和panic的使用
func test() {
	defer func() { //定义一个错误捕获匿名函数
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	b := 0
	a := 100 / b
	fmt.Println(a)
	//主动抛出异常
	err := initConfig()
	if err != nil {
		panic(err)
	}
	return
}

func main() {
	for {
		test()
		time.Sleep(time.Second)
	}
	var a []int                 //定义一个切片，[]如果是空的话是切片，[5]是一个长度为5的数组
	a = append(a, 10, 20, 3242) //在a这个切片里面追加一些数值
	a = append(a, a...)         //a... 意思是把a这个切片展开
	fmt.Println(a)

}
