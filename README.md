[TOC]

# golang语言特性


1. 垃圾回收
- 内存自动回收，再也不需要开发人员管理内存
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
# 编译go

## go build 

```go
go build -o 指定编译后的名字 需要编译的pkg


```

## go run

像执行脚本文件一样执行Go代码

## go install

`go install`分为两步：

	1. 先编译得到一个可执行文件
	
	2. 将可执行文件拷贝到`GOPATH/bin`

## go doc

`go doc builtin.delete ` 查看builtin.delete的用法

## go get [-alrtAFR]
```
    # 显示操作流程日志及信息
    -v
    # 下载丢失的包，但不更新已经存在的包
    -u
    # 只下载，不自动安装
    -d
    # 允许使用 HTTP 方式进行下载操作
    -insecure
```
## go env
查看环境变量

```shell 
GOROOT
GOPATH
当我们导入一个包xxx时:
go系统会优先在GOROOT/src中寻找，然后在GOPATH/src中寻找

```

## go fmt

格式化go代码文件

## go list

列出全部安装的package

## go tool

 查看汇编语言 go tool compile -S  main.go 



## 交叉编译

Go支持跨平台编译

例如：在windows平台编译一个能在linux平台执行的可执行文件

```bash
SET CGO_ENABLED=0  // 禁用CGO
SET GOOS=linux  // 目标平台是linux
SET GOARCH=amd64  // 目标处理器架构是amd64
```

执行`go build`

```shell
# mac上编译linux和windows二进制
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build 
# 如果你想在Windows 32位系统下运行
CGO_ENABLED=0 GOOS=windows GOARCH=386 go build test.go
# 如果你想在Windows 64位系统下运行
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build 

CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=7 CC=arm-linux-gnueabi-gcc-4.7 go build



CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows GOARCH=amd64 go build -x -v -ldflags "-s -w" ./cmd/foot-api/main.go      


# linux上编译mac和windows二进制
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build 
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build

# windows上编译mac和linux二进制
SET CGO_ENABLED=0 SET GOOS=darwin SET GOARCH=amd64 go build main.go
SET CGO_ENABLED=0 SET GOOS=linux SET GOARCH=amd64 go build main.go
```

# go module

> 使用go1.11 go1.12的版本，需要手动开启go module支持
>
> 开启办法：
>
> `set GO111MODULE=on //windows`
>
> `export GO111MODULE=on //mac` 

## goproxy

设置代理：

```go
export GOPROXY=https://goproxy.cn    Mac
SET GOPROXY=https://goproxy.cn       Windows

// windows:
go env -w GO111MODULE=on 
go env -w GOPROXY=https://goproxy.cn,direct

// linux mac
export GO111MODULE=on
export GOPROXY=https://goproxy.cn
```

## go.mod文件

记录了当前项目依赖的第三方包信息和版本信息

第三方的依赖包都下载到了 `GOPATH/pkg/mod`目录下。

## go.sum文件

详细包名和版本信息

## 常用的命令

```go
go mod init [包名] // 初始化项目

go mod tidy // 检查代码里的依赖去更新go.mod文件中的依赖

go get 或者go mod download

go get -u [包名] // 下载包
go get github.com/wilk/uuid@0.0.1 // 指定版本

go clean -i -n  [包名]  // 卸载包
```

# 学习golang

```
ctrl+alt+m toc
ctrl+alt+x 粘贴image
clt+ w 预览
ctrl+shift+p 查看全面命令

Command + Shift + [ 折叠代码块
Command + Shift + ] 展开代码块
Command + K Command + [ 折叠全部子代码块
Command + K Command + ] 展开全部子代码块
Command + K Command + 0 折叠全部代码块
Command + K Command + J 展开全部代码块
```

