<!-- TOC -->

- [strings和strconv使用](#strings%E5%92%8Cstrconv%E4%BD%BF%E7%94%A8)
	- [strings使用](#strings%E4%BD%BF%E7%94%A8)
	- [strconv 整数-字符串转换](#strconv-%E6%95%B4%E6%95%B0-%E5%AD%97%E7%AC%A6%E4%B8%B2%E8%BD%AC%E6%8D%A2)

<!-- /TOC -->
# strings和strconv使用
## strings使用
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

## strconv 整数-字符串转换
```go
// 1. strconv.Itoa(i int): 把几个整数i转换为字符串
// 2. strconv.Atoi(str string)(int,error): 把一个字符串转换为整数
```　