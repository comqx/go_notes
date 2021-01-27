# 安装

`go get -u github.com/jteeuwen/go-bindata/...`

# 使用

>指定生成的代码的文件名`cmd/pipelineService/pipelineBindata.go`，以及指定静态目录的位置`static`

```go
# 生成代码文件、包名需要改下，默认是main
go-bindata -o=asset/asset.go -pkg asset  static/...


# 代码里面这么使用文件，解析文件的byte内容
import (	"cronK8sService/asset")
fileByte, err := asset.Asset(filePath)
```

