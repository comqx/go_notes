# udp客户端与服务端

## 服务端

```go
package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 40000,
	})
	if err != nil {
		fmt.Println("listen udp faild err:", err)
	}
	defer conn.Close()

	//不需要建立连接，直接收发数据
	var data [1024]byte
	for {
		n, addr, err := conn.ReadFromUDP(data[:])
		if err != nil {
			fmt.Println("read from UDP faild,err:", err)
		}
		fmt.Println(data[:n])
		reply := strings.ToUpper(string(data[:n]))
		//发送消息
		conn.WriteToUDP([]byte(reply), addr)
	}
}
```

## 客户端

```go
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

//udp client
func main() {
	socket, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 40000,
	})
	if err != nil {
		fmt.Println("连接服务端失败，err:", err)
	}
	defer socket.Close()
	var reply [1024]byte
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("请输入内容：")
		msg, _ := reader.ReadString('\n')
		socket.Write([]byte(msg))
		// 收回复的数据
		n, _, err := socket.ReadFromUDP(reply[:])
		if err != nil {
			fmt.Println("redv reply msg failed,err:", err)
			return
		}
		fmt.Println("收到回复信息：", string(reply[:n]))
	}
}
```

