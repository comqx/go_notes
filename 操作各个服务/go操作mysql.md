[TOC]

​       go的标准库database/sql包提供了保证sql或者类sql数据库的通用接口，但是没有提供具体数据库的驱动。因此使用database/sql包的时候必须注入至少一个数据库驱动

# sql标准库和驱动部署与使用

[mysql驱动](https://github.com/go-sql-driver/mysql)

```go
go get -u github.com/go-sql-driver/mysql

func Open(driverName, dataSourceName string) (*DB, error) //第一个参数指定了数据库类型，第二个指定数据库的连接方式

//创建数据库
CREATE DATABASE sql_test;
use sql_test；
//创建表
CREATE TABLE `user` (
    `id` BIGINT(20) NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(20) DEFAULT '',
    `age` INT(11) DEFAULT '0',
    PRIMARY KEY(`id`)
)ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;
```

## 初始化连接

### SetMaxOpenConns

### SetMaxIdleConns

```go
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" //使用包的init()方法
)

var db *sql.DB //设置一个全局的连接池

// 初始化数据库,建立数据库连接,指定一个全局的连接池
func initdb() {
	// 数据库信息
	//用户名:密码@tcp(ip:端口)/数据库的名称
	dsn := "root:123456@tcp(127.0.0.1:3306)/goday10"
  
  // Open打开一个dirverName指定的数据库，dataSourceName指定数据源，一般包至少括数据库文件名和（可能的）连接信息。
	// 连接数据库
	db, err := sql.Open("mysql", dsn) //如果db是一个全局变量,那么这里就不能重复声明了
	if err != nil {
		fmt.Printf("dsn:%s err:%v\n", dsn, err)
	}
	err = db.Ping() //尝试连接校验用户名密码
	if err != nil {
		fmt.Printf("open %s failed, err:%s\n", dsn, err)
	}	
  db.SetMaxOpenConns(10) //设置数据库连接池的最大连接数
	db.SetMaxIdleConns(5) //设置最大空闲连接数
  
	fmt.Println("连接数据库成功")
}
```

## 增删改查

```go
type user struct {
	id   int
	name string
	age  int
}

// 查询单个记录
func queryOne(id int) {
	var u1 user
	// 1. 写查询单条记录的sql语句
	sqlStr := "select id,name,age from user where id=?;"
	// 2. 执行并拿到结果
	// 必须对rowObj对象调用Scan方法,因为该方法会释放数据库链接 // 从连接池里拿一个连接出来去数据库查询单条记录
	db.QueryRow(sqlStr, id).Scan(&u1.id, &u1.name, &u1.age)
	// 3.打印
	fmt.Printf("u1:%#v\n", u1)
}

// 查询多个记录
func queryMore(n int) {
	var u1 user
	// 1. 定义sql
	sqlStr := "select id,name,age from user where id>?;"
	// 2. 执行sql
	rows, err := db.Query(sqlStr, n)
	if err != nil {
		fmt.Printf("exec %s query faild, err:%v\n", sqlStr, err)
	}
	// 3. 一定要关闭rows
	defer rows.Close()
	// 4. 循环取值
	for rows.Next() {
		err := rows.Scan(&u1.id, &u1.name, &u1.age)
		if err != nil {
			fmt.Printf("scan faild,err:%s\n", err)
		}
		fmt.Printf("u1:%#v\n", u1)
	}
}

// 增加记录
func insertRow() {
	// 1. sql
	sqlStr := `insert into user(name,age) values("lll",20);`
	// 2. exec
	ret, err := db.Exec(sqlStr)
	if err != nil {
		fmt.Printf("insert faild,err:%v\n", err)
	}
	// 3. 如果是插入数据的操作，能够拿到插入数据的id
	id, err := ret.LastInsertId()
	if err != nil {
		fmt.Printf("get id failed,err:%v\n", err)
		return
	}
	fmt.Println("id:", id)
}

// 修改记录
func updateRow(newage, id int) {
	sqlStr := `update user set age=? where id =?;`
	ret, err := db.Exec(sqlStr, newage, id)
	if err != nil {
		fmt.Printf("update failed, err:%v\n", err)
		return
	}
	//获取影响的行数
	n, err := ret.RowsAffected()
	if err != nil {
		fmt.Printf("get id failed,err:%v\n", err)
		return
	}
	fmt.Printf("更新了%d行数据\n", n)
}

// 删除记录
func deleteRow(id int) {
	sqlStr := `delete from user where id=?`
	ret, err := db.Exec(sqlStr, id)
	if err != nil {
		fmt.Printf("delete failed,err:%v\n", err)
		return
	}
	n, err := ret.RowsAffected()
	if err != nil {
		fmt.Printf("get id failed,err:%v\n", err)
		return
	}
	fmt.Printf("删除了%d行数据\n", n)
}
```

## mysql预处理

### 什么是预处理

普通SQL语句执行过程：

1. 客户端对SQL语句进行占位符替换得到完整的SQL语句。
2. 客户端发送完整SQL语句到MySQL服务端
3. MySQL服务端执行完整的SQL语句并将结果返回给客户端。

预处理执行过程：

1. 把SQL语句分成两部分，命令部分与数据部分。
2. 先把命令部分发送给MySQL服务端，MySQL服务端进行SQL预处理。
3. 然后把数据部分发送给MySQL服务端，MySQL服务端对SQL语句进行占位符替换。
4. MySQL服务端执行完整的SQL语句并将结果返回给客户端。

### 为什么需要预处理

1. 优化MySQL服务器重复执行SQL的方法，可以提升服务器性能，提前让服务器编译，一次编译多次执行，节省后续编译的成本。
2. 避免SQL注入问题。

```GO
// 使用预处理插入多条数据
func prepareInsert() {
	sqlStr := `insert into user(name,age) values(?,?);`
	stmt, err := db.Prepare(sqlStr) // 把SQL语句先发给MySQL预处理一下
	if err != nil {
		fmt.Printf("prepare failed ,err:%v\n", err)
		return
	}
	defer stmt.Close()
	// 后续只需要拿到stmt去执行一些操作
	var m = map[string]int{
		"yw":  20,
		"lqx": 21,
		"xk":  33,
	}
	for k, v := range m {
		stmt.Exec(k, v)
	}
}
```

## mysql事务

### 什么是事务（仅innodb支持）

事务：一个最小的不可再分的工作单元；通常一个事务对应一个完整的业务(例如银行账户转账业务，该业务就是一个最小的工作单元)，同时这个完整的业务需要执行多次的DML(insert、update、delete)语句共同联合完成。A转账给B，这里面就需要执行两次update操作。

### 事务的ACID

|  条件  |                             解释                             |
| :----: | :----------------------------------------------------------: |
| 原子性 | 一个事务（transaction）中的所有操作，要么全部完成，要么全部不完成，不会结束在中间某个环节。事务在执行过程中发生错误，会被回滚（Rollback）到事务开始前的状态，就像这个事务从来没有执行过一样。 |
| 一致性 | 在事务开始之前和事务结束以后，数据库的完整性没有被破坏。这表示写入的资料必须完全符合所有的预设规则，这包含资料的精确度、串联性以及后续数据库可以自发性地完成预定的工作。 |
| 隔离性 | 数据库允许多个并发事务同时对其数据进行读写和修改的能力，隔离性可以防止多个事务并发执行时由于交叉执行而导致数据的不一致。事务隔离分为不同级别，包括读未提交（Read uncommitted）、读提交（read committed）、可重复读（repeatable read）和串行化（Serializable）。 |
| 持久性 | 事务处理结束后，对数据的修改就是永久的，即便系统故障也不会丢失。 |

```GO
// mysql事务
func transactionDemo() {
	// 1. 开启事务
	tx, err := db.Begin()
	if err != nil {
		fmt.Printf("begin failed,err:%v\n", err)
		return
	}
	// 2. 执行多条sql
	sqlStr1 := `update user set age=age-2 where id=3`
	sqlStr2 := `update user set age=age+2 where id=5`
	// 3. 执行sql1
	_, err = tx.Exec(sqlStr1)
	if err != nil {
		// 发生回滚
		tx.Rollback()
		fmt.Println("执行sql1出错,要回滚")
		return
	}
	// 4. 执行sql2
	_, err = tx.Exec(sqlStr2)
	if err != nil {
		// 发生回滚
		tx.Rollback()
		fmt.Println("执行sql2出错,要回滚")
		return
	}
	// 5. 都返回成功,就提交本次事务
	// 上面两步SQL都执行成功，就提交本次事务
	err = tx.Commit()
	if err != nil {
		// 要回滚
		tx.Rollback()
		fmt.Println("提交出错啦，要回滚！")
		return
	}
	fmt.Println("事务执行成功！")
}
```

# sqlx第三方库部署使用

## 初始化

## 增删改查

## 事务操作

```GO
package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql" // init()
	"github.com/jmoiron/sqlx"
)

// Go连接MySQL示例
var db *sqlx.DB // 是一个连接池对象

func initDB() (err error) {
	// 数据库信息
	// 用户名:密码@tcp(ip:端口)/数据库的名字
	dsn := "root:root@tcp(127.0.0.1:3306)/sql_test"
	// 连接数据库
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		return
	}
	db.SetMaxOpenConns(10) // 设置数据库连接池的最大连接数
	db.SetMaxIdleConns(5)  // 设置最大空闲连接数
	return
}

type user struct {
	ID   int
	Name string
	Age  int
}

// 插入数据
func queryRow() {
	// 查询单条记录
	sqlStr1 := `select id, name, age from user where id=1`
	var u user
	db.Get(&u, sqlStr1)
	fmt.Printf("u:%#v\n", u)
	// 查询多条记录
	var userList []user
	sqlStr2 := `select id,name, age from user`
	err := db.Select(&userList, sqlStr2)
	if err != nil {
		fmt.Printf("select failed, err:%v\n", err)
		return
	}
	fmt.Printf("userList:%#v\n", userList)
}

// 插入数据
func insertRowDemo() {
	sqlStr := "insert into user(name, age) values (?,?)"
	ret, err := db.Exec(sqlStr, "沙河小王子", 19)
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return
	}
	theID, err := ret.LastInsertId() // 新插入数据的id
	if err != nil {
		fmt.Printf("get lastinsert ID failed, err:%v\n", err)
		return
	}
	fmt.Printf("insert success, the id is %d.\n", theID)
}

// 更新数据
func updateRowDemo() {
	sqlStr := "update user set age=? where id = ?"
	ret, err := db.Exec(sqlStr, 39, 6)
	if err != nil {
		fmt.Printf("update failed, err:%v\n", err)
		return
	}
	n, err := ret.RowsAffected() // 操作影响的行数
	if err != nil {
		fmt.Printf("get RowsAffected failed, err:%v\n", err)
		return
	}
	fmt.Printf("update success, affected rows:%d\n", n)
}

// 删除数据
func deleteRowDemo() {
	sqlStr := "delete from user where id = ?"
	ret, err := db.Exec(sqlStr, 6)
	if err != nil {
		fmt.Printf("delete failed, err:%v\n", err)
		return
	}
	n, err := ret.RowsAffected() // 操作影响的行数
	if err != nil {
		fmt.Printf("get RowsAffected failed, err:%v\n", err)
		return
	}
	fmt.Printf("delete success, affected rows:%d\n", n)
}

// 事务操作
func transactionDemo() {
	tx, err := db.Beginx() // 开启事务
	if err != nil {
		if tx != nil {
			tx.Rollback()
		}
		fmt.Printf("begin trans failed, err:%v\n", err)
		return
	}
	sqlStr1 := "Update user set age=40 where id=?"
	tx.MustExec(sqlStr1, 2)
	sqlStr2 := "Update user set age=50 where id=?"
	tx.MustExec(sqlStr2, 4)
	err = tx.Commit() // 提交事务
	if err != nil {
		tx.Rollback() // 回滚
		fmt.Printf("commit failed, err:%v\n", err)
		return
	}
	fmt.Println("exec trans success!")
}
func main() {
	err := initDB()
	if err != nil {
		fmt.Printf("init DB failed, err:%v\n", err)
	}
}
```



# 注意事项

## sql的占位符

不同数据库占位符是不同的

|   数据库   |  占位符语法  |
| :--------: | :----------: |
|   MySQL    |     `?`      |
| PostgreSQL | `$1`, `$2`等 |
|   SQLite   |  `?` 和`$1`  |
|   Oracle   |   `:name`    |

