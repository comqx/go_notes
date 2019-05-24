package goroute

func Add(a int, b int, c chan int) {
	sum := a + b
	c <- sum   //把sum的值传到管道c里面
}
