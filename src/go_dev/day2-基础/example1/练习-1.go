
//随机给个数a，打印出来全部相加等于a的全部相加的过程
package main

import (
	"fmt"
)

func assig_1(a int) {
	for i := 0; i <= a; i++ {
		fmt.Printf("%d+%d=%d\n",i,a-i,a)
	}
}

func main() {
	assig_1(100)
}
