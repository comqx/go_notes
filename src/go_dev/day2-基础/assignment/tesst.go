package main

import (
	"fmt"
)

func sums(n int) uint64{
	var s uint64 = 1
	var sum uint64 = 0
	for i :=1;i<=n;i++{
		s1 := s * uint64(i)
		fmt.Printf("%d!=%v*%d,,,%d\n",i,s,i,s1)
		sum+=s1
		s=s1
	}
	return sum
}
//1!=1
//2!=2*1=2  2!=1!*2
//3!=3*2*1=6  3!=2!*3
//4!=4*3*2*1=24  4!=3!*4
//5!=5*4*3*2*1=120  5!=4!*5

func main(){
	var a int
	fmt.Scanf("%d",&a)
	for i:=1;i<=a;i++{
		sums(i)
	}
}
