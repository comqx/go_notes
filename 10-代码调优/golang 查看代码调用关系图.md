[TOC]

go-callvis 是github上一个开源项目，可以用来查看golang代码调用关系。

# 安装

## 安装graphviz

```shell
$ brew install graphviz
apt install graphviz
yum install graphviz

```
## 安装go-callvis

```go
go get -u github.com/ofabry/go-callvis
cd $GOPATH/src/github.com/ofabry/go-callvis && make install
```
### 用法

``` shell
$ go-callvis [flags] package
```

### 示例

**以orchestrator项目为例，其代码已经下载到本地。**

```shell
// go-callvis main.go

go-callvis github.com/github/orchestrator/go/cmd/orchestrator
```
如果没有focus标识，默认是main

例如，查看`github.com/github/orchestrator/go/http` 这个package下面的调用关系：
```shell
$ go-callvis -focus github.com/github/orchestrator/go/http  github.com/github/orchestrator/go/cmd/orchestrator
```
浏览器跳出页面http://localhost:7878，可以看到代码调用关系图。
