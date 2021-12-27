[toc]

# 本篇概览

> 作为《kubebuilder实战》系列的第六篇，前面已完成了编码，现在到了验证功能的环节，请确保您的docker和kubernetes环境正常，然后咱们一起完成以下操作：

1. 部署CRD
2. 本地运行Controller
3. 通过yaml文件新建elasticweb资源对象
4. 通过日志和kubectl命令验证elasticweb功能是否正常
5. 浏览器访问web，验证业务服务是否正常
6. 修改singlePodQPS，看elasticweb是否自动调整pod数量
7. 修改totalQPS，看elasticweb是否自动调整pod数
8. 删除elasticweb，看相关的service和deployment被自动删除
9. 构建Controller镜像，在kubernetes运行此Controller，验证上述功能是否正常

看似简单的部署验证操作，零零散散加起来居然有这么多…好吧不感慨了，立即开始吧；

# 部署CRD

```shell
# 从控制台进入Makefile所在目录，执行命令make install，即可将CRD部署到kubernetes：
zhaoqin@zhaoqindeMBP-2 elasticweb % make install

# 可以用命令kubectl api-versions验证CRD部署是否成功：
zhaoqin@zhaoqindeMBP-2 elasticweb % kubectl api-versions|grep elasticweb
```

# 本地运行Controller

先尝试用最简单的方式来验证Controller的功能，如下图，Macbook电脑是我的开发环境，直接用elasticweb工程中的Makefile，可以将Controller的代码在本地运行起来里面：

<img src="https://picgo-img.oss-cn-beijing.aliyuncs.com/md-img/2021-12-02/1638429680.png" alt="image-20211202152120165" style="zoom:50%;" />

```shell 
# 进入Makefile文件所在目录，执行命令make run即可编译运行controller
zhaoqin@zhaoqindeMBP-2 elasticweb % make run
```

# 新建elasticweb资源对象

- 负责处理elasticweb的Controller已经运行起来了，接下来就开始创建elasticweb资源对象吧，用yaml文件来创建；

- 在`config/samples`目录下，kubebuilder为咱们创建了demo文件`elasticweb_v1_elasticweb.yaml`，不过这里面spec的内容不是咱们定义的那四个字段，需要改成以下内容：
```yaml
  apiVersion: v1
  kind: Namespace
  metadata:
    name: dev
    labels:
      name: dev
  ---
  apiVersion: elasticweb.com.bolingcavalry/v1
  kind: ElasticWeb
  metadata:
    namespace: dev
    name: elasticweb-sample
  spec:
    # Add fields here
    image: tomcat:8.0.18-jre8
    port: 30003
    singlePodQPS: 500
    totalQPS: 600
```

- 对上述配置的几个参数做如下说明：
  1. 使用的namespace为dev
  2. 本次测试部署的应用为tomcat
  3. service使用宿主机的30003端口暴露tomcat的服务
  4. 假设单个pod能支撑500QPS，外部请求的QPS为600
- 执行命令`kubectl apply -f config/samples/elasticweb_v1_elasticweb.yaml`，即可在kubernetes创建elasticweb实例：

```shell
zhaoqin@zhaoqindeMBP-2 elasticweb % kubectl apply -f config/samples/elasticweb_v1_elasticweb.yaml
```

- 去controller的窗口发现打印了不少日志，通过分析日志发现Reconcile方法执行了两次，第一执行时创建了deployment和service等资源：

```shell
2021-02-21T10:03:57.108+0800    INFO    controllers.ElasticWeb  1. start reconcile logic        {"elasticweb": "dev/elasticweb-sample"}
2021-02-21T10:03:57.108+0800    INFO    controllers.ElasticWeb  3. instance : Image [tomcat:8.0.18-jre8], Port [30003], SinglePodQPS [500], TotalQPS [600], RealQPS [nil]       {"elasticweb": "dev/elasticweb-sample"}
2021-02-21T10:03:57.210+0800    INFO    controllers.ElasticWeb  4. deployment not exists        {"elasticweb": "dev/elasticweb-sample"}
2021-02-21T10:03:57.313+0800    INFO    controllers.ElasticWeb  set reference   {"func": "createService"}
2021-02-21T10:03:57.313+0800    INFO    controllers.ElasticWeb  start create service    {"func": "createService"}
2021-02-21T10:03:57.364+0800    INFO    controllers.ElasticWeb  create service success  {"func": "createService"}
2021-02-21T10:03:57.365+0800    INFO    controllers.ElasticWeb  expectReplicas [2]      {"func": "createDeployment"}
2021-02-21T10:03:57.365+0800    INFO    controllers.ElasticWeb  set reference   {"func": "createDeployment"}
2021-02-21T10:03:57.365+0800    INFO    controllers.ElasticWeb  start create deployment {"func": "createDeployment"}
2021-02-21T10:03:57.382+0800    INFO    controllers.ElasticWeb  create deployment success       {"func": "createDeployment"}
2021-02-21T10:03:57.382+0800    INFO    controllers.ElasticWeb  singlePodQPS [500], replicas [2], realQPS[1000] {"func": "updateStatus"}
2021-02-21T10:03:57.407+0800    DEBUG   controller-runtime.controller   Successfully Reconciled {"controller": "elasticweb", "request": "dev/elasticweb-sample"}
2021-02-21T10:03:57.407+0800    INFO    controllers.ElasticWeb  1. start reconcile logic        {"elasticweb": "dev/elasticweb-sample"}
2021-02-21T10:03:57.407+0800    INFO    controllers.ElasticWeb  3. instance : Image [tomcat:8.0.18-jre8], Port [30003], SinglePodQPS [500], TotalQPS [600], RealQPS [1000]      {"elasticweb": "dev/elasticweb-sample"}
2021-02-21T10:03:57.407+0800    INFO    controllers.ElasticWeb  9. expectReplicas [2], realReplicas [2] {"elasticweb": "dev/elasticweb-sample"}
2021-02-21T10:03:57.407+0800    INFO    controllers.ElasticWeb  10. return now  {"elasticweb": "dev/elasticweb-sample"}
2021-02-21T10:03:57.407+0800    DEBUG   controller-runtime.controller   Successfully Reconciled {"controller": "elasticweb", "request": "dev/elasticweb-sample"}
```

- 再用kubectl get命令详细检查资源对象，一切符合预期，elasticweb、service、deployment、pod都是正常的：

```
zhaoqin@zhaoqindeMBP-2 elasticweb % kubectl apply -f config/samples/elasticweb_v1_elasticweb.yaml
namespace/dev created
elasticweb.elasticweb.com.bolingcavalry/elasticweb-sample created
zhaoqin@zhaoqindeMBP-2 elasticweb % kubectl get elasticweb -n dev                                 
NAME                AGE
elasticweb-sample   35s
zhaoqin@zhaoqindeMBP-2 elasticweb % kubectl get service -n dev                                    
NAME                TYPE       CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
elasticweb-sample   NodePort   10.107.177.158   <none>        8080:30003/TCP   41s
zhaoqin@zhaoqindeMBP-2 elasticweb % kubectl get deployment -n dev                                 
NAME                READY   UP-TO-DATE   AVAILABLE   AGE
elasticweb-sample   2/2     2            2           46s
zhaoqin@zhaoqindeMBP-2 elasticweb % kubectl get pod -n dev                                        
NAME                                 READY   STATUS    RESTARTS   AGE
elasticweb-sample-56fc5848b7-l5thk   1/1     Running   0          50s
elasticweb-sample-56fc5848b7-lqjk5   1/1     Running   0          50s
```

# 验证业务功能

 ```
 curl 10.107.177.158:8080
 ```

# 修改单个Pod的QPS

- 如果自身优化，或者外界依赖变化(如缓存、数据库扩容)，这些都可能导致当前服务的QPS提升，假设单个Pod的QPS从500提升到了800，看看咱们的Operator能不能自动做出调整（总QPS是600，因此pod数应该从2降到1）

- 在`config/samples/`目录下新增名为`update_single_pod_qps.yaml`的文件，内容如下：

```yaml
spec:
singlePodQPS: 800
```

- 执行以下命令，即可将单个Pod的QPS从500更新为800（注意，参数type很重要别漏了）：

```
kubectl patch elasticweb elasticweb-sample \
-n dev \
--type merge \
--patch "$(cat config/samples/update_single_pod_qps.yaml)"
```
-  此时去看controller日志，如下图，红框1表示spec已经更新，红框2则表示用最新的参数计算出来的pod数量，符合预期：

![image-20211202152736672](https://picgo-img.oss-cn-beijing.aliyuncs.com/md-img/2021-12-02/1638430056.png)

- 用kubectl get命令检查pod，可见已经降到1个了：

```shell
zhaoqin@zhaoqindeMBP-2 elasticweb % kubectl get pod -n dev                                                                                       
NAME                                 READY   STATUS    RESTARTS   AGE
elasticweb-sample-56fc5848b7-l5thk   1/1     Running   0          30m
```

- 记得用浏览器检查tomcat是否正常；

# 修改总QPS

- 外部QPS也在频繁变化中，咱们的operator也需要根据总QPS及时调节pod实例，以确保整体服务质量，接下来咱们就修改总QPS，看operator是否生效：
- 在`config/samples/`目录下新增名为`update_total_qps.yaml`的文件，内容如下：

```yaml
spec:
  totalQPS: 2600
```

- 执行以下命令，即可将总QPS从600更新为2600（注意，参数type很重要别漏了）：

```shell
kubectl patch elasticweb elasticweb-sample \
-n dev \
--type merge \
--patch "$(cat config/samples/update_total_qps.yaml)"
```

此时去看controller日志，如下图，红框1表示spec已经更新，红框2则表示用最新的参数计算出来的pod数量，符合预期：

![image-20211202154405871](https://picgo-img.oss-cn-beijing.aliyuncs.com/md-img/2021-12-02/1638431046.png)

- 用kubectl get命令检查pod，可见已经增长到4个，4个pd的能支撑的QPS为3200，满足了当前2600的要求：

```shell
zhaoqin@zhaoqindeMBP-2 elasticweb % kubectl get pod -n dev
NAME                                 READY   STATUS    RESTARTS   AGE
elasticweb-sample-56fc5848b7-8n7tq   1/1     Running   0          8m22s
elasticweb-sample-56fc5848b7-f2lpb   1/1     Running   0          8m22s
elasticweb-sample-56fc5848b7-l5thk   1/1     Running   0          48m
elasticweb-sample-56fc5848b7-q8p5f   1/1     Running   0          8m22s
```

- 记得用浏览器检查tomcat是否正常；

# 删除验证

- 目前整个dev这个namespace下有service、deployment、pod、elasticweb这些资源对象，如果要全部删除，只需删除elasticweb即可，因为service和deployment都和elasticweb建立的关联关系，代码如下图红框：

<img src="https://picgo-img.oss-cn-beijing.aliyuncs.com/md-img/2021-12-02/1638431105.png" alt="image-20211202154505668" style="zoom:50%;" />

- 执行删除elasticweb的命令：

```shell
kubectl delete elasticweb elasticweb-sample -n dev
```

- 再去查看其他资源，都被自动删除了：

```shell
zhaoqin@zhaoqindeMBP-2 elasticweb % kubectl delete elasticweb elasticweb-sample -n dev
elasticweb.elasticweb.com.bolingcavalry "elasticweb-sample" deleted
zhaoqin@zhaoqindeMBP-2 elasticweb % kubectl get pod -n dev                            
NAME                                 READY   STATUS        RESTARTS   AGE
elasticweb-sample-56fc5848b7-9lcww   1/1     Terminating   0          45s
elasticweb-sample-56fc5848b7-n7p7f   1/1     Terminating   0          45s
zhaoqin@zhaoqindeMBP-2 elasticweb % kubectl get pod -n dev
NAME                                 READY   STATUS        RESTARTS   AGE
elasticweb-sample-56fc5848b7-n7p7f   0/1     Terminating   0          73s
zhaoqin@zhaoqindeMBP-2 elasticweb % kubectl get pod -n dev
No resources found in dev namespace.
zhaoqin@zhaoqindeMBP-2 elasticweb % kubectl get deployment -n dev
No resources found in dev namespace.
zhaoqin@zhaoqindeMBP-2 elasticweb % kubectl get service -n dev   
No resources found in dev namespace.
zhaoqin@zhaoqindeMBP-2 elasticweb % kubectl get namespace dev 
NAME   STATUS   AGE
dev    Active   97s
```

# 构建镜像并发布

```shell
# 构建镜像
make docker-build docker-push IMG=bolingcavalry/elasticweb:002

# 发布到k8s
make deploy IMG=bolingcavalry/elasticweb:002

# 接下来像之前那样创建elasticweb资源对象，验证所有资源是否创建成功：
kubectl apply -f config/samples/elasticweb_v1_elasticweb.yaml 


# 查看部署的服务的日志
kubectl logs -f \
elasticweb-controller-manager-5795d4d98d-t6jvc \
-c manager \
-n elasticweb-system
```

# 卸载和清理

- 体验完毕后，如果想把前面创建的资源全部清理掉(注意，是清理资源，不是资源对象)，可以执行以下命令：

```shell
make uninstall
```

 至此，整个operator的设计、开发、部署、验证流程就全部完成了，在您的operator开发过程中，希望本文能给您带来一些参考；

