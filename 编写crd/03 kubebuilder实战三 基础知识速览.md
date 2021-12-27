[toc]

# 本篇概览

- 作为《kubebuilder实战》系列的第三篇，本该进入真枪实弹的operator开发环节，却突然发现kubebuilder涉及的知识点太多太零散，如果现在就敲命令写代码去实战，即便完成了一次operator开发，但缺失大量信息（例如操作顺序怎么安排、步骤之间如何关联等），不但《kubebuilder实战》系列失去参考价值，过几个月就连我自己都看不懂这些内容了，因此，本篇暂缓实战，咱们一起对kubebuilder开发过程中的知识点做一次速记，再从容的启动开发工作；
- 特别说明：webhook是operator中的重要功能，其理论和实战都需要大量篇幅，因此后面会有这方面专门的文章，本文不会涉及webhook的知识点；
- 接下来，大串讲开始；
  

# 知识储备

能看懂kubebuilder官方demo的代码、会用client对象操作kubernetes资源，以上两点是胜任operator开发的最基本要求，否则在开发过程中会有种寸步难行的感觉，达到这些条件需要少量的知识储备，现在欣宸已经为您准备好了，希望您能简单浏览一下：

- 《Kubernetes的Group、Version、Resource学习小记》
- 《client-go实战之一：准备工作》
- 《client-go实战之二:   RESTClient》
- 《client-go实战之三：Clientset》
- 《client-go实战之四：dynamicClient》
- 《client-go实战之五：DiscoveryClient》

# 初始化相关知识点

- 创建operator的第一步，就是用kubebuilder命令行创建整个项目，这个在[《kubebuilder实战之二：初次体验kubebuilder》](https://xinchen.blog.csdn.net/article/details/113089414)已经试过，当时执行的是如下三行命令：

```shell
mkdir -p $GOPATH/src/helloworld
cd $GOPATH/src/helloworld
kubebuilder init --domain com.bolingcavalry
```

- 在用上module之后，大家已经脱离了\$GOPATH的束缚，像上面那样中规中矩的去$GOPATH/src下面操作就略有些别扭了，来试试不用\$GOPATH的初始化方式；
  1. 随处新建一个目录(路径中不要有中文和空格)，例如`/Users/zhaoqin/temp/202102/15/elasticweb`
  2. 在目录中用`go mod init elasticweb`命令新建名为elasticweb的工程；
  3. 再执行`kubebuilder init --domain com.bolingcavalry`，即可新建operator工程；

# 基础设施

- operator工程新建完成后，会新增不少文件和目录，以下几个是官方提到的基础设施：
  1. go.mod：module的配置文件，里面已经填充了几个重要依赖；
  2. Makefile：非常重要的工具，前文咱们也用过了，编译构建、部署、运行都会用到；
  3. PROJECT：kubebuilder工程的元数据，在生成各种API的时候会用到这里面的信息；
  4. config/default：基于kustomize制作的配置文件，为controller提供标准配置，也可以按需要去修改调整；
  5. config/manager：一些和manager有关的细节配置，例如镜像的资源限制；
  6. config/rbac：顾名思义，如果像限制operator在kubernetes中的操作权限，就要通过rbac来做精细的权限配置了，这里面就是权限配置的细节；

## main.go

- main.go是kubebuilder自动生成的代码，这是operator的启动代码，里面有几处值得注意：

1. 两个全局变量，如下所示，setupLog用于输出日志无需多说，scheme也是常用的工具，它提供了Kind和Go代码中的数据结构的映射，：

```golang
var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)
```

2. 另外还有些设置，例如监控指标相关的，以及管理controller和webhook的manager，它会一直运行下去直到被外部终止，关于这个manage还有一处要注意的地方，就是它的参数，下图是默认的参数，如果您想让operator在指定namespace范围内生效，还可以在下午的地方新增Namespace参数，如果要指定多个nanespace，就使用cache.MultiNamespacedCacheBuilder(namespaces)参数：
   ![image-20211202110500437](https://picgo-img.oss-cn-beijing.aliyuncs.com/md-img/2021-12-02/1638414300.png)

- main.go的内容在大多数场景无需改动，了解即可，接下来的API是重头戏；

## API相关（数据核心）

- API是operator的核心，当您决定使用operator时，就应该从真实需求出发，开始设计整个CRD，而这些设计最终体现在CRD的数据结构，以及对真实值和期望值的处理逻辑中；
- 在《kubebuilder实战之二：初次体验kubebuilder》咱们创建过API，当时的命令是：

```shell
kubebuilder create api \
--group webapp \
--version v1 \
--kind Guestbook
```

kubebuilder自动新增了很多内容，如下图，都是为了这个CRD服务的：

<img src="https://picgo-img.oss-cn-beijing.aliyuncs.com/md-img/2021-12-02/1638414647.png" alt="image-20211202111047638" style="zoom:50%;" />

- 新增的内容中，最核心的当然是CRD了，也就是上图中Guestbook数据结构所在的guestbook_types.go文件，这个最重要的数据结构如下：

```go
type Guestbook struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GuestbookSpec   `json:"spec,omitempty"`
	Status GuestbookStatus `json:"status,omitempty"`
}
```

1. metav1.TypeMeta：保存了资源的Group、Version、Kind
2. metav1.ObjectMeta：保存了资源对象的名称和namespace
3. Spec：期望状态，例如deployment在创建时指定了pod有三个副本
4. Status：真实状态，例如deployment在创建后只有一个副本(其他的还没有创建成功)，大多数资源对象都有此字段，不过ConfigMap是个例外（想想也是，配置信息嘛，配成啥就是啥，没有什么期望值和真实值的说法）；
5. 还有一个数据结构，就是Guestbook对应的列表GuestbookList，就是单个资源对象的集合；
   guestbook_types.go所在目录下还有两个文件：groupversion_info.go定义了Group和Version，以及注册到scheme时用到的实例SchemeBuilder，zz_generated.deepcopy.go用于实现实例的深拷贝，它们都无需修改，了解即可；

## controller相关（业务核心）

- 前面聊过了数据核心，接下来要讨论如何实现业务需求了，在operator开发过程中，尽管业务逻辑各不相同，但有两个共性：
  1. Status(真实状态)是个数据结构，其字段是业务定义的，其字段值也是业务代码执行自定义的逻辑算出来的；
  2. 业务核心的目标，是确保Status与Spec达成一致，例如deployment指定了pod的副本数为3，如果真实的pod没有三个，deployment的controller代码就去创建pod，如果真实的pod超过了三个，deployment的controller代码就去删除pod;
- 以上就是咱们的controller要做的事情，接下来看看代码的细节，kubebuilder创建的guestbook_controller.go就是controller，咱们的业务代码都写在这个文件中，来看看kubebuilder帮我们准备了什么：
  3. 数据结构定义，如下所示，操作资源对象时用到的客户端工具client.Client、日志工具、Kind和数据结构的关系Scheme，这些都帮我们准备好了，真贴心：

```go
type GuestbookReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}
```
  4. SetupWithManager方法，在main.go中有调用，指定了Guestbook这个资源的变化会被manager监控，从而触发Reconcile方法：

```go
func (r *GuestbookReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webappv1.Guestbook{}).
		Complete(r)
}
```

5. 如下图，Reconcile方法前面有一些**+kubebuilder:rbac**前缀的注释，这些是用来确保controller在运行时有对应的资源操作权限，例如红框中就是我自己添加的(**注意这里需要保持与func的行之间的空行**)，这样controller就有权查询pod资源对象了：

<img src="https://picgo-img.oss-cn-beijing.aliyuncs.com/md-img/2021-12-02/1638414931.png" alt="image-20211202111531811" style="zoom:50%;" />

- guestbook_controller.go是operator的业务核心，而controller的核心是其Reconcile方法，将来咱们的大部分代码都是写在这里面的，主要做的事情就是获取status，然后让status和spec达成一致；

- 关于status，官方的一段描述值得重视，如下图红框，主要是说资源对象的状态应该是每次重新计算出来的，这里以deployment为例，想知道当前有多少个pod，有两种方法，第一种准备个字段记录，每次对pod的增加和删除都修改这个字段，于是读这个字段就知道pod数量了，第二种方法是每次用client工具去实时查询pod数量，目前官方明确推荐使用第二种方法：

---

至此，基础知识串讲就完成了，咱们按照官方资料的顺序把知识点过了一遍，接下来，就是按照官方资料的顺序去实战了，让大家久等了，下一篇，operator实战；

