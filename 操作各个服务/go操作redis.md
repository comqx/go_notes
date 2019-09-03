[TOC]

# redis介绍

## 支持的数据结构

Redis支持诸如字符串（strings）、哈希（hashes）、列表（lists）、集合（sets）、带范围查询的排序集合（sorted sets）、位图（bitmaps）、hyperloglogs、带半径查询和流的地理空间索引等数据结构（geospatial indexes）。

## 应用场景

- 缓存系统，减轻主数据库的压力
- 计数场景，比如微薄、抖音的关注度和粉丝量
- 热门排行榜，需要排序的场景特别适合使用zset
- 利用list可以实现队列的功能

# 操作redis

## Redis第三方库安装与连接初始化

`go get -u github.com/go-redis/redis`

```GO
package main

import (
	"fmt"
	"github.com/go-redis/redis"
)

var redisDB *redis.Client

// 初始化连接
func initClient() (err error) {
	redisDB = redis.NewClient(&redis.Options{  //建立连接
		Addr:     "localhost:6379",
		Password: "123",
		DB:       0,
	})
	ok, err := redisDB.Ping().Result() //检测是否可用
	if err != nil {
		return err
	}
	fmt.Printf("redis status is %s\n", ok)
	return nil
}

```

## 基本使用

### set、get

```GO
// set,get
func redisExample() {
	//func (redis.cmdable).Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	// set 传入 key,value,存活时间(0是永久),返回查看状态的地址
	err := redisDB.Set("score", 100, 0).Err()
	if err != nil {
		fmt.Printf("set score failed, err:%v\n", err)
		return
	}
	//获取指定key的结果
	val, err := redisDB.Get("score").Result()
	if err != nil {
		fmt.Printf("get score failed, err:%v\n", err)
		return
	}
	fmt.Println("score", val)

	// 查看一个不存在的key,err返回redis.Nil
	val2, err := redisDB.Get("name").Result()
	if err == redis.Nil {
		fmt.Println("name does not exist")
	} else if err != nil {
		fmt.Printf("get name failed, err:%v\n", err)
		return
	} else {
		fmt.Println("name", val2)
	}
}
```

### zset示例

```GO
//ZSET
func redisExample2() {
	zsetKey := "language_rank"
	language := []*redis.Z{
		&redis.Z{Score: 90.0, Member: "Golang"},
		&redis.Z{Score: 98.0, Member: "Java"},
		&redis.Z{Score: 95.0, Member: "Python"},
		&redis.Z{Score: 97.0, Member: "JavaScript"},
		&redis.Z{Score: 99.0, Member: "C/C++"},
	}
	//ZADD
	num, err := redisDB.ZAdd(zsetKey, language...).Result()
	if err != nil {
		fmt.Printf("Zadd faild,err:%v\n", err)
	}
	fmt.Printf("zadd %d succ.\n", num)

	// 把golang的分数加10
	newScore, err := redisDB.ZIncrBy(zsetKey, 10.0, "Golang").Result()
	if err != nil {
		fmt.Printf("zincrby failed, err:%v\n", err)
		return
	}
	fmt.Printf("Golang's score is %f now.\n", newScore)

	//取分数最高的3位
	ret, err := redisDB.ZRevRangeWithScores(zsetKey, 0, 2).Result()
	if err != nil {
		fmt.Printf("zrevrange failed, err:%v\n", err)
		return
	}
	for _, z := range ret {
		fmt.Println("取分数最高的3位:", z.Member, z.Score)
	}

	// 取95~100分的
	op := &redis.ZRangeBy{
		Min: "95",
		Max: "100",
	}
	ret, err = redisDB.ZRangeByScoreWithScores(zsetKey, op).Result()
	if err != nil {
		fmt.Printf("zrangebyscore failed, err:%v\n", err)
		return
	}
	for _, z := range ret {
		fmt.Println("取95-100分的", z.Member, z.Score)
	}
}
```

