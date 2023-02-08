```golang
package main

import (
   "bytes"
   "fmt"
   "io/ioutil"
   "net/http"
)

func dockerToApi() {
	url := "http://10.20.3.102:2375/build?t=golang_image"
	// dockerTestBuild.tar的内容需要是同级目录的内容
	// ./docker/
	// ./abcd/
	// ./Dockerfile
	content, err := ioutil.ReadFile("./dockerTestBuild.tar")
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(content))
	req.Header.Set("Content-Type", "application/x-tar")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("status", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}


作者：liedmirror
链接：https://juejin.cn/post/6992540701607067656
来源：稀土掘金
著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。
```

