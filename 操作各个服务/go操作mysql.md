[TOC]

​       go的标准库database/sql包提供了保证sql或者类sql数据库的通用接口，但是没有提供具体数据库的驱动。因此使用database/sql包的时候必须注入至少一个数据库驱动

# sql标准库和驱动部署与使用

[mysql驱动](https://github.com/go-sql-driver/mysql)

```go
go get -u github.com/go-sql-driver/mysql

func Open(driverName, dataSourceName string) (*DB, error) //第一个参数指定了数据库类型，第二个指定数据库的连接方式
```

```go
import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// Open打开一个dirverName指定的数据库，dataSourceName指定数据源，一般包至少括数据库文件名和（可能的）连接信息。
func main() {
   // DSN:Data Source Name
	dsn := "user:password@tcp(127.0.0.1:3306)/dbname"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()
}
```

## 初始化连接

### SetMaxOpenConns

### SetMaxIdleConns

## 增删改查

### 单行查询

### 多行查询

### 插入数据

### 更新数据

### 删除数据

## mysql预处理

## mysql事务

# sqlx第三方库部署使用

## 初始化

## 增删改查

## 事务操作

# 注意事项

## sql的占位符

## sql注入

