```shell
go get -u github.com/jinzhu/gorm
# 官方：
https://gorm.io/zh_CN/docs/index.html
```

# gorm快速入门


```go
package main

import (
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Product struct {
  gorm.Model
  Code string
  Price uint
}

func main() {
  db, err := gorm.Open("sqlite3", "test.db")
  if err != nil {
    panic("failed to connect database")
  }
  defer db.Close()

  // Migrate the schema 迁移ｓｃｈｅｍａ
  db.AutoMigrate(&Product{})

  // 创建
  db.Create(&Product{Code: "L1212", Price: 1000})

  // 读取
  var product Product
  db.First(&product, 1) // 查询id为1的product
  db.First(&product, "code = ?", "L1212") // 查询code为l1212的product

  // 更新 - 更新product的price为2000
  db.Model(&product).Update("Price", 2000)

  // 删除 - 删除product
  db.Delete(&product)
}
```

