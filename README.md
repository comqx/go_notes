


<!-- TOC -->

- [golang语言特性](#golang语言特性)
- [包的概念](#包的概念)
- [go程序目录结构](#go程序目录结构)
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

# 学习golang
```
ctrl+alt+m toc
ctrl+alt+x 粘贴image
```