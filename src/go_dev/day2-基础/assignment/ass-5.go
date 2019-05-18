package main
import "fmt"

func sums(n int) uint64{
	var s uint64 = 1
	var sums uint64 = 0
	for i :=1;i<=n;i++{
		s1 := s * uint64(i)
		fmt.Printf("%d!=%v*%d,,,%d\n",i,s,i,s1)
		sums+=s1
		s=s1
		}
	return sums
}
func main(){
	var n int
	fmt.Scanf("%d",&n)
	s := sums(n)
	fmt.Println(s)
}