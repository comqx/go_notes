[toc]

# 简介

在本次课程中，我们来学习使用WebSocket来打造一个实时聊天系统。我们会从一下几个方面来进行学习：

什么是websocket；

Websocket与传统的HTTP协议有什么区别；

Websocket有哪些优点；

如何建立连接；

如何维持连接；

Golang实战项目—实时聊天系统；

# 总结；

## 什么是websocket？

WebSocket协议是基于TCP的一种新的网络协议。它实现了**浏览器与服务器全双工(full-duplex)通信**——**允许服务器主动发送信息给客户端。**

WebSocket通信协议于2011年被IETF定为标准RFC 6455，并被RFC7936所补充规范。

WebSocket协议支持（在受控环境中运行不受信任的代码的）客户端与（选择加入该代码的通信的）远程主机之间进行全双工通信。用于此的安全模型是Web浏览器常用的基于原始的安全模式。 协议包括一个开放的握手以及随后的TCP层上的消息帧。 该技术的目标是为基于浏览器的、需要和服务器进行双向通信的（服务器不能依赖于打开多个HTTP连接（例如，使用XMLHttpRequest或iframe和长轮询））应用程序提供一种通信机制。

## Websocket与传统的HTTP协议有什么区别?

**http，websocket**都是应用层协议，他们规定的是数据怎么封装，而他们传输的通道是下层提供的。就是说无论是 http 请求，还是 WebSocket 请求，他们用的连接都是**传输层**提供的，即 tcp 连接（传输层还有 udp 连接）。只是说 http1.0 协议规定，你一个请求获得一个响应后，你要把连接关掉。所以你用 http 协议发送的请求是无法做到一直连着的（如果服务器一直不返回也可以保持相当一段时间，但是也会有超时而被断掉）。而 WebSocket 协议规定说等握手完成后我们的连接不能断哈。虽然 WebSocket 握手用的是 http 请求，但是请求头和响应头里面都有特殊字段，当浏览器或者服务端收到后会做相应的协议转换。所以 http 请求被 hold 住不返回的长连接和 WebSocket 的连接是有本质区别的。

## WebSocket有哪些优点？

说到优点，这里的对比参照物是HTTP协议，概括地说就是：支持双向通信，更灵活，更高效，可扩展性更好。

**支持双向通信，实时性更强。**

**更好的二进制支持。**

**较少的控制开销。**连接创建后，ws客户端、服务端进行数据交换时，协议控制的数据包头部较小。在不包含头部的情况下，服务端到客户端的包头只有2~10字节（取决于数据包长度），客户端到服务端的的话，需要加上额外的4字节的掩码。而HTTP协议每次通信都需要携带完整的头部。

**支持扩展**。ws协议定义了扩展，用户可以扩展协议，或者实现自定义的子协议。（比如支持自定义压缩算法等）

对于后面两点，没有研究过WebSocket协议规范的同学可能理解起来不够直观，但不影响对WebSocket

## 如何建立连接？

客户端通过HTTP请求与WebSocket服务端协商升级协议。协议升级完成后，后续的数据交换则遵照WebSocket的协议。

### 1. 客户端：申请协议升级

首先，客户端发起协议升级请求。可以看到，采用的是标准的HTTP报文格式，且只支持GET方法。

```shell
GET / HTTP/1.1

Host: localhost:8080

Origin: http://127.0.0.1:3000

Connection: Upgrade  

Upgrade: websocket

Sec-WebSocket-Version: 13

Sec-WebSocket-Key: w4v7O6xFTi36lq3RNcgctw==
```

重点请求首部意义如下：

**Connection: Upgrade**：表示要升级协议

**Upgrade: websocket**：表示要升级到websocket协议。

**Sec-WebSocket-Version:** 13：表示websocket的版本。如果服务端不支持该版本，需要返回一个Sec-WebSocket-Versionheader，里面包含服务端支持的版本号。

**Sec-WebSocket-Key**：与后面服务端响应首部的Sec-WebSocket-Accept是配套的，提供基本的防护，比如恶意的连接，或者无意的连接。

### 2. 服务器：响应协议升级

服务端返回内容如下，状态代码101表示协议切换。到此完成协议升级，后续的数据交互都按照新的协议来。

```shell
HTTP/1.1 101 Switching Protocols

Connection:Upgrade

Upgrade: websocket

Sec-WebSocket-Accept: Oy4NRAQ13jhfONC7bP8dTKb4PTU=
```

###  3. Sec-WebSocket-Accept的计算

**Sec-WebSocket-Accept**根据客户端请求首部的Sec-WebSocket-Key计算出来。

计算公式为：

将Sec-WebSocket-Key跟258EAFA5-E914-47DA-95CA-C5AB0DC85B11拼接。

通过SHA1计算出摘要，并转成base64字符串。

## 如何维持连接?

WebSocket为了保持客户端、服务端的实时双向通信，需要确保客户端、服务端之间的TCP通道保持连接没有断开。然而，对于长时间没有数据往来的连接，如果依旧长时间保持着，可能会浪费包括的连接资源。

但不排除有些场景，客户端、服务端虽然长时间没有数据往来，但仍需要保持连接。这个时候，可以采用心跳来实现。

**发送方->接收方**：ping

**接收方->发送方**：pong

ping、pong的操作，对应的是WebSocket的两个控制帧，opcode分别是0x9、0xA。

举例，WebSocket服务端向客户端发送ping，只需要如下代码（采用ws模块）

Golang实战项目—实时聊天系统

这里是使用Github上的一个开源项目作为案例。

## 获取Golang的websocket库

> go get github.com/gorilla/websocket

###  获取测试程序

> git clone https://github.com/scotch-io/go-realtime-chat.git

###  Server

```go
package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan Message)           // broadcast channel
// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Define our message object
type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()
	// Register our new client
	clients[ws] = true
	for {
		var msg Message
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		// Send the newly received message to the broadcast channel
		broadcast <- msg
	}
}
func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		// Send it out to every client that is currently connected
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
func main() {
	// Create a simple file server
	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)
	// Configure websocket route
	http.HandleFunc("/ws", handleConnections)
	// Start listening for incoming chat messages
	go handleMessages()
	// Start the server on localhost port 8000 and log any errors
	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

```



## 客户端

```vue
new Vue({
el: '#app',
data: {
    ws: null, // Our websocket
    newMsg: '', // Holds new messages to be sent to the server
    chatContent: '', // A running list of chat messages displayed on the screen
    email: null, // Email address used for grabbing an avatar
    username: null, // Our username
    joined: false // True if email and username have been filled in
},
created: function() {
    var self = this;
    this.ws = new WebSocket('ws://' + window.location.host + '/ws');
    this.ws.addEventListener('message', function(e) {
      var msg = JSON.parse(e.data);
      self.chatContent += '<div class="chip">'
          \+ '<img src="' + self.gravatarURL(msg.email) + '">' // Avatar
          \+ msg.username
        \+ '</div>'
        \+ emojione.toImage(msg.message) + '<br/>'; // Parse emojis
      var element = document.getElementById('chat-messages');
      element.scrollTop = element.scrollHeight; // Auto scroll to the bottom
    });
},
methods: {
    send: function () {
      if (this.newMsg != '') {
        this.ws.send(
          JSON.stringify({
            email: this.email,
            username: this.username,
            message: $('<p>').html(this.newMsg).text() // Strip out html
          }
        ));
        this.newMsg = ''; // Reset newMsg
      }
    },
    join: function () {
      if (!this.email) {
        Materialize.toast('You must enter an email', 2000);
        return
      }
      if (!this.username) {
        Materialize.toast('You must choose a username', 2000);
        return
      }
      this.email = $('<p>').html(this.email).text();
      this.username = $('<p>').html(this.username).text();
      this.joined = true;
    },
    gravatarURL: function(email) {
      return 'http://www.gravatar.com/avatar/' + CryptoJS.MD5(email);
    }
}
});
```



# 总结

具体使用什么技术是需要根据使用场景进行选择的，在这里我为大家总结了http和websocket不同的使用场景，请大家参考。

## HTTP 

检索资源（Retrieve Resource）

高度可缓存的资源（Highly Cacheable Resource）

幂等性和安全性（Idempotency and Safety）

错误方案（Error Scenarios）

## websockt

快速响应时间（Fast Reaction Time）

持续更新（Ongoing Updates）

Ad-hoc消息传递（Ad-hoc Messaging）

错误的HTTP应用场景

依赖于客户端轮询服务，而不是由用户主动发起。

需要频繁的服务调用来发送小消息。

客户端需要快速响应对资源的更改，并且，无法预测更改何时发生。

错误的WebSockets应用场景

连接仅用于极少数事件或非常短的时间，客户端无需快速响应事件。

需要一次打开多个WebSockets到同一服务。

打开WebSocket，发送消息，然后关闭它 - 然后再重复该过程。

消息传递层中重新实现请求/响应模式。





# 使用websocket完成实时打印日志功能

## 服务端

```go
package main

import (
	"bufio"
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func echoHandler(ws *websocket.Conn) {
	var (
		err  error
		m    int
		line string
	)

	fileObj, err := os.Open("/Users/liuqixiang/project/multCloudDev/attachments/index.html")
	if err != nil {
		fmt.Printf("open file faild, err:%v\n", err)
		return
	}
	//关闭文件
	defer fileObj.Close()
	reader := bufio.NewReader(fileObj)
	msg := make([]byte, 512)
	for {

		//n, err = ws.Read(msg)
		//if err != nil {
		//	log.Fatal(err)
		//}
		//fmt.Printf("Receive: %s\n", msg[:n])

		line, err = reader.ReadString('\n') //按照字符’\n‘来分割每次读取长度
		if err == io.EOF {
			time.Sleep(10 * time.Millisecond)
			continue
		}
		if err != nil {
			fmt.Printf("read from file failed, err:%v\n", err)
			return
		}
		fmt.Print(line)

		//send_msg := "[" + string(msg[:n]) + "]"
		m, err = ws.Write([]byte(line))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Send: %s\n", msg[:m])
	}

}

func main() {
	http.Handle("/echo", websocket.Handler(echoHandler))
	http.Handle("/", http.FileServer(http.Dir("/Users/liuqixiang/project/multCloudDev/attachments/")))

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
```
## 前端页面

```html
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8"/>
    <title>Sample of websocket with golang</title>
    <script src="http://apps.bdimg.com/libs/jquery/2.1.4/jquery.min.js"></script>

    <script>
        $(function() {
            var ws = new WebSocket("ws://localhost:8080/echo");
            ws.onmessage = function(e) {
                $('<li>').text(event.data).appendTo($ul);
            };
            var $ul = $('#msg-list');
            $('#sendBtn').click(function(){
                var data = $('#name').val();
                ws.send(data);
            });
        });
    </script>
</head>
<body>
<input id="name" type="text"/>
<input type="button" id="sendBtn" value="send"/>
<ul id="msg-list"></ul>
</body>
</html>
```

