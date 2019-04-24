//2. 一个程序包含俩个包add和main,其中add包有俩个变量：Name和age，请问main包如何访问Name和age
//-> Name可以被跨包访问到，age在编译的时候出错（未定义），也就是说大写可以跨包访问到，小写的只能是个私有的

//3. 包别名的应用，开发一个程序，使用包别名来访问包中的函数

//4. 每个源代码文件都可以包含一个init函数，这个init函数自动被go运行框架调用

package main

import (
	"fmt"

	a "../add" //3. 包别名
)

func main() {
	fmt.Println("name=", a.Name)
	fmt.Println("age=", a.Age)
}
