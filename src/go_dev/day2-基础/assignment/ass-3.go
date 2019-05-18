package main
import (
	"fmt"
	"math"
)

func isPrime(n int ) bool {
	for i:=2;i <= int(math.Sqrt(float64(n))); i++{  //开方
		if n%i==0 {
			return false
		}
	}
	return true
}

func main(){
	var n int
	var m int 
	fmt.Scanf("%d%d",&n,&m) //从终端输入字符串，自动转换为int存入变量里面
	for i := n;i <m; i++{
		if isPrime(i) == true{
			fmt.Printf("%d\n",i)
			continue
		}
	}
}