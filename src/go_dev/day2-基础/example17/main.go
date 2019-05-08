package main
import "fmt"

func reverse(str string) string {
	var result string
	strlen:= len(str)
	for i :=0;i < strlen;i++ {
		result=result + fmt.Sprintf("%c",str[strlen-i-1])
	}
	return result
}

func reverse1(str string) string {
	var result []byte //数组
	tmp := []byte(str)
	lenth := len(str)
	for i :=0 ;i < lenth;i++{
		result = append(result,tmp[lenth-i-1])
	}
	return string(result)
}

func main(){
	var str1="hello"
	str2 := "world"
	// str3 := str1 +" "+ str2 //字符串拼接 方式1
	str3 := fmt.Sprintf("%s %s",str1,str2) //字符串拼接 方式2
	n := len(str3) //字符串长度
	fmt.Println(str3)
	fmt.Printf("len(str3)=%d\n",n)

	substr := str3[0:5] //切片
	fmt.Println(substr)

	substr = str3[6:] //切片
	fmt.Println(substr)
	//str反转，方法1
	result := reverse(str3)
	fmt.Println(result)
	fmt.Printf("%v",result)
	//str反转，方法2
	result = reverse1(result)
	fmt.Println(result)
}