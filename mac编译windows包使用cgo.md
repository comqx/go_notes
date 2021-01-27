> 主要问题是go里面使用了sqlite3包导致需要使用cgo来编译

# runtime/cgo

gcc_libinit_windows.c:7:10: fatal error: 'windows.h' file not found

- 安装MacPorts
  https://www.macports.org/install.php#installing
  从上述连接选择对应的系统版本,下载对应的pkg文件
  双击安装，一路next加同意
- 安装使用
  http://mingw-w64.org/doku.php/download/macports
  `port not found!`
  解决方案：find / -name "port" 有时候会出现在：`/opt/local/bin/port`下面
- 执行
  `sudo /opt/local/bin/port install mingw-w64`
  等待完成安装
- 安装结果会在以下路径产生程序
  `/opt/local/bin/x86_64-w64-mingw32-gcc`

以上是使用了c的编译方式进行交叉不同系统之间的编译

以上是实例：

```shell
CGO_ENABLED=0 GOOS=windows GOARCH=amd64  go build -o web.exe -ldflags "-H windowsgui" web.go 

# 有时候会报  undefined: SQLiteConn

# 接下来使用
CGO_ENABLED=1 GOOS=windows GOARCH=amd64  go build -o web.exe -ldflags "-H windowsgui" web.go 

# 会提示
# runtime/cgo
gcc_libinit_windows.c:7:10: fatal error: 'windows.h' file not found

# 这时候就用到了上面安装的MacPorts了

CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=/opt/local/bin/x86_64-w64-mingw32-gcc go build -o web.exe -ldflags "-H windowsgui" web.go 
```

- web.go是包含sqlite3的驱动的，这个是在mac下进行编译windwos碰到的问题。主要是使用了交叉C编译器代替主机编译器。

https://github.com/mattn/go-sqlite3/issues/444
https://github.com/mattn/go-sqlite3/issues/372

github有人提问，解决方案就是使用c编译器代替主机编译器