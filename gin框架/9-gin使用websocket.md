[toc]

# gin使用websocket

Gin 框架默认不支持 websocket，可以使用 `github.com/gorilla/websocket` 实现。

Talk is cheap. Show me the code，代码如下：

项目布局：

```shell
github.com
└── leffss
    └── ginWebsocket
        ├── go.mod
        ├── go.sum
        ├── main.go
        └── ws
            └── ws.go
```

具体原理就不讲了，可以看代码注释，比较详细了。

## `ws.go`

```go
package ws

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

// Manager 所有 websocket 信息
type Manager struct {
	Group map[string]map[string]*Client
	groupCount, clientCount uint
	Lock sync.Mutex
	Register, UnRegister chan *Client
	Message	chan *MessageData
	GroupMessage chan *GroupMessageData
	BroadCastMessage chan *BroadCastMessageData
}

// Client 单个 websocket 信息
type Client struct {
	Id, Group string
	Socket *websocket.Conn
	Message   chan []byte
}

// messageData 单个发送数据信息
type MessageData struct {
	Id, Group string
	Message    []byte
}

// groupMessageData 组广播数据信息
type GroupMessageData struct {
	Group string
	Message    []byte
}

// 广播发送数据信息
type BroadCastMessageData struct {
	Message    []byte
}

// 读信息，从 websocket 连接直接读取数据
func (c *Client) Read() {
	defer func() {
		WebsocketManager.UnRegister <- c
		log.Printf("client [%s] disconnect", c.Id)
		if err := c.Socket.Close(); err != nil {
			log.Printf("client [%s] disconnect err: %s", c.Id, err)
		}
	}()

	for {
		messageType, message, err := c.Socket.ReadMessage()
		if err != nil || messageType == websocket.CloseMessage {
			break
		}
		log.Printf("client [%s] receive message: %s", c.Id, string(message))
		c.Message <- message
	}
}

// 写信息，从 channel 变量 Send 中读取数据写入 websocket 连接
func (c *Client) Write() {
	defer func() {
		log.Printf("client [%s] disconnect", c.Id)
		if err := c.Socket.Close(); err != nil {
			log.Printf("client [%s] disconnect err: %s", c.Id, err)
		}
	}()

	for {
		select {
		case message, ok := <-c.Message:
			if !ok {
				_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			log.Printf("client [%s] write message: %s", c.Id, string(message))
			err := c.Socket.WriteMessage(websocket.BinaryMessage, message)
			if err != nil {
				log.Printf("client [%s] writemessage err: %s", c.Id, err)
			}
		}
	}
}

// 启动 websocket 管理器
func (manager *Manager) Start() {
	log.Printf("websocket manage start")
	for {
		select {
		// 注册
		case client := <-manager.Register:
			log.Printf("client [%s] connect", client.Id)
			log.Printf("register client [%s] to group [%s]", client.Id, client.Group)

			manager.Lock.Lock()
			if manager.Group[client.Group] == nil {
				manager.Group[client.Group] = make(map[string]*Client)
				manager.groupCount += 1
			}
			manager.Group[client.Group][client.Id] = client
			manager.clientCount += 1
			manager.Lock.Unlock()

		// 注销
		case client := <-manager.UnRegister:
			log.Printf("unregister client [%s] from group [%s]", client.Id, client.Group)
			manager.Lock.Lock()
			if _, ok := manager.Group[client.Group]; ok {
				if _, ok := manager.Group[client.Group][client.Id]; ok {
					close(client.Message)
					delete(manager.Group[client.Group], client.Id)
					manager.clientCount -= 1
					if len(manager.Group[client.Group]) == 0 {
						//log.Printf("delete empty group [%s]", client.Group)
						delete(manager.Group, client.Group)
						manager.groupCount -= 1
					}
				}
			}
			manager.Lock.Unlock()

			// 发送广播数据到某个组的 channel 变量 Send 中
			//case data := <-manager.boardCast:
			//	if groupMap, ok := manager.wsGroup[data.GroupId]; ok {
			//		for _, conn := range groupMap {
			//			conn.Send <- data.Data
			//		}
			//	}
		}
	}
}

// 处理单个 client 发送数据
func (manager *Manager) SendService() {
	for {
		select {
		case data := <-manager.Message:
			if groupMap, ok := manager.Group[data.Group]; ok {
				if conn, ok := groupMap[data.Id]; ok {
					conn.Message <- data.Message
				}
			}
		}
	}
}

// 处理 group 广播数据
func (manager *Manager) SendGroupService() {
	for {
		select {
		// 发送广播数据到某个组的 channel 变量 Send 中
		case data := <-manager.GroupMessage:
			if groupMap, ok := manager.Group[data.Group]; ok {
				for _, conn := range groupMap {
					conn.Message <- data.Message
				}
			}
		}
	}
}

// 处理广播数据
func (manager *Manager) SendAllService() {
	for {
		select {
		case data := <-manager.BroadCastMessage:
			for _, v := range manager.Group {
				for _, conn := range v {
					conn.Message <- data.Message
				}
			}
		}
	}
}

// 向指定的 client 发送数据
func (manager *Manager) Send(id string, group string, message []byte) {
	data := &MessageData{
		Id: id,
		Group: group,
		Message:    message,
	}
	manager.Message <- data
}

// 向指定的 Group 广播
func (manager *Manager) SendGroup(group string, message []byte) {
	data := &GroupMessageData{
		Group: group,
		Message:    message,
	}
	manager.GroupMessage <- data
}

// 广播
func (manager *Manager) SendAll(message []byte) {
	data := &BroadCastMessageData{
		Message:    message,
	}
	manager.BroadCastMessage <- data
}

// 注册
func (manager *Manager) RegisterClient(client *Client) {
	manager.Register <- client
}

// 注销
func (manager *Manager) UnRegisterClient(client *Client) {
	manager.UnRegister <- client
}

// 当前组个数
func (manager *Manager) LenGroup() uint {
	return manager.groupCount
}

// 当前连接个数
func (manager *Manager) LenClient() uint {
	return manager.clientCount
}

// 获取 wsManager 管理器信息
func (manager *Manager) Info() map[string]interface{} {
	managerInfo := make(map[string]interface{})
	managerInfo["groupLen"] = manager.LenGroup()
	managerInfo["clientLen"] = manager.LenClient()
	managerInfo["chanRegisterLen"] = len(manager.Register)
	managerInfo["chanUnregisterLen"] = len(manager.UnRegister)
	managerInfo["chanMessageLen"] = len(manager.Message)
	managerInfo["chanGroupMessageLen"] = len(manager.GroupMessage)
	managerInfo["chanBroadCastMessageLen"] = len(manager.BroadCastMessage)
	return managerInfo
}

// 初始化 wsManager 管理器
var WebsocketManager = Manager{
	Group: make(map[string]map[string]*Client),
	Register:    make(chan *Client, 128),
	UnRegister:  make(chan *Client, 128),
	GroupMessage:   make(chan *GroupMessageData, 128),
	Message:   make(chan *MessageData, 128),
	BroadCastMessage: make(chan *BroadCastMessageData, 128),
	groupCount: 0,
	clientCount: 0,
}

// gin 处理 websocket handler
func (manager *Manager) WsClient(ctx *gin.Context) {
	upGrader := websocket.Upgrader{
		// cross origin domain
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		// 处理 Sec-WebSocket-Protocol Header
		Subprotocols: []string{ctx.GetHeader("Sec-WebSocket-Protocol")},
	}

	conn, err := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("websocket connect error: %s", ctx.Param("channel"))
		return
	}

	client := &Client{
		Id:     uuid.NewV4().String(),
		Group:  ctx.Param("channel"),
		Socket: conn,
		Message:   make(chan []byte, 1024),
	}

	manager.RegisterClient(client)
	go client.Read()
	go client.Write()
	time.Sleep(time.Second * 15)
	// 测试单个 client 发送数据
	manager.Send(client.Id, client.Group, []byte("Send message ----" + time.Now().Format("2006-01-02 15:04:05")))
}

// 测试组广播
func TestSendGroup() {
	for {
		time.Sleep(time.Second * 20)
		WebsocketManager.SendGroup("leffss", []byte("SendGroup message ----" + time.Now().Format("2006-01-02 15:04:05")))
	}
}

// 测试广播
func TestSendAll() {
	for {
		time.Sleep(time.Second * 25)
		WebsocketManager.SendAll([]byte("SendAll message ----" + time.Now().Format("2006-01-02 15:04:05")))
		fmt.Println(WebsocketManager.Info())
	}
}


```
## `main.go`

```go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leffss/ginWebsocket/ws"
)

func main() {
	go ws.WebsocketManager.Start()
	go ws.WebsocketManager.SendService()
	go ws.WebsocketManager.SendService()
	go ws.WebsocketManager.SendGroupService()
	go ws.WebsocketManager.SendGroupService()
	go ws.WebsocketManager.SendAllService()
	go ws.WebsocketManager.SendAllService()
	go ws.TestSendGroup()
	go ws.TestSendAll()

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome Gin Server")
	})

	wsGroup := router.Group("/ws")
	{
		wsGroup.GET("/:channel", ws.WebsocketManager.WsClient)
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server Start Error: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown Error:", err)
	}
	log.Println("Server Shutdown")
}
```

## 测试 websocket 代码 `main.go` 

```go
package main

import (
    "flag"
    "fmt"
    "net/url"
    "time"

    "github.com/gorilla/websocket"
)

var addr = flag.String("addr", "127.0.0.1:8080", "http service address")

func main() {
    u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws/leffss"}
    var dialer *websocket.Dialer

    conn, _, err := dialer.Dial(u.String(), nil)
    if err != nil {
        fmt.Println(err)
        return
    }

    go timeWriter(conn)

    for {
        _, message, err := conn.ReadMessage()
        if err != nil {
            fmt.Println("read:", err)
            return
        }

        fmt.Printf("received: %s\n", message)
    }
}

func timeWriter(conn *websocket.Conn) {
    for {
        time.Sleep(time.Second * 5)
        conn.WriteMessage(websocket.TextMessage, []byte(time.Now().Format("2006-01-02 15:04:05")))
    }
}
```

# web页面测试ws协议

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no">
    <meta name="renderer" content="webkit">
    <title>websocket test</title>
</head>
<body>
<script src="http://cdn.bootcss.com/stomp.js/2.3.3/stomp.min.js"></script>
<script src="https://cdn.bootcss.com/sockjs-client/1.1.4/sockjs.min.js"></script>
<script type="text/javascript">
    // 建立连接对象（还未发起连接）
    var socket = new SockJS("http://geip-infra-bi-inner-proxy.glodon.com/base/ws");

    // 获取 STOMP 子协议的客户端对象
    var stompClient = Stomp.over(socket);

   var dateFormat= function(fmt, date) {
    let ret;
    const opt = {
        "Y+": date.getFullYear().toString(),        // 年
        "m+": (date.getMonth() + 1).toString(),     // 月
        "d+": date.getDate().toString(),            // 日
        "H+": date.getHours().toString(),           // 时
        "M+": date.getMinutes().toString(),         // 分
        "S+": date.getSeconds().toString()          // 秒
        // 有其他格式化字符需求可以继续添加，必须转化成字符串
    };
    for (let k in opt) {
        ret = new RegExp("(" + k + ")").exec(fmt);
        if (ret) {
            fmt = fmt.replace(ret[1], (ret[1].length == 1) ? (opt[k]) : (opt[k].padStart(ret[1].length, "0")))
        };
    };
    return fmt;
}
    // 向服务器发起websocket连接并发送CONNECT帧
    stompClient.connect(
        {'Authorization': ''},
        function connectCallback(frame) {
            // 连接成功时（服务器响应 CONNECTED 帧）的回调方法
	    let date = new Date();
	    let fdate = dateFormat("YYYY-mm-dd HH:MM:SS", date);
            setMessageInnerHTML("connect success.<br>connection time:"+fdate);
            stompClient.subscribe('/topic/348394616033792', function (response) {
				if(null!= response){
				  setMessageInnerHTML("348394616033792 raw data:"+ response.body);
				  setMessageInnerHTML("<br>");
				}
            });
        },
        function errorCallBack(error) {
            // 连接失败时（服务器响应 ERROR 帧）的回调方法
 	    let date = new Date();
	    let fdate = dateFormat("YYYY-mm-dd HH:MM:SS", date);
            setMessageInnerHTML("connection failed:"+ error+ "<br>error time:"+ fdate);
        }
    );
    //将消息显示在网页上
    function setMessageInnerHTML(innerHTML) {
         document.getElementById('message').innerHTML += innerHTML + '<br/>';
    }
</script>
<div id="message"></div>
</body>
</html>
```

