package main

import (
	"fmt"
)
func main(){
	var n int16 = 34
	var m int32

	// m=n // cannot use n (type int16) as type int32 in assignment
	m=int32(n)
	fmt.Println("m=%d,n=%d\n",m,n)
}