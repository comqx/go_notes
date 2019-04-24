package test

import (
	"fmt"
)
var Name string = "test-name"
var Age int = 100

func init() {
	fmt.Println("this is a test package")
	fmt.Println("test.package.Name=", Name)
	fmt.Println("test.package.Age=", Age)
	Age = 10
	fmt.Println("test.package.Age=",Age)
}

