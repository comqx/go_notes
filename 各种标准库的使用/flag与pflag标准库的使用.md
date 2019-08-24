[TOC]

pflag 包的设计目的就是替代标准库中的 flag 包，因此它具有更强大的功能并且与标准的兼容性更好。本文将介绍 pflag 包与 flag 包相比的主要优势，

# pflag

## pflag 包的主要特点

pflag 包与 flag 包的工作原理甚至是代码实现都是类似的，下面是 pflag 相对 flag 的一些优势：

- 支持更加精细的参数类型：例如，flag 只支持 uint 和 uint64，而 pflag 额外支持 uint8、uint16、int32 等类型。
- 支持更多参数类型：ip、ip mask、ip net、count、以及所有类型的 slice 类型。
- 兼容标准 flag 库的 Flag 和 FlagSet：pflag 更像是对 flag 的扩展。
- 原生支持更丰富的功能：支持 shorthand、deprecated、hidden 等高级功能。

## pflag的使用

### NoOptDefVal-为参数设置默认值之外的值

### MarkDeprecated 指定被启用的参数

### MarkShorthandDeprecated 指定被启用的参数简写

### MarkHidden 隐藏参数

### SetNormalizeFunc 解决传参不规范的问题

```go
package main

import (
	"fmt"
	"strings"

	flag "github.com/spf13/pflag"
)

// 定义命令行参数对应的变量
var cliName = flag.StringP("name", "n", "nick", "Input your name")
var cliAge = flag.IntP("age", "a", 22, "Input your age")
var cliOK = flag.BoolP("ok", "o", false, "Are you ok")
var cliDes = flag.StringP("des-detail", "d", "", "Input Description")
var cliOldFlag = flag.StringP("badflag", "b", "", "Input badflag")
var cliGender = flag.StringP("gender", "g", "male", "Input Your Gender")

func wordSepNormailzeFunc(f *flag.FlagSet, name string) flag.NormalizedName {
	from := []string{"-", "_"}
	to := "."
	for _, sep := range from {
		name = strings.Replace(name, sep, to, -1)
	}
	return flag.NormalizedName(name)
}

func main() {
	// 设置标准化参数名称的函数
  // 如果我们创建了名称为 --des-detail 的参数，但是用户却在传参时写成了 --des_detail 或 --des.detail 会怎么样？默认情况下程序会报错退出，但是我们可以通过 pflag 提供的 SetNormalizeFunc 功能轻松的解决这个问题
	flag.CommandLine.SetNormalizeFunc(wordSepNormailzeFunc)

	// 为 age 参数设置 NoOptDefVal 默认值，通过简便的方式为参数设置默认值之外的值
	// 使用-a 参数，不用添加任何值默认就是25，不使用-a 默认就是22
	flag.Lookup("age").NoOptDefVal = "25"

	//把 badflag 参数标记为即将废弃的，请用户使用 des-detail 参数
	//./main  -b asd   Flag shorthand -b has been deprecated, please use -d instead
	flag.CommandLine.MarkDeprecated("badflag", "please use --des-detail instead")

	//把 badflag 参数的 shorthand 标记为即将废弃的，请用户使用 des-detail 的 shorthand 参数
	///main  -b asd    Flag --badflag has been deprecated, please use --des-detail instead
	flag.CommandLine.MarkShorthandDeprecated("badflag", "please use -d instead")

	//在帮助文档中隐藏参数 gender
	flag.CommandLine.MarkHidden("gender")

	// 把用户传递的命令行参数解析为对应变量的值
	flag.Parse()

	fmt.Println(*cliName, *cliAge, *cliOK, *cliDes, *cliOldFlag)
}

----
localhost:11pflag liuqixiang$ ./main  -h
Usage of ./main:
  -a, --age int[=25]        Input your age (default 22)
  -d, --des.detail string   Input Description
  -n, --name string         Input your name (default "nick")
  -o, --ok                  Are you ok
pflag: help requested
                                             
localhost:11pflag liuqixiang$ ./main  -b abc
Flag shorthand -b has been deprecated, please use -d instead
Flag --badflag has been deprecated, please use --des-detail instead
nick 22 false  abc
                                             
localhost:11pflag liuqixiang$ ./main  -g abc
nick 22 false  
```

# flag

## flag包练习

```go
func main() {
	//定义命令行参数方式1
	var name string
	var age int
	var married bool
	var delay time.Duration
	flag.StringVar(&name, "name", "张三", "姓名")
	flag.IntVar(&age, "age", 18, "年龄")
	flag.BoolVar(&married, "married", false, "婚否")
	flag.DurationVar(&delay, "d", 0, "延迟的时间间隔")

	//解析命令行参数
	flag.Parse()
	fmt.Println(name, age, married, delay)
	//返回命令行参数后的其他参数
	fmt.Println(flag.Args())
	//返回命令行参数后的其他参数个数
	fmt.Println(flag.NArg())
	//返回使用的命令行参数个数
	fmt.Println(flag.NFlag())
}
---------------------------------------------------

$ ./flag_demo -help
Usage of ./flag_demo:
  -age int
        年龄 (default 18)
  -d duration
        时间间隔
  -married
        婚否
  -name string
        姓名 (default "张三")
```