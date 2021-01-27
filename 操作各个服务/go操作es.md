[TOC]

> ​		前言：elasticsearch 是一个基于Lucene构建的开源的、分布式、restful接口的全文搜索引擎。es还是一个分布式的文档数据库，其中每个字段均可被索引，而且每个字段的数据均可被搜索。es扩展能力很强，可以扩展到几百台服务器以及处理PB级的数据。可以在短时间内存储、搜索和分析大量的数据。通常应用于复杂搜索场景情况下的核心发动机



# elasticsearch

快速理解es：http://developer.51cto.com/art/201904/594615.htm

## es能做什么

1. 当你经营一家网上商店，你可以让你的客户搜索你卖的商品。在这种情况下，你可以使用ElasticSearch来存储你的整个产品目录和库存信息，为客户提供精准搜索，可以为客户推荐相关商品。
2. 当你想收集日志或者交易数据的时候，需要分析和挖掘这些数据，寻找趋势，进行统计，总结，或发现异常。在这种情况下，你可以使用Logstash或者其他工具来进行收集数据，当这引起数据存储到ElasticsSearch中。你可以搜索和汇总这些数据，找到任何你感兴趣的信息。
3. 对于程序员来说，比较有名的案例是GitHub，GitHub的搜索是基于ElasticSearch构建的，在github.com/search页面，你可以搜索项目、用户、issue、pull request，还有代码。共有40~50个索引库，分别用于索引网站需要跟踪的各种数据。虽然只索引项目的主分支（master），但这个数据量依然巨大，包括20亿个索引文档，30TB的索引文件。

## elasticsearch基本概念

### Near Realtime(NRT)几乎实时

​		elasticsearch是一个几乎实时的搜索平台。从索引一个文档到这个文档可被搜索只需要很短的时间，一般为毫秒级

### cluster集群

​		集群是一个或者多个节点的集合，**这些节点共同保存整个数据，并在所有节点上提供联合索引和搜索功能。**一个集群由一个唯一集群ID确定，并指定一个集群名（默认为elasticsearch）。该集群名非常重要，因为节点可通过这个集群名来加入集群，一个节点只能是集群的一部分。<br>

​		确保在不同的环境中不要使用相同的集群名称，否则可能导致连接到错误的集群节点。

### Node节点

​		节点是单个服务器实例，它是群集的一部分，**可以存储数据**，**并参与群集的索引和搜索功能**。就像一个集群，节点的名称默认为一个随机的通用唯一标识符（UUID），确定在启动时分配给该节点。如果不希望默认，可以定义任何节点名。这个名字对管理很重要，目的是要确定你的网络服务器对应于你的ElasticSearch群集节点。<br>

​		我们可以通过群集名配置节点以连接特定的群集。默认情况下，每个节点设置加入名为“elasticSearch”的集群。这意味着如果你启动多个节点在网络上，假设他们能发现彼此都会自动形成和加入一个名为“elasticsearch”的集群。<br>

​		在单个群集中，你可以拥有尽可能多的节点。此外，**如果“elasticsearch”在同一个网络中，没有其他节点正在运行，从单个节点的默认情况下会形成一个新的单节点名为”elasticsearch”的集群**。<br>

### Index索引（类似于库）

​		索引是具有相似特性的文档集合。例如，可以为客户数据提供索引，为产品目录建立另一个索引，以及为订单数据建立另一个索引。索引由名称（必须全部为小写）标识，该名称用于在对其中的文档执行索引、搜索、更新和删除操作时引用索引。在单个群集中，你可以定义尽可能多的索引。

### Type类型（类似于表）

​		在索引中，可以定义一个或多个类型。**类型是索引的逻辑类别/分区**，其语义完全取决于你。一般来说，类型定义为具有公共字段集的文档。例如，假设你运行一个博客平台，并将所有数据存储在一个索引中。在这个索引中，你可以为用户数据定义一种类型，为博客数据定义另一种类型，以及为注释数据定义另一类型。

### Document文档（类型于数据行）

​		**文档是可以被索引的信息的基本单位**。例如，你可以为单个客户提供一个文档，单个产品提供另一个文档，以及单个订单提供另一个文档。本文件的表示形式为JSON（JavaScript Object Notation）格式，这是一种非常普遍的互联网数据交换格式。<br>

​		在索引/类型中，你可以存储尽可能多的文档。请注意，尽管文档物理驻留在索引中，文档实际上必须索引或分配到索引中的类型。

### Shards & Replicas分片与副本

​		索引可以存储大量的数据，这些数据可能超过单个节点的硬件限制。例如，十亿个文件占用磁盘空间1TB的单指标可能不适合对单个节点的磁盘或可能太慢服务仅从单个节点的搜索请求。<br>

​		为了解决这一问题，Elasticsearch提供**细分你的指标分成多个块称为分片的能力。当你创建一个索引，你可以简单地定义你想要的分片数量。每个分片本身是一个全功能的、独立的“指数”，可以托管在集群中的任何节点。**

#### Shards分片的重要性主要体现俩个特征

1. 分片允许你水平拆分或缩放内容的大小<br>

2. 分片允许你分配和并行操作的碎片（可能在多个节点上）从而提高性能/吞吐量 这个机制中的碎片是分布式的以及其文件汇总到搜索请求是完全由ElasticSearch管理，对用户来说是透明的。<br>

  
    在同一个集群网络或云环境上，故障是任何时候都会出现的，拥有一个故障转移机制以防分片和节点因为某些原因离线或消失是非常有用的，并且被强烈推荐。为此，Elasticsearch允许你创建一个或多个拷贝，你的索引分片进入所谓的副本或称作复制品的分片，简称Replicas。

#### Replicas副本的重要性主要体现俩个特征

1. 副本为分片或节点失败提供了高可用性。为此，需要注意的是，一个副本的分片不会分配在同一个节点作为原始的或主分片，副本是从主分片那里复制过来的。
2. 副本允许用户扩展你的搜索量或吞吐量，因为搜索可以在所有副本上并行执行。

## ES基本概念与关系型数据库的比较

| ES概念                                         | 关系型数据库       |
| :--------------------------------------------- | :----------------- |
| Index（索引）支持全文检索                      | Database（数据库） |
| Type（类型）                                   | Table（表）        |
| Document（文档），不同文档可以有不同的字段集合 | Row（数据行）      |
| Field（字段）                                  | Column（数据列）   |
| Mapping（映射）                                | Schema（模式）     |

# es restful API

## 查询相关

```shell
# 查看健康状态
curl -X GET 127.0.0.1:9200/_cat/health?v
>epoch      timestamp cluster        status node.total node.data shards pri relo init unassign pending_tasks max_task_wait_time active_shards_percent
1570719865 15:04:25  docker-cluster green           1         1      0   0    0    0        0             0                  -                100.0%

# 查询当前es集群中所有的indices （索引）
curl -X GET 127.0.0.1:9200/_cat/indices?v
>health status index uuid pri rep docs.count docs.deleted store.size pri.store.size
```

## 创建相关

```SHELL
# 创建索引
curl -X PUT 127.0.0.1:9200/www
```

## 删除相关

```shell
# 删除索引
curl -X DELETE 127.0.0.1:9200/www
```

## 插入相关

```shell
# 插入记录
curl -H "Content-Type:application/json" -X POST 127.0.0.1:9200/user/person -d '
{
	"name": "dsb",
	"age": 9000,
	"married": true
}'

#输出：
{
    "_index": "user", # 索引
    "_type": "person", # 类型
    "_id": "MLcwUWwBvEa8j5UrLZj4", # 生成的id
    "_version": 1,
    "result": "created", 
    "_shards": {
        "total": 2, 
        "successful": 1,
        "failed": 0
    },
    "_seq_no": 3, 
    "_primary_term": 1
}

# 也可以使用put，但是需要传入id
curl -H "Content-Type:application/json" -X PUT 127.0.0.1:9200/user/person/4 -d '
{
	"name": "sb",
	"age": 9,
	"married": false
}'
```

## 检索相关

>Elasticsearch的检索语法比较特别，使用GET方法携带JSON格式的查询条件。

```shell
# 全检索
curl -X GET 127.0.0.1:9200/user/person/_search

# 按条件检索：
curl -H "Content-Type:application/json" -X PUT 127.0.0.1:9200/user/person/4 -d '
{
	"query":{
		"match": {"name": "sb"}
	}	
}'

# ElasticSearch默认一次最多返回10条结果，可以像下面的示例通过size字段来设置返回结果的数目。
curl -H "Content-Type:application/json" -X PUT 127.0.0.1:9200/user/person/4 -d '
{
	"query":{
		"match": {"name": "sb"},
		"size": 2
	}	
}'
```



# go操作elasticsearch

## elastic client

使用第三方库https://github.com/olivere/elastic来连接ES并进行操作。

注意下载与你的ES相同版本的client，例如我们这里使用的ES是7.2.1的版本，那么我们下载的client也要与之对应为`github.com/olivere/elastic/v7`。

使用`go.mod`来管理依赖：

```go
require (
    github.com/olivere/elastic/v7 v7.0.4
)
```

## 实例

>更多操作查看：https://godoc.org/github.com/olivere/elastic

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "reflect"

    "gopkg.in/olivere/elastic.v7" //这里使用的是版本5，最新的是6，有改动
)

var client *elastic.Client
var host = "http://127.0.0.1:9200/"

type Employee struct {
    FirstName string   `json:"first_name"`
    LastName  string   `json:"last_name"`
    Age       int      `json:"age"`
    About     string   `json:"about"`
    Interests []string `json:"interests"`
}

//初始化
func init() {
    errorlog := log.New(os.Stdout, "APP", log.LstdFlags)
    var err error
    client, err = elastic.NewClient(elastic.SetErrorLog(errorlog), elastic.SetURL(host))
    if err != nil {
        panic(err)
    }
    info, code, err := client.Ping(host).Do(context.Background())
    if err != nil {
        panic(err)
    }
    fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

    esversion, err := client.ElasticsearchVersion(host)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Elasticsearch version %s\n", esversion)

}

/*下面是简单的CURD*/

//创建
func create() {

    //使用结构体
    e1 := Employee{"Jane", "Smith", 32, "I like to collect rock albums", []string{"music"}}
    put1, err := client.Index().
        Index("megacorp").
        Type("employee").
        Id("1").
        BodyJson(e1).
        Do(context.Background())
    if err != nil {
        panic(err)
    }
    fmt.Printf("Indexed tweet %s to index s%s, type %s\n", put1.Id, put1.Index, put1.Type)

    //使用字符串
    e2 := `{"first_name":"John","last_name":"Smith","age":25,"about":"I love to go rock climbing","interests":["sports","music"]}`
    put2, err := client.Index().
        Index("megacorp").
        Type("employee").
        Id("2").
        BodyJson(e2).
        Do(context.Background())
    if err != nil {
        panic(err)
    }
    fmt.Printf("Indexed tweet %s to index s%s, type %s\n", put2.Id, put2.Index, put2.Type)

    e3 := `{"first_name":"Douglas","last_name":"Fir","age":35,"about":"I like to build cabinets","interests":["forestry"]}`
    put3, err := client.Index().
        Index("megacorp").
        Type("employee").
        Id("3").
        BodyJson(e3).
        Do(context.Background())
    if err != nil {
        panic(err)
    }
    fmt.Printf("Indexed tweet %s to index s%s, type %s\n", put3.Id, put3.Index, put3.Type)

}

//删除
func delete() {

    res, err := client.Delete().Index("megacorp").
        Type("employee").
        Id("1").
        Do(context.Background())
    if err != nil {
        println(err.Error())
        return
    }
    fmt.Printf("delete result %s\n", res.Result)
}

//修改
func update() {
    res, err := client.Update().
        Index("megacorp").
        Type("employee").
        Id("2").
        Doc(map[string]interface{}{"age": 88}).
        Do(context.Background())
    if err != nil {
        println(err.Error())
    }
    fmt.Printf("update age %s\n", res.Result)

}

//查找
func gets() {
    //通过id查找
    get1, err := client.Get().Index("megacorp").Type("employee").Id("2").Do(context.Background())
    if err != nil {
        panic(err)
    }
    if get1.Found {
        fmt.Printf("Got document %s in version %d from index %s, type %s\n", get1.Id, get1.Version, get1.Index, get1.Type)
    }
}

//搜索
func query() {
    var res *elastic.SearchResult
    var err error
    //取所有
    res, err = client.Search("megacorp").Type("employee").Do(context.Background())
    printEmployee(res, err)

    //字段相等
    q := elastic.NewQueryStringQuery("last_name:Smith")
    res, err = client.Search("megacorp").Type("employee").Query(q).Do(context.Background())
    if err != nil {
        println(err.Error())
    }
    printEmployee(res, err)

    //条件查询
    //年龄大于30岁的
    boolQ := elastic.NewBoolQuery()
    boolQ.Must(elastic.NewMatchQuery("last_name", "smith"))
    boolQ.Filter(elastic.NewRangeQuery("age").Gt(30))
    res, err = client.Search("megacorp").Type("employee").Query(q).Do(context.Background())
    printEmployee(res, err)

    //短语搜索 搜索about字段中有 rock climbing
    matchPhraseQuery := elastic.NewMatchPhraseQuery("about", "rock climbing")
    res, err = client.Search("megacorp").Type("employee").Query(matchPhraseQuery).Do(context.Background())
    printEmployee(res, err)

    //分析 interests
    aggs := elastic.NewTermsAggregation().Field("interests")
    res, err = client.Search("megacorp").Type("employee").Aggregation("all_interests", aggs).Do(context.Background())
    printEmployee(res, err)

}

//简单分页
func list(size, page int) {
    if size < 0 || page < 1 {
        fmt.Printf("param error")
        return
    }
    res, err := client.Search("megacorp").
        Type("employee").
        Size(size).
        From((page - 1) * size).
        Do(context.Background())
    printEmployee(res, err)

}

//打印查询到的Employee
func printEmployee(res *elastic.SearchResult, err error) {
    if err != nil {
        print(err.Error())
        return
    }
    var typ Employee
    for _, item := range res.Each(reflect.TypeOf(typ)) { //从搜索结果中取数据的方法
        t := item.(Employee)
        fmt.Printf("%#v\n", t)
    }
}

func main() {
    create()
    delete()
    update()
    gets()
    query()
    list(1, 3)
}
```

## 问题整理

>通过 elastic.v5 驱动连接 ES 有具体报错：“no Elasticsearch node available”
>通过 http 访问的访问仍然能够正常打印出来 ES 的提示信息

```GO
// 连接es抛出异常：panic: no active connection found: no Elasticsearch node available

这里需要关闭Sniff程序：
	Sniff程序自动查找默认群集的所有节点

// 解决：
client, err := elastic.NewClient(elastic.SetSniff(false),elastic.SetURL("http://192.168.1.7:9200"))
```





