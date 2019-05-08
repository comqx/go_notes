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