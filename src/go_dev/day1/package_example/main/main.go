package main

import (
	"fmt"

	"../calc"
)

func main() {
	sum := calc.Add(100, 300) //注意calc包里面的add函数必须是大写，这边才能调用
	sub := calc.Sub(100, 200)
	fmt.Println("sum=", sum)
	fmt.Println("sub=", sub)
}
