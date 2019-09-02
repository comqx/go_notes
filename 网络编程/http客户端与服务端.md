# http客户端与服务端

## 服务端

```go
package main

import (
	"fmt"
	"net/http"
)

func sayHello(resp http.ResponseWriter, r *http.Request) {
	var respIntContent string
	if r.Method == "POST" {
		respIntContent = "post"
	} else if r.Method == "GET" {
		respIntContent = "get"
	}
	fmt.Fprint(resp, respIntContent) //把respIntContent写入resp这个对象中，返回给http前端
}

func main() {
	http.HandleFunc("/hello", sayHello) //指定url路径以及执行的函数
	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		fmt.Println("start http server failed ,err", err)
	}
}
```

## 客户端

```go
func main() {
	//get 请求
	resp, err := http.Get("http://qixiang-liu.github.io")
	if err != nil {
		fmt.Println("get url failed err", err)
	}
	defer resp.Body.Close()
	//读取body
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("readall body err ", err)
	}
	fmt.Println(string(b))
}
```

## 爬虫实例

```go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

//定义客户端，然后关闭keepalive
var client = http.Client{
	Transport: &http.Transport{
		DisableKeepAlives: true, //当需要频繁发生请求的时候，需要把keepalive关闭掉。避免无限建立连接。 
	},
}

//定义需要接收的字段数据的一个bewlist类型的 struct
type bwcList struct {
	Data struct {
		Detail []struct {
			OfflineActivityId int `json:"offlineActivityId"`
		} `json:"detail"`
	} `json:"data"`
}

// 判断平台打开，服务启动就自动打开图片
func openbrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func login() http.Cookie {
	Qurl := "https://account.dianping.com/account/getqrcodeimg" //定义url
	var cookies http.Cookie  
	req, err := http.NewRequest("GET", Qurl, nil) //创建一个请求，是一个get请求，Qurl，需要传入的参数是nil
	if err != nil {
		fmt.Println("newrequest err:", err)
	}
  // 定义get请求的header
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.86 Safari/537.36")
	req.Header.Add("Referer", "https://account.dianping.com/account/iframeLogin?callback=EasyLogin_frame_callback0&wide=false&protocol=https:&redir=http://www.dianping.com/")
	req.Header.Add("Host", "account.dianping.com")
  //获取respoes
	res, _ := client.Do(req)
  //获取返回的内容
	resb, _ := ioutil.ReadAll(res.Body)
	newFile, _ := os.OpenFile("ewm.jpg", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	newFile.Write(resb) //因为是一个图片，需要写入byte类型的切片
	defer func() {
		newFile.Close()
		res.Body.Close()
	}()
	openbrowser("./ewm.jpg")
	//获取cookies
	cookies.Value = res.Cookies()[0].Value
	cookies.Name = res.Cookies()[0].Name
	return cookies
}

//拿这二维码的cook，去循环获取认证成功后的cook
func login_sm(cookies http.Cookie) []*http.Cookie {
	Qurl := "https://account.dianping.com/account/ajax/queryqrcodestatus"
	cookieFormat := fmt.Sprintf("%s=%s", cookies.Name, cookies.Value)
	var cookieSlice []*http.Cookie
	fmt.Println(cookieFormat)
	for {
		// for {
		req, err := http.NewRequest("POST", Qurl, strings.NewReader(cookieFormat)) //post 请求，请求的url是Qurl，携带的数据是strings.NewReader(cookieFormat)
		if err != nil {
			fmt.Println("get err :", err)
		}
		//post请求，必须要设定Content-Type为application/x-www-form-urlencoded，post参数才可正常传递。
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded") //表示传入的参数可以是一个a=1&b=2的类型
		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.86 Safari/537.36")
		req.Header.Add("Referer", "https://account.dianping.com/account/iframeLogin?callback=EasyLogin_frame_callback0&wide=false&protocol=https:&redir=http://www.dianping.com/")
		req.Header.Add("Origin", "https://account.dianping.com")
		req.Header.Add("X-Requested-With", "XMLHttpRequest")
		req.AddCookie(&cookies) //添加cookies

		time.Sleep(time.Second * 2)
    //获取respoes
		r, err := client.Do(req)
		if err != nil {
			fmt.Println("do err", err)
		}
		cookieSlice = r.Cookies()
		fmt.Println(cookieSlice)
		if len(cookieSlice) != 0 {
			return cookieSlice
		}
	}
}

func bwc_list(page int) (bwcList, bool) {
	url := "http://m.dianping.com/activity/static/pc/ajaxList"
	bodyStr := fmt.Sprintf("{\"cityId\":2,\"type\":1,\"mode\":\"\",\"page\":%d}", page)
	bodyData := strings.NewReader(bodyStr)
	//#cityid:1 上海；2 北京
	//# type： 1 ，全部 2， 美食， 3 丽人 4 玩乐
	req, err := http.NewRequest("POST", url, bodyData)
	if err != nil {
		fmt.Println("newrequest err,", err)
	}
	// req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.86 Safari/537.36")
	rep, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败，err:", err)
	}
	bData, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		fmt.Println("读取body失败:", err)
	}
	var jsonStruct bwcList
	err = json.Unmarshal(bData, &jsonStruct)
	if err != nil {
		fmt.Println("解析失败：", err)
	}
	defer func() {
		rep.Body.Close()
	}()
	if len(jsonStruct.Data.Detail) == 0 {
		return jsonStruct, false
	}
	return jsonStruct, true
}

func baoming(jsonStructP bwcList, cookieSlice []*http.Cookie) {
	url := "http://s.dianping.com/ajax/json/activity/offline/saveApplyInfo"
	for _, value := range jsonStructP.Data.Detail {
		time.Sleep(time.Millisecond * 500)
		fmt.Println(value.OfflineActivityId)
		dataStr := fmt.Sprintf("offlineActivityId=%d&phoneNo=183****6263&shippingAddress=&extraCount=&birthdayStr=&email=&marryDayStr=&babyBirths=&pregnant=&marryStatus=0&comboId=&branchId=&usePassCard=0&passCardNo=&isShareSina=false&isShareQQ=false", value.OfflineActivityId)
		Referer := fmt.Sprintf("http://s.dianping.com/event/%d", value.OfflineActivityId)
		req, _ := http.NewRequest("POST", url, strings.NewReader(dataStr))
		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.86 Safari/537.36")
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8;")
		req.Header.Add("Referer", Referer)
		req.Header.Add("Origin", "http://s.dianping.com")
		req.Header.Add("X-Request", "JSON")
		req.Header.Add("X-Requested-With", "XMLHttpRequest")
		req.Header.Add("Accept", "application/json,text/javascript")
		for _, v := range cookieSlice {
			req.AddCookie(v)
		}
    
		res, err := client.Do(req)
		if err != nil {
			fmt.Println("请求失败")
		}
		resB, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(resB))
		defer func() {
			res.Body.Close()
		}()
	}
}

func main() {
	cook := login()
	cookieSlice := login_sm(cook)
	cookieSlice = append(cookieSlice, &cook)
	for num := 1; num <= 4; num++ {
		jsonStructP, ok := bwc_list(num)
		if !ok {
			break
		}
		baoming(jsonStructP, cookieSlice)
	}
}
```



