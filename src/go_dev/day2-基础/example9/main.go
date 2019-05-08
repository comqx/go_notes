package main

//3种方式，实现a,b值的交换
//第一种是，交换指针，赋值内存地址的形式
//第二种是，函数的形式，返回值去交换值
//第三种是，直接交换赋值
import (
	"fmt"
)
func swap(a *int ,b *int ){  // *指针的意思，
	tmp:= *a   //*a,tmp指向a的指针
	*a=*b   //a的指针等于b的指针
	*b=tmp  //b的指针等于tmp
	return
}

func swap1(a int,b int)(int,int ){
	return b,a 
}

func test(){
	var a=100  //a 的作用域是在函数内部
	fmt.Println(a)
	for i := 0 ;i <100 ;i++{
		var b = i * 2  //b的作用域在这个语句块里面
		fmt.Println(b)
	}
	// fmt.println(b)  //在这一层就不能出现b的变量

	if(a >0){
		var c int = 100  //这个c的作用域也在这个if语句块里面
		fmt.Println(c)
	}
}

func main(){
	first := 100
	second := 200
	//1/ swap(&first,&second) //传入内存地址，传入的是一个副本
	//2/ first,second = swap1(first,second) //利用函数返回值的形式

	first,second=second,first  //3/第三种方式交换first，second的值

	fmt.Println("first=",first)
	fmt.Println("second=",second)

	test()
}