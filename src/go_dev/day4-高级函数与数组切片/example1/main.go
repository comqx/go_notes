package main

import "fmt"


func initconfig() (err error) {
	return errors.New("init config faild")
}
func test() {
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		fmt.Println(err)
	// 	}
	// }()

	err := initconfig()
	if err != nil {
		panic(err)
	}
	b := 0
	a := 100
	fmt.Println(a / b)
	return
}

func main() {
	test()

	fmt.Println(1, 2, 3, 4, 5)
}
