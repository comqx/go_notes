[toc]

# 本篇概览

本文是《kubebuilder实战》系列的第七篇，之前的文章咱们完成了一个Operator的设计、开发、部署、验证过程，为了让整个过程保持简洁并且篇幅不膨胀，实战中刻意跳过了一个重要的知识点：webhook，如今是时候学习它了，这是个很重要的功能；
本篇由以下部分构成：

1. 介绍webhook；
2. 结合前面的elasticweb项目，设计一个使用webhook的场景；
3. 准备工作
4. 生成webhook
5. 开发(配置)
6. 开发(编码)
7. 部署
8. 验证Defaulter(添加默认值)
9. 验证Validator(合法性校验)

# 关于webhook

 熟悉java开发的读者大多知道过滤器(Servlet Filter)，如下图，外部请求会先到达过滤器，做一些统一的操作，例如转码、校验，然后才由真正的业务逻辑处理请求：

<img src="https://picgo-img.oss-cn-beijing.aliyuncs.com/md-img/2021-12-03/1638512013.png" alt="image-20211203141333066" style="zoom:50%;" />

Operator中的webhook，其作用与上述过滤器类似，外部对CRD资源的变更，在Controller处理之前都会交给webhook提前处理，流程如下图，该图来自《Getting Started with Kubernetes | Operator and Operator Framework》：

<img src="https://picgo-img.oss-cn-beijing.aliyuncs.com/md-img/2021-12-03/1638512035.png" alt="image-20211203141354964" style="zoom:50%;" />

再来看看webhook具体做了哪些事情，如下图，kubernetes[官方博客](https://kubernetes.io/blog/2019/03/21/a-guide-to-kubernetes-admission-controllers/)明确指出webhook可以做两件事：修改(mutating)和验证(validating)

![image-20211203141421311](https://picgo-img.oss-cn-beijing.aliyuncs.com/md-img/2021-12-03/1638512061.png)

- kubebuilder为我们提供了生成webhook的基础文件和代码的工具，与制作API的工具类似，极大地简化了工作量，咱们只需聚焦业务实现即可；
- 基于kubebuilder制作的webhook和controller，如果是同一个资源，那么它们在同一个进程中；

# 设计实战场景

- 为了让实战有意义，咱们为前面的elasticweb项目上增加需求，让webhook发挥实际作用；

1. 如果用户忘记输入总QPS，系统webhook负责设置默认值**1300**，操作如下图：

<img src="https://picgo-img.oss-cn-beijing.aliyuncs.com/md-img/2021-12-03/1638512089.png" alt="image-20211203141449439" style="zoom:50%;" />

2. 为了保护系统，给单个pod的QPS设置上限**1000**，如果外部输入的singlePodQPS值超过**1000**，就创建资源对象失败，如下图所示：

<img src="https://picgo-img.oss-cn-beijing.aliyuncs.com/md-img/2021-12-03/1638512114.png" alt="image-20211203141514505" style="zoom:50%;" />

# 准备工作

- 和controller类似，webhook既能在kubernetes环境中运行，也能在kubernetes环境之外运行；
- 如果webhook在kubernetes环境之外运行，是有些麻烦的，需要将证书放在所在环境，默认地址是：

```shell
/tmp/k8s-webhook-server/serving-certs/tls.{crt,key}
```

- 为了省事儿，也为了更接近生产环境的用法，接下来的实战的做法是将webhook部署在kubernetes环境中
- 为了让webhook在kubernetes环境中运行，咱们要做一点准备工作安装cert manager，执行以下操作：

```shell 
kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.2.0/cert-manager.yaml
```

- 上述操作完成后会新建很多资源，如namespace、rbac、pod等，以pod为例如下：

```
[root@hedy ~]# kubectl get pods --all-namespaces
NAMESPACE        NAME                                       READY   STATUS    RESTARTS   AGE
cert-manager     cert-manager-6588898cb4-nvnz8              1/1     Running   1          5d14h
cert-manager     cert-manager-cainjector-7bcbdbd99f-q645r   1/1     Running   1          5d14h
cert-manager     cert-manager-webhook-5fd9f9dd86-98tm9      1/1     Running   1          5d14h
...
```

- 操作完成后，准备工作结束，可以开始实战了；

# 生成webhook

- 进入elasticweb工程下，执行以下命令创建webhook：

```shell
kubebuilder create webhook \
--group elasticweb \
--version v1 \
--kind ElasticWeb \
--defaulting \
--programmatic-validation
```

- 上述命令执行完毕后，先去看看main.go文件，如下图红框1所示，自动增加了一段代码，作用是让webhook生效：

<img src="https://picgo-img.oss-cn-beijing.aliyuncs.com/md-img/2021-12-03/1638512215.png" alt="image-20211203141655661" style="zoom:50%;" />

- 上图红框2中的`elasticweb_webhook.go`就是新增文件，内容如下：

```golang
package v1

import (
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var elasticweblog = logf.Log.WithName("elasticweb-resource")

func (r *ElasticWeb) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-elasticweb-com-bolingcavalry-v1-elasticweb,mutating=true,failurePolicy=fail,groups=elasticweb.com.bolingcavalry,resources=elasticwebs,verbs=create;update,versions=v1,name=melasticweb.kb.io

var _ webhook.Defaulter = &ElasticWeb{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *ElasticWeb) Default() {
	elasticweblog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// +kubebuilder:webhook:verbs=create;update,path=/validate-elasticweb-com-bolingcavalry-v1-elasticweb,mutating=false,failurePolicy=fail,groups=elasticweb.com.bolingcavalry,resources=elasticwebs,versions=v1,name=velasticweb.kb.io

var _ webhook.Validator = &ElasticWeb{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *ElasticWeb) ValidateCreate() error {
	elasticweblog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *ElasticWeb) ValidateUpdate(old runtime.Object) error {
	elasticweblog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *ElasticWeb) ValidateDelete() error {
	elasticweblog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
```

- 上述代码有两处需要注意，第一处和填写默认值有关，如下图：

<img src="https://picgo-img.oss-cn-beijing.aliyuncs.com/md-img/2021-12-03/1638512258.png" alt="image-20211203141738658" style="zoom:50%;" />

- 第二处和校验有关，如下图：

<img src="https://picgo-img.oss-cn-beijing.aliyuncs.com/md-img/2021-12-03/1638512282.png" alt="image-20211203141802553" style="zoom:50%;" />

- 咱们要实现的业务需求就是通过修改上述elasticweb_webhook.go的内容来实现，不过代码稍后再写，先把配置都改好；

# 开发(配置)

- 打开文件config/default/kustomization.yaml，下图四个红框中的内容原本都被注释了，现在请将注释符号都删掉，使其生效：

<img src="https://picgo-img.oss-cn-beijing.aliyuncs.com/md-img/2021-12-03/1638512324.png" alt="image-20211203141844776" style="zoom:50%;" />

- 还是文件config/default/kustomization.yaml，节点vars下面的内容，原本全部被注释了，现在请全部放开，放开后的效果如下图：

<img src="https://picgo-img.oss-cn-beijing.aliyuncs.com/md-img/2021-12-03/1638512345.png" alt="image-20211203141905396" style="zoom:50%;" />

- 配置已经完成，可以编码了；

# 开发(编码)

- 打开文件elasticweb_webhook.go
- 新增依赖：

```go
apierrors "k8s.io/apimachinery/pkg/api/errors"
```

- 找到Default方法，改成如下内容，可见代码很简单，判断TotalQPS是否存在，若不存在就写入默认值，另外还加了两行日志：

```go
func (r *ElasticWeb) Default() {
	elasticweblog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
	// 如果创建的时候没有输入总QPS，就设置个默认值
	if r.Spec.TotalQPS == nil {
		r.Spec.TotalQPS = new(int32)
		*r.Spec.TotalQPS = 1300
		elasticweblog.Info("a. TotalQPS is nil, set default value now", "TotalQPS", *r.Spec.TotalQPS)
	} else {
		elasticweblog.Info("b. TotalQPS exists", "TotalQPS", *r.Spec.TotalQPS)
	}
}
```

- 接下来开发校验功能，咱们把校验功能封装成一个`validateElasticWeb`方法，然后在新增和修改的时候各调用一次，如下，可见最终是调用`apierrors.NewInvalid`生成错误实例的，而此方法接受的是多个错误，因此要为其准备切片做入参，当然了，如果是多个参数校验失败，可以都放入切片中：

```go
func (r *ElasticWeb) validateElasticWeb() error {
	var allErrs field.ErrorList

	if *r.Spec.SinglePodQPS > 1000 {
		elasticweblog.Info("c. Invalid SinglePodQPS")

		err := field.Invalid(field.NewPath("spec").Child("singlePodQPS"),
			*r.Spec.SinglePodQPS,
			"d. must be less than 1000")

		allErrs = append(allErrs, err)

		return apierrors.NewInvalid(
			schema.GroupKind{Group: "elasticweb.com.bolingcavalry", Kind: "ElasticWeb"},
			r.Name,
			allErrs)
	} else {
		elasticweblog.Info("e. SinglePodQPS is valid")
		return nil
	}
}
```

- 再找到新增和修改资源对象时被调用的方法，在里面调用validateElasticWeb：

```go
// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *ElasticWeb) ValidateCreate() error {
	elasticweblog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.

	return r.validateElasticWeb()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *ElasticWeb) ValidateUpdate(old runtime.Object) error {
	elasticweblog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return r.validateElasticWeb()
}

```

- 编码完成，可见非常简单，接下来，咱们把以前实战遗留的东西清理一下，再开始新的部署和验证；

# 清理工作

 ```shell
 # 删除elasticweb资源对象：
 kubectl delete -f config/samples/elasticweb_v1_elasticweb.yaml
 
 # 删除controller
 kustomize build config/default | kubectl delete -f -
 
 # 删除CRD
 make uninstall
 ```

# 部署

```shell
# 部署CRD
make install

# 构建镜像并推送到仓库
make docker-build docker-push IMG=registry.cn-hangzhou.aliyuncs.com/bolingcavalry/elasticweb:001

# 部署集成了webhook功能的controller：
make deploy IMG=registry.cn-hangzhou.aliyuncs.com/bolingcavalry/elasticweb:001

# 查看pod，确认启动成功：
zhaoqin@zhaoqindeMBP-2 ~ % kubectl get pods --all-namespaces
NAMESPACE           NAME                                             READY   STATUS    RESTARTS   AGE
cert-manager        cert-manager-6588898cb4-nvnz8                    1/1     Running   1          5d21h
cert-manager        cert-manager-cainjector-7bcbdbd99f-q645r         1/1     Running   1          5d21h
cert-manager        cert-manager-webhook-5fd9f9dd86-98tm9            1/1     Running   1          5d21h
elasticweb-system   elasticweb-controller-manager-7dcbfd4675-898gb   2/2     Running   0          20s
```

#  验证Defaulter(添加默认值)

 ```yaml 
 # 修改文件config/samples/elasticweb_v1_elasticweb.yaml，修改后的内容如下，可见totalQPS字段已经被注释掉了：
 
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
   # totalQPS: 600
 ```

```shell
 # 创建一个elasticweb资源对象：
 kubectl apply -f config/samples/elasticweb_v1_elasticweb.yaml
 
 # 1. 此时单个pod的QPS是500，如果webhook的代码生效的话，总QPS就是1300，而对应的pod数应该是3个，接下来咱们看看是否符合预期；
 # 2. 先看elasticweb、deployment、pod等资源对象是否正常，如下所示，全部符合预期：
 zhaoqin@zhaoqindeMBP-2 ~ % kubectl get elasticweb -n dev                                                                 
 NAME                AGE
 elasticweb-sample   89s
 zhaoqin@zhaoqindeMBP-2 ~ % kubectl get deployments -n dev
 NAME                READY   UP-TO-DATE   AVAILABLE   AGE
 elasticweb-sample   3/3     3            3           98s
 zhaoqin@zhaoqindeMBP-2 ~ % kubectl get service -n dev    
 NAME                TYPE       CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
 elasticweb-sample   NodePort   10.105.125.125   <none>        8080:30003/TCP   106s
 zhaoqin@zhaoqindeMBP-2 ~ % kubectl get pod -n dev    
 NAME                                 READY   STATUS    RESTARTS   AGE
 elasticweb-sample-56fc5848b7-5tkxw   1/1     Running   0          113s
 elasticweb-sample-56fc5848b7-blkzg   1/1     Running   0          113s
 elasticweb-sample-56fc5848b7-pd7jg   1/1     Running   0          113s
 
 
 # 用kubectl describe命令查看elasticweb资源对象的详情，如下所示，TotalQPS字段被webhook设置为1300，RealQPS也计算正确：
zhaoqin@zhaoqindeMBP-2 ~ % kubectl describe elasticweb elasticweb-sample -n dev
Name:         elasticweb-sample
Namespace:    dev
Labels:       <none>
Annotations:  <none>
API Version:  elasticweb.com.bolingcavalry/v1
Kind:         ElasticWeb
Metadata:
  Creation Timestamp:  2021-02-27T16:07:34Z
  Generation:          2
  Managed Fields:
    API Version:  elasticweb.com.bolingcavalry/v1
    Fields Type:  FieldsV1
    fieldsV1:
      f:metadata:
        f:annotations:
          .:
          f:kubectl.kubernetes.io/last-applied-configuration:
      f:spec:
        .:
        f:image:
        f:port:
        f:singlePodQPS:
    Manager:      kubectl-client-side-apply
    Operation:    Update
    Time:         2021-02-27T16:07:34Z
    API Version:  elasticweb.com.bolingcavalry/v1
    Fields Type:  FieldsV1
    fieldsV1:
      f:status:
        f:realQPS:
    Manager:         manager
    Operation:       Update
    Time:            2021-02-27T16:07:34Z
  Resource Version:  687628
  UID:               703de111-d859-4cd2-b3c4-1d201fb7bd7d
Spec:
  Image:           tomcat:8.0.18-jre8
  Port:            30003
  Single Pod QPS:  500
  Total QPS:       1300
Status:
  Real QPS:  1500
Events:      <none>
```

- 再来看看controller的日志，其中的webhook部分是否符合预期，如下图红框所示，发现TotalQPS字段为空，就将设置为默认值，并且在检测的时候SinglePodQPS的值也没有超过1000：

<img src="https://picgo-img.oss-cn-beijing.aliyuncs.com/md-img/2021-12-03/1638512660.png" alt="image-20211203142420069" style="zoom:50%;" />

- 最后别忘了用浏览器验证web服务是否正常，我这里的完整地址是：http://192.168.50.75:30003/
- 至此，咱们完成了webhook的Defaulter验证，接下来验证Validator

# 验证Validator

- 接下来该验证webhook的参数校验功能了，先验证修改时的逻辑；
- 编辑文件config/samples/update_single_pod_qps.yaml，值如下：

```yaml
spec:
  singlePodQPS: 1100
```

- 用patch命令使之生效：

```shell
kubectl patch elasticweb elasticweb-sample \
-n dev \
--type merge \
--patch "$(cat config/samples/update_single_pod_qps.yaml)"
```

- 此时，控制台会输出错误信息：

```shell
Error from server (ElasticWeb.elasticweb.com.bolingcavalry "elasticweb-sample" is invalid: spec.singlePodQPS: Invalid value: 1100: d. must be less than 1000): admission webhook "velasticweb.kb.io" denied the request: ElasticWeb.elasticweb.com.bolingcavalry "elasticweb-sample" is invalid: spec.singlePodQPS: Invalid value: 1100: d. must be less than 1000
```

- 接下来再试试webhook在新增时候的校验功能；

- 清理前面创建的elastic资源对象，执行命令：

```shell
kubectl delete -f config/samples/elasticweb_v1_elasticweb.yaml
```

  修改文件，如下图红框所示，咱们将singlePodQPS的值改为超过1000，看看webhook是否能检查到这个错误，并阻止资源对象的创建：

<img src="https://picgo-img.oss-cn-beijing.aliyuncs.com/md-img/2021-12-03/1638512784.png" alt="image-20211203142624098" style="zoom:50%;" />

- 执行以下命令开始创建elasticweb资源对象：

```
kubectl apply -f config/samples/elasticweb_v1_elasticweb.yaml
```

- 控制台提示以下信息，包含了咱们代码中写入的错误描述，证明elasticweb资源对象创建失败，证明webhook的Validator功能已经生效：

```shell
namespace/dev created
Error from server (ElasticWeb.elasticweb.com.bolingcavalry "elasticweb-sample" is invalid: spec.singlePodQPS: Invalid value: 1500: d. must be less than 1000): error when creating "config/samples/elasticweb_v1_elasticweb.yaml": admission webhook "velasticweb.kb.io" denied the request: ElasticWeb.elasticweb.com.bolingcavalry "elasticweb-sample" is invalid: spec.singlePodQPS: Invalid value: 1500: d. must be less than 1000
```

- 不放心的话执行kubectl get命令检查一下，发现空空如也：

```shell
zhaoqin@zhaoqindeMBP-2 elasticweb % kubectl get elasticweb -n dev       
No resources found in dev namespace.
zhaoqin@zhaoqindeMBP-2 elasticweb % kubectl get deployments -n dev
No resources found in dev namespace.
zhaoqin@zhaoqindeMBP-2 elasticweb % kubectl get service -n dev
No resources found in dev namespace.
zhaoqin@zhaoqindeMBP-2 elasticweb % kubectl get pod -n dev
No resources found in dev namespace.
```

还要看下controller日志，如下图红框所示，符合预期：

<img src="https://picgo-img.oss-cn-beijing.aliyuncs.com/md-img/2021-12-03/1638512836.png" alt="image-20211203142715949" style="zoom:50%;" />

- 至此，operator的webhook的开发、部署、验证咱们就完成了，整个elasticweb也算是基本功能齐全，希望能为您的operator开发提供参考；
