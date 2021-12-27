[toc]

# 本篇概览

- 本篇是《kubebuilder实战》系列的第五篇，前面的一切努力(环境准备、知识储备、需求分析、数据结构和业务逻辑设计)，都是为了将之前的设计用编码实现；
- 既然已经充分准备，如今无需太多言语，咱们开始动手吧！

# 新建项目elasticweb
新建名为elasticweb的文件夹，在里面执行以下命令即可创建名为elasticweb的项目，domain为com.bolingcavalry：
```shell
go mod init elasticweb
kubebuilder init --domain com.bolingcavalry
```
然后是CRD，执行以下命令即可创建相关资源：

```shell
kubebuilder create api \
--group elasticweb \
--version v1 \
--kind ElasticWeb
```
然后用IDE打开整个工程，我这里是goland：
<img src="https://picgo-img.oss-cn-beijing.aliyuncs.com/md-img/2021-12-02/1638428145.png" alt="image-20211202145545954" style="zoom:50%;" />

# CRD编码

- 打开文件`api/v1/elasticweb_types.go`，做以下几步改动：
  1. 修改数据结构ElasticWebSpec，增加前文设计的四个字段；
  2. 修改数据结构ElasticWebStatus，增加前文设计的一个字段；
  3. 增加String方法，这样打印日志时方便我们查看，注意RealQPS字段是指针，因此可能为空，需要判空；
- 完整的`elasticweb_types.go`如下所示：
  
```go
package v1

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
)

// 期望状态
type ElasticWebSpec struct {
	// 业务服务对应的镜像，包括名称:tag
	Image string `json:"image"`
	// service占用的宿主机端口，外部请求通过此端口访问pod的服务
	Port *int32 `json:"port"`

	// 单个pod的QPS上限
	SinglePodQPS *int32 `json:"singlePodQPS"`
	// 当前整个业务的总QPS
	TotalQPS *int32 `json:"totalQPS"`
}

// 实际状态，该数据结构中的值都是业务代码计算出来的
type ElasticWebStatus struct {
	// 当前kubernetes中实际支持的总QPS
  // omitempty 需要允许该参数为nil，不然在webhook测试的时候回出现报错
	RealQPS *int32 `json:"realQPS,omitempty"`
}

// +kubebuilder:object:root=true

// ElasticWeb is the Schema for the elasticwebs API
type ElasticWeb struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ElasticWebSpec   `json:"spec,omitempty"`
	Status ElasticWebStatus `json:"status,omitempty"`
}

func (in *ElasticWeb) String() string {
	var realQPS string

	if nil == in.Status.RealQPS {
		realQPS = "nil"
	} else {
		realQPS = strconv.Itoa(int(*(in.Status.RealQPS)))
	}

	return fmt.Sprintf("Image [%s], Port [%d], SinglePodQPS [%d], TotalQPS [%d], RealQPS [%s]",
		in.Spec.Image,
		*(in.Spec.Port),
		*(in.Spec.SinglePodQPS),
		*(in.Spec.TotalQPS),
		realQPS)
}

// +kubebuilder:object:root=true

// ElasticWebList contains a list of ElasticWeb
type ElasticWebList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ElasticWeb `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ElasticWeb{}, &ElasticWebList{})
}
```

- 在elasticweb目录下执行`make install`即可部署CRD到kubernetes：

```shell
zhaoqin@zhaoqindeMBP-2 elasticweb % make install
/Users/zhaoqin/go/bin/controller-gen "crd:trivialVersions=true" rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases
kustomize build config/crd | kubectl apply -f -
Warning: apiextensions.k8s.io/v1beta1 CustomResourceDefinition is deprecated in v1.16+, unavailable in v1.22+; use apiextensions.k8s.io/v1 CustomResourceDefinition
customresourcedefinition.apiextensions.k8s.io/elasticwebs.elasticweb.com.bolingcavalry created
```

- 部署成功后，用api-versions命令可以查到该GV：

![image-20211202150850331](https://picgo-img.oss-cn-beijing.aliyuncs.com/md-img/2021-12-02/1638428930.png)

# 回顾业务逻辑

- 核心数据结构设计编码完毕，接下来该编写业务逻辑代码了，大家还记得前文设计的业务流程吧，简单回顾一下，如下图：

<img src="https://picgo-img.oss-cn-beijing.aliyuncs.com/md-img/2021-12-02/1638428648.png" alt="image-20211202150407854" style="zoom:50%;" />

- 打开文件elasticweb_controller.go，接下来咱们逐渐添加内容；

# 添加资源访问权限

- 咱们的elasticweb会对service、deployment这两种资源做查询、新增、修改等操作，因此需要这些资源的操作权限，增加下图红框中的两行注释，这样代码生成工具就会在RBAC配置中增加对应的权限：

![image-20211202151035532](https://picgo-img.oss-cn-beijing.aliyuncs.com/md-img/2021-12-02/1638429035.png)

# 常量定义

- 先把常量准备好，可见每个pod使用的CPU和内存都是在此固定的，您也可以改成在Spec中定义，这样就可以从外部传入了，另外这里为每个pod只分配了0.1个CPU，主要是因为我穷买不起好的CPU，您可以酌情调整该值：

```go
const (
	// deployment中的APP标签名
	APP_NAME = "elastic-app"
	// tomcat容器的端口号
	CONTAINER_PORT = 8080
	// 单个POD的CPU资源申请
	CPU_REQUEST = "100m"
	// 单个POD的CPU资源上限
	CPU_LIMIT = "100m"
	// 单个POD的内存资源申请
	MEM_REQUEST = "512Mi"
	// 单个POD的内存资源上限
	MEM_LIMIT = "512Mi"
)
```

# 方法getExpectReplicas

- 有个很重要的逻辑：根据单个pod的QPS和总QPS，计算需要多少个pod，咱们将这个逻辑封装到一个方法中以便使用：

```go
// 根据单个QPS和总QPS计算pod数量
func getExpectReplicas(elasticWeb *elasticwebv1.ElasticWeb) int32 {
	// 单个pod的QPS
	singlePodQPS := *(elasticWeb.Spec.SinglePodQPS)

	// 期望的总QPS
	totalQPS := *(elasticWeb.Spec.TotalQPS)

	// Replicas就是要创建的副本数
	replicas := totalQPS / singlePodQPS

	if totalQPS%singlePodQPS > 0 {
		replicas++
	}

	return replicas
}
```

# 方法createServiceIfNotExists

- 将创建service的操作封装到一个方法中，是的主干代码的逻辑更清晰，可读性更强；

- 创建service的时候，有几处要注意：

  1. 先查看service是否存在，不存在才创建；
  2. 将service和CRD实例elasticWeb建立关联(controllerutil.SetControllerReference方法)，这样当elasticWeb被删除的时候，service会被自动删除而无需我们干预；
  3. 创建service的时候用到了client-go工具，推荐您阅读《client-go实战系列》,工具越熟练，编码越尽兴；

- 创建service的完整方法如下：

```golang
// 新建service
func createServiceIfNotExists(ctx context.Context, r *ElasticWebReconciler, elasticWeb *elasticwebv1.ElasticWeb, req ctrl.Request) error {
	log := r.Log.WithValues("func", "createService")

	service := &corev1.Service{}

	err := r.Get(ctx, req.NamespacedName, service)

	// 如果查询结果没有错误，证明service正常，就不做任何操作
	if err == nil {
		log.Info("service exists")
		return nil
	}

	// 如果错误不是NotFound，就返回错误
	if !errors.IsNotFound(err) {
		log.Error(err, "query service error")
		return err
	}

	// 实例化一个数据结构
	service = &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: elasticWeb.Namespace,
			Name:      elasticWeb.Name,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{{
				Name:     "http",
				Port:     8080,
				NodePort: *elasticWeb.Spec.Port,
			},
			},
			Selector: map[string]string{
				"app": APP_NAME,
			},
			Type: corev1.ServiceTypeNodePort,
		},
	}

	// 这一步非常关键！
	// 建立关联后，删除elasticweb资源时就会将deployment也删除掉
	log.Info("set reference")
	if err := controllerutil.SetControllerReference(elasticWeb, service, r.Scheme); err != nil {
		log.Error(err, "SetControllerReference error")
		return err
	}

	// 创建service
	log.Info("start create service")
	if err := r.Create(ctx, service); err != nil {
		log.Error(err, "create service error")
		return err
	}

	log.Info("create service success")

	return nil
}
```

# 方法createDeployment

- 将创建deployment的操作封装在一个方法中，同样是为了将主干逻辑保持简洁；
- 创建deployment的方法也有几处要注意：
1. 调用getExpectReplicas方法得到要创建的pod的数量，该数量是创建deployment时的一个重要参数；
2. 每个pod所需的CPU和内存资源也是deployment的参数；
3. 将deployment和elasticweb建立关联，这样删除elasticweb的时候deplyment就会被自动删除了；
4. 同样是使用client-go客户端工具创建deployment资源；
```go
// 新建deployment
func createDeployment(ctx context.Context, r *ElasticWebReconciler, elasticWeb *elasticwebv1.ElasticWeb) error {
	log := r.Log.WithValues("func", "createDeployment")

	// 计算期望的pod数量
	expectReplicas := getExpectReplicas(elasticWeb)

	log.Info(fmt.Sprintf("expectReplicas [%d]", expectReplicas))

	// 实例化一个数据结构
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: elasticWeb.Namespace,
			Name:      elasticWeb.Name,
		},
		Spec: appsv1.DeploymentSpec{
			// 副本数是计算出来的
			Replicas: pointer.Int32Ptr(expectReplicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": APP_NAME,
				},
			},

			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": APP_NAME,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name: APP_NAME,
							// 用指定的镜像
							Image:           elasticWeb.Spec.Image,
							ImagePullPolicy: "IfNotPresent",
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									Protocol:      corev1.ProtocolSCTP,
									ContainerPort: CONTAINER_PORT,
								},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									"cpu":    resource.MustParse(CPU_REQUEST),
									"memory": resource.MustParse(MEM_REQUEST),
								},
								Limits: corev1.ResourceList{
									"cpu":    resource.MustParse(CPU_LIMIT),
									"memory": resource.MustParse(MEM_LIMIT),
								},
							},
						},
					},
				},
			},
		},
	}

	// 这一步非常关键！
	// 建立关联后，删除elasticweb资源时就会将deployment也删除掉
	log.Info("set reference")
	if err := controllerutil.SetControllerReference(elasticWeb, deployment, r.Scheme); err != nil {
		log.Error(err, "SetControllerReference error")
		return err
	}

	// 创建deployment
	log.Info("start create deployment")
	if err := r.Create(ctx, deployment); err != nil {
		log.Error(err, "create deployment error")
		return err
	}

	log.Info("create deployment success")

	return nil
}
```

# 方法updateStatus

- 不论是创建deployment资源对象，还是对已有的deployment的pod数量做调整，这些操作完成后都要去修改Status，既实际的状态，这样外部才能随时随地知道当前elasticweb支持多大的QPS，因此需要将修改Status的操作封装到一个方法中，给多个场景使用，Status的计算逻辑很简单：pod数量乘以每个pod的QPS就是总QPS了，代码如下：

```go
// 完成了pod的处理后，更新最新状态
func updateStatus(ctx context.Context, r *ElasticWebReconciler, elasticWeb *elasticwebv1.ElasticWeb) error {
	log := r.Log.WithValues("func", "updateStatus")

	// 单个pod的QPS
	singlePodQPS := *(elasticWeb.Spec.SinglePodQPS)

	// pod总数
	replicas := getExpectReplicas(elasticWeb)

	// 当pod创建完毕后，当前系统实际的QPS：单个pod的QPS * pod总数
	// 如果该字段还没有初始化，就先做初始化
	if nil == elasticWeb.Status.RealQPS {
		elasticWeb.Status.RealQPS = new(int32)
	}

	*(elasticWeb.Status.RealQPS) = singlePodQPS * replicas

	log.Info(fmt.Sprintf("singlePodQPS [%d], replicas [%d], realQPS[%d]", singlePodQPS, replicas, *(elasticWeb.Status.RealQPS)))

	if err := r.Update(ctx, elasticWeb); err != nil {
		log.Error(err, "update instance error")
		return err
	}

	return nil
}
```

# 主干代码

- 前面细枝末节都处理完毕，可以开始主流程了，有了前面的流程图的赋值，主流程的代码很容就写出来了，如下所示，已经添加了足够的注释，就不再赘述了：

```go
func (r *ElasticWebReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	// 会用到context
	ctx := context.Background()
	log := r.Log.WithValues("elasticweb", req.NamespacedName)

	// your logic here

	log.Info("1. start reconcile logic")

	// 实例化数据结构
	instance := &elasticwebv1.ElasticWeb{}

	// 通过客户端工具查询，查询条件是
	err := r.Get(ctx, req.NamespacedName, instance)

	if err != nil {

		// 如果没有实例，就返回空结果，这样外部就不再立即调用Reconcile方法了
		if errors.IsNotFound(err) {
			log.Info("2.1. instance not found, maybe removed")
			return reconcile.Result{}, nil
		}

		log.Error(err, "2.2 error")
		// 返回错误信息给外部
		return ctrl.Result{}, err
	}

	log.Info("3. instance : " + instance.String())

	// 查找deployment
	deployment := &appsv1.Deployment{}

	// 用客户端工具查询
	err = r.Get(ctx, req.NamespacedName, deployment)

	// 查找时发生异常，以及查出来没有结果的处理逻辑
	if err != nil {
		// 如果没有实例就要创建了
		if errors.IsNotFound(err) {
			log.Info("4. deployment not exists")

			// 如果对QPS没有需求，此时又没有deployment，就啥事都不做了
			if *(instance.Spec.TotalQPS) < 1 {
				log.Info("5.1 not need deployment")
				// 返回
				return ctrl.Result{}, nil
			}

			// 先要创建service
			if err = createServiceIfNotExists(ctx, r, instance, req); err != nil {
				log.Error(err, "5.2 error")
				// 返回错误信息给外部
				return ctrl.Result{}, err
			}

			// 立即创建deployment
			if err = createDeployment(ctx, r, instance); err != nil {
				log.Error(err, "5.3 error")
				// 返回错误信息给外部
				return ctrl.Result{}, err
			}

			// 如果创建成功就更新状态
			if err = updateStatus(ctx, r, instance); err != nil {
				log.Error(err, "5.4. error")
				// 返回错误信息给外部
				return ctrl.Result{}, err
			}

			// 创建成功就可以返回了
			return ctrl.Result{}, nil
		} else {
			log.Error(err, "7. error")
			// 返回错误信息给外部
			return ctrl.Result{}, err
		}
	}

	// 如果查到了deployment，并且没有返回错误，就走下面的逻辑

	// 根据单QPS和总QPS计算期望的副本数
	expectReplicas := getExpectReplicas(instance)

	// 当前deployment的期望副本数
	realReplicas := *deployment.Spec.Replicas

	log.Info(fmt.Sprintf("9. expectReplicas [%d], realReplicas [%d]", expectReplicas, realReplicas))

	// 如果相等，就直接返回了
	if expectReplicas == realReplicas {
		log.Info("10. return now")
		return ctrl.Result{}, nil
	}

	// 如果不等，就要调整
	*(deployment.Spec.Replicas) = expectReplicas

	log.Info("11. update deployment's Replicas")
	// 通过客户端更新deployment
	if err = r.Update(ctx, deployment); err != nil {
		log.Error(err, "12. update deployment replicas error")
		// 返回错误信息给外部
		return ctrl.Result{}, err
	}

	log.Info("13. update status")

	// 如果更新deployment的Replicas成功，就更新状态
	if err = updateStatus(ctx, r, instance); err != nil {
		log.Error(err, "14. update status error")
		// 返回错误信息给外部
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}
```

- 至此，整个elasticweb operator编码就完成了，限于篇幅，咱们把部署、运行、镜像制作等操作放在下一篇文章吧；

