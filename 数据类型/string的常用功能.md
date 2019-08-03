- [strings方法使用](#strings方法使用)
- [*内置函数len()](#内置函数len)
- [*rune计算string的字符数](#rune计算string的字符数)
# strings方法使用

```go

package main

import (
	"fmt"
	"strings"
)
// 1. strings.HasPrefix(s string, prefix string) bool: 判断字符串s是否以prefix开头
func urlProcess(url string) string {
	result := strings.HasPrefix(url,"http://") //HasPrefix判断是否以指定字符串开头
	if !result {
		url = fmt.Sprintf("http://%s",url)
	}
	return url
}
func pathProcess(path string) string{
	result := strings.HasPrefix(path,"/")
	if !result{
		path = fmt.Sprintf("%s/",path)
	}
	return path
}
func prefix(){
	var (
		url string
		path string
	)
	fmt.Scanf("%s%s",&url,&path)
	url = urlProcess(url )
	path = pathProcess(path )
	fmt.Printf("%s,%s",url,path)

}
// 2. strings.HasSuffix(s string, Suffix string) bool 判断字符串s是否以Suffix结尾
// 3. strings.Index(s string,str string) int: 判断str在s中首次出现的位置，如果没有出现，返回-1
// 4. strings.LastIndex(s string,str string) int: 判断str在s中最后出现的位置，如果没有出现，返回-1
// 5. strings.Replace(str string, old string, new string, n int): 字符串替换,n是替换几次
func replace(){
	result := strings.Replace("heheheworld","he","12",2)
	// results := strings.Trim("heheheworld","asdasd")
	fmt.Printf("%s",result)
}
// 6. strings.Count(str string, substr string) int: 字符串计数
// 7. strings.Repeat(str string,count int)string: 重复count次str
// 8. strings.ToLower(str string)string: 转为小写
// 9. strings.ToUpper(str string)string: 转为大写
// 10. strings.TrimSpace(str string): 去掉字符串首尾空白字符
//     strings.Trim(str string,cut string): 去掉字符串首cut字符,没有返回空
//     strings.TrimLeft(str string,cut string): 去掉字符串左边（首）cut字符
//     strings.TrimRight(str string,cut string): 去掉字符串右边（末尾）cut字符
// 11. strings.Field(str string): 返回str空格分隔的所有字符串的slice
//     strings.Split(str string,split string): 返回str split分隔的所有字符串的slice
func field(){
	str := "helloworld"
	splitResult:= strings.Fields(str)
	for i :=0;i<len(splitResult);i++{
		fmt.Println(splitResult[i])
	}

	splitResult= strings.Split(str,"w")
	for i :=0;i<len(splitResult);i++{
		fmt.Println(splitResult[i])
	}
}
// 12. strings.Join(s1 []string,sep string):用sep把s1中的所有元素链接起来

func main() {
	// prefix()
	// replace()
	field()
}
```

# *内置函数len()

>len()计算的是字符串的byte类型的长度

```go
  var s = "helloworld,温暖"
	lenStr := len(s)           //len()计算的是字符串的byte类型的长度
	fmt.Printf("%d\n", lenStr) //17
```

# *rune计算string的字符数

> runne 会把string转化成int32切片，并且存储的是字符对应的acsii码

```go
  var s = "helloworld,温暖"
	runeStr := []rune(s)   //runne 会把string转化成int32切片，并且存储的是字符对应的acsii码
	fmt.Printf("%T,%d\n", runeStr, runeStr) //[]int32,[104 101 108 108 111 119 111 114 108 100 44 28201 26262]
	fmt.Printf("%c\n", runeStr)             //使用%c，获取acsii码对应的字符
	lenRuneStr := len(runeStr)
	fmt.Printf("s字符串里面有%d字符\n", lenRuneStr) //s字符串里面有13字符
```