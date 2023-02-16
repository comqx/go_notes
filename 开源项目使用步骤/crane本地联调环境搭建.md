


# 本地启动

```shell
 go run cmd/craned/main.go   --recommendation-configuration-file=./testing/config.yaml  --prometheus-address=http://10.0.169.10:30670 --kubeconfig=./testing/kube-config  --webhook-enabled=false --feature-gates=Analysis=true,TimeSeriesPrediction=true,Autoscaling=true --leader-elect=false
```

# ide启动

<img src="https://picgo-img.oss-cn-beijing.aliyuncs.com/md-img/2023-02-15/1676449344.png" alt="image-20230215162224467" style="zoom:50%;" />

