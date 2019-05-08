1、判断101-200之间有多个素数，并输出所有素数
```go
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
```
2、打印出100-999中所有的水仙花数，所谓水仙花数是指一个3位数，期个位数字立方和等于该数本身。例如：153是一个水仙花数，因为153=1的三次方+5的三次方+3的三次方
```go
package main
import "fmt"
func isNumber(n int) bool {
	var i,j,k int
	i = n %10
	j = (n /10) %10
	k = (n / 100)%10
	//fmt.Printf("%d,%d,%d\n",i,j,k)
	sum := i*i*i + j*j*j + k*k*k
	return sum == n
}
func main(){
	var n int
	var m int
	fmt.Scanf("%d,%d",&n,&m)
	for i :=n;i<=m;i++{
		if isNumber(i) == true{
			fmt.Println(i,"is shuixianhua")
		}
	}
}
```
3、对于一个数n,求n的阶乘之和，即1！+2！+3！+---n!
```go
1!=1
2!=2*1=2  2!=1!*2
3!=3*2*1=6  3!=2!*3
4!=4*3*2*1=24  4!=3!*4
5!=5*4*3*2*1=120  5!=4!*5

package main
import "fmt"

func sum(n int) uint64{
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
func main(){
	var n int
	fmt.Scanf("%d",&n)
	s := sum(n)
	fmt.Println(s)
}
```
