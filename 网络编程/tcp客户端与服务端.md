# TCP客户端和服务端

## 服务端

```go
package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func processCon(conn net.Conn) {
	defer conn.Close()
	// 3. 与客户端通信
	var tmp [128]byte
	reader := bufio.NewReader(os.Stdin) //抓取终端输入
	for {
		n, err := conn.Read(tmp[:]) //传入一个固定大小的数组，然后去获取内容，n表示获取的内容长度
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("读取失败：", err)
		}
		fmt.Println(string(tmp[:n])) //转换为字符串
		fmt.Println("请回复：")
		msg, _ := reader.ReadString('\n') //获取输入信息按照\n分割
		msg = strings.TrimSpace(msg)
		if msg == "exit" {
			break
		}
		conn.Write([]byte(msg)) //写入conn里面，也就是发送给对端
	}
}
func runNet() {
	var tcpAddr = net.TCPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 8080,
	}

	// 1.本地端口启动
	tcpListen, err := net.ListenTCP("tcp", &tcpAddr)
	if err != nil {
		fmt.Println("监听端口err:", err)
	}
	defer tcpListen.Close()
	// 2.等待别人跟我连接
	for {
		conn, err := tcpListen.Accept() //获取连接的对象
		if err != nil {
			fmt.Println("连接失败: ", err)
		}
		go processCon(conn) //开启线程去接收和发送信息
	}
}

func main() {
	runNet()
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
	"strings"
)

func runConn() {
	// 1.连接对端端口
	conn, err := net.Dial("tcp", "127.0.0.1:8080") //建立连接
	if err != nil {
		fmt.Println("监听端口err:", err)
	}
	defer conn.Close()
	// 2.发送消息，然后循环读信息
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("请说话：")
		msg, _ := reader.ReadString('\n') //获取输入信息
		msg = strings.TrimSpace(msg)
		if msg == "exit" {
			break
		}
		conn.Write([]byte(msg)) //写入输入信息
	}
}

func main() {
	runConn()
}
```



# TCP黏包解决

### 服务端

```go
package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	proto "oldbody.com/day8/07nianbao_jiejue/protocol"
)

func process(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		recvStr, err := proto.Decode(reader) //直接调用这个包来解决
		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Println("decode failed,err:", err)
			return
		}
		fmt.Println("收到client发来的数据：", recvStr)
	}
}

func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:30000")
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("accept failed, err:", err)
			continue
		}
		go process(conn)
	}
}
```

### 客户端

```go
package main

import (
	"fmt"
	"net"

	proto "oldbody.com/day8/07nianbao_jiejue/protocol"
)

// 黏包 client
func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:30000")
	if err != nil {
		fmt.Println("dial failed, err", err)
		return
	}
	defer conn.Close()
	for i := 0; i < 20; i++ {
		msg := `Hello, Hello. How are you?`
		// 调用协议编码数据
		b, err := proto.Encode(msg)
		if err != nil {
			fmt.Println("encode failed,err:", err)
			return
		}
		conn.Write(b)
		// time.Sleep(time.Second)
	}
}
```

### 协议

```go
package proto

import (
	"bufio"
	"bytes"
	"encoding/binary"
)

// Encode 将消息编码
func Encode(message string) ([]byte, error) {
	// 读取消息的长度，转换成int32类型（占4个字节）
	var length = int32(len(message))
	var pkg = new(bytes.Buffer)
	// 写入消息头
	err := binary.Write(pkg, binary.LittleEndian, length)
	if err != nil {
		return nil, err
	}
	// 写入消息实体
	err = binary.Write(pkg, binary.LittleEndian, []byte(message))
	if err != nil {
		return nil, err
	}
	return pkg.Bytes(), nil
}

// Decode 解码消息
func Decode(reader *bufio.Reader) (string, error) {
	// 读取消息的长度
	lengthByte, _ := reader.Peek(4) // 读取前4个字节的数据
	lengthBuff := bytes.NewBuffer(lengthByte)
	var length int32
	err := binary.Read(lengthBuff, binary.LittleEndian, &length)
	if err != nil {
		return "", err
	}
	// Buffered返回缓冲中现有的可读取的字节数。
	if int32(reader.Buffered()) < length+4 {
		return "", err
	}

	// 读取真正的消息数据
	pack := make([]byte, int(4+length))
	_, err = reader.Read(pack)
	if err != nil {
		return "", err
	}
	return string(pack[4:]), nil
}
```

