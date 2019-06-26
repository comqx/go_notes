


<!-- TOC -->

- [golang语言特性](#golang语言特性)
- [包的概念](#包的概念)
- [go程序目录结构](#go程序目录结构)
- [go程序基本结构](#go程序基本结构)
- [文件名&关键字&标识符](#文件名关键字标识符)
- [函数的声明和注释](#函数的声明和注释)
- [学习golang](#学习golang)

<!-- /TOC -->
# golang语言特性
1. 垃圾回收
- 内存自动回收，再也不需要开发人员管理内存
- 开发人员专注业务实现，降低了心智负担
- 只需要new分配内存，不需要释放
2. 天然并发
- 从语言层面支持并发，非常简单
- goroute，轻量级线程，创建成千上万个goroute称为可能
- 基于csp（communicating sequential process）模型实现
3. channel
- 管道，类似uninx/linux中的pipe
- 多个goroute之间通过channel进行通信
- 支持任何类型
4. 多返回值
- 一个函数返回多个值
# 包的概念
1. 和python一样，把相同功能的代码放到一个目录，称之为包
2. 包可以被其他包引用
3. main包是用来生成可执行文件，每个程序只有一个main包
4. 包的主要用途是提高代码的可复用性
# go程序目录结构
```go
-d:
 -/project #项目目录
  -/src #放的是我们的代码
   -/go_dev
    -/day1
     -/example1/
      -/
  -/bin #放的是可执行文件
  -/vender #放的是第三方包
  -/pkg #静态库
export GOPATH=d:/project/ #指定项目位置
```
# go程序基本结构
```go
package main //说明是个包

import "fmt" //导入fmt包

func main() {
    fmt.Println("hello,world")
}

/*注释：
1. 任何一个代码文件隶属于一个包
2. import关键字， 引用其他包
    import("fmt")
    import("os") 通常写为：
    import(
        "fmt"
        "os"
    )
3. golang可执行程序，package main，并且有且只有一个main入口函数
4. 包中函数调用：
    a. 同一个包中函数，直接调用
    b. 不同包中函数，通过包名+点+函数名进行调用
5. 包访问控制规则：
    a. 大写意味着这个函数/变量是可导出的
    b. 小写意味着这个函数/变量是私有的，包外不可访问
    */
```
# 文件名&关键字&标识符
1. 所有go源码以.go结尾
2. 标识符以字母或下划线开头，大小写敏感，比如：
```go
a. boy   √
b. Boy   √
c. a+b   x
d. Oboy  √
e. _boy  √
f. =_boy  x
g. _      √  //特殊的标识符，a,_,b，意思是不取值
```
3. _是特殊标识符，用来忽略结果
4. 保留关键字
5. 
break | default | func | interface | select（用在chann管道里面）
------|---------|------|-----------|--------------------
case | defer | go | map | struct
chan | else | goto | package | switch
const（声明一个常量） | fallthough | if | range | type
continue | for | import | return | var

# 函数的声明和注释
```go
1.函数声明：func 函数名称（参数列表）（返回值列表）{}
// 举例：
func add(){}
func add(a int,b int) int {}
func add(a int,b int) (int,int) {}

2.注释，俩种注释，单行注释：//和多行注释/**/
```
# 学习golang
```
ctrl+alt+m toc
ctrl+alt+x 粘贴image
clt+ w 预览
ctrl+shift+p 查看全面命令
```

