package add

import (
	_ "../test" //只想引用这个包，做初始化，但是不使用包内容
)

var Name string = "add-name"

// var age int = 123  //变量必须是大写，才能在这个包外面引用
var Age int = 123

// var Name string //name=   string默认值是空
// var Age int // age= 0  int 默认值是0

//包里面的init函数，一定是在包之前执行的
func init() {
	Name = "add-init-name"
	Age = 10
}

//> name= helloword
//> age= 10
