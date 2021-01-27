# 利用系统信号来中断服务

> 利用系统信号中断服务，避免业务chan中的缓存的消息，没有完全被消费

```go
// 建立系统信号接收管道
	signals := make(chan os.Signal, 1)
// 监听系统信号
	signal.Notify(signals, os.Interrupt)


	for {
		select {
		case part, ok := <-consumer.Partitions():
			if !ok {
				return
			}
			go func(pc cluster.PartitionConsumer) {
				for msg := range pc.Messages() {
					fmt.Println("kafka堆积情况：", len(KafkaDataClan))
					KafkaDataClan <- msg.Value
				}
			}(part)
      // 如果有信号过来，关闭循环
		case <-signals:
			return
		}
	}
```

