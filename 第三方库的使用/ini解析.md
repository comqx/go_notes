# ini使用

`gopkg.in/ini.v1`

## 基本使用

```GO
package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	//https://ini.unknwon.io/docs/advanced/map_and_reflect
	"os"
)

//基本操作
func simoper() {
	cfg, err := ini.Load("conf.ini")
	if err != nil {
		fmt.Println("load err:", err)
		os.Exit(1)
	}
	//典型读取操作，默认分区可以使用空字符串表示
	fmt.Println("app mode:", cfg.Section("").Key("app_mode").String())
	fmt.Println("data path:", cfg.Section("paths").Key("data").String())

	//我们可以做一些候选值限制的操作
	fmt.Println("server protocol:",
		cfg.Section("server").Key("protocol").In("http", []string{"http", "https"}))

	// 如果读取的值不在候选列表内，则会回退使用提供的默认值
	fmt.Println("Email Protocol:",
		cfg.Section("server").Key("protocol").In("smtp", []string{"imap", "smtp"}))

	// 试一试自动类型转换
	fmt.Printf("Port Number: (%[1]T) %[1]d\n", cfg.Section("server").Key("http_port").MustInt(9999))
	fmt.Printf("Enforce Domain: (%[1]T) %[1]v\n", cfg.Section("server").Key("enforce_domain").MustBool(false))

	// 修改某个值然后进行保存
	cfg.Section("").Key("app_mode").SetValue("production")
	cfg.SaveTo("my.ini.local")
}
```

## 使用struct接收值

```GO
//使用struct来映射
func main() {
	//节定义为一个struct类型
	type Paths struct {
		Data string `ini:"data"`
	}
	type Server struct {
		Protocol       string `ini:"protocol"`
		Http_port      int    `ini:"http_port"`
		Enforce_domain bool   `ini:"enforce_demain"`
	}
	type Conf struct {
		App_mode string `ini:"app_mode"`
		Paths    `ini:"paths"`
		Server   `ini:"server"`
	}
	// 加载配置文件
	cfg, err := ini.Load("conf.ini")
	if err != nil {
		fmt.Println("load err:", err)
		os.Exit(1)
	}
	//new  给struct初始化一块内存地址
	p := new(Conf)
	//反射给对应的值
	err = cfg.MapTo(p)
	if err != nil {
		fmt.Println("mapto err:", err)
	}
	fmt.Println(p)
}
```

