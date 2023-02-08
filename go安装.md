[toc]

# Centos 安装golang并安装使用中文文档

## 下载yum源并安装golang

```shell
rpm --import https://mirror.go-repo.io/centos-unstable/RPM-GPG-KEY-GO-REPO
curl -s https://mirror.go-repo.io/centos-unstable/go-repo-unstable.repo | tee /etc/yum.repos.d/go-repo-unstable.repo
yum install golang
```



## 安装中文版gotour（在线学习golang工具）

> 本地golang学习工具

```shell
go get github.com/Go-zh/tour/gotour

cd $GOPATH/bin
./gotour

# 浏览器访问 127.0.0.1:60712
```

