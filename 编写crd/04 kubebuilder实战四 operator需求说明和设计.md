[toc]

# 本篇概览

- 作为《kubebuilder实战》系列的第四篇，经历了前面的充分准备，从本篇开始，咱们来开发一个有实际作用的operator，该operator名为elasticweb，既弹性web服务；
- 这将是一次完整的operator开发实战，设计、编码、部署等环节都会参与到，与《kubebuilder实战之二：初次体验kubebuilder》的不同之处在于，elasticweb从CRD设计再到controller功能都有明确的业务含义，能执行业务逻辑，而《kubebuilder实战之二》仅仅是一次开发流程体验；
- 为了做好这个operator，本篇不急于编码，而是认真的做好设计工作，咱们的operator有什么功能，解决了什么问题，有哪些核心内容，都将在本篇整理清楚，有了这样的准备，才能在下一章写出符合要求的代码；
- 接下来咱们先聊一些背景知识，以便更好的进入正题；

# 需求背景

1. **QPS：**Queries-per-second，既每秒查询率，就是说服务器在一秒的时间内处理了多少个请求；
2. **背景：**做过网站开发的同学对横向扩容应该都了解，简单的说，假设一个tomcat的QPS上限为500，如果外部访问的QPS达到了600，为了保障整个网站服务质量，必须再启动一个同样的tomcat来共同分摊请求，如下图所示(简单起见，假设咱们的后台服务是无状态的，也就是说不依赖宿主机的IP、本地磁盘之类)：

<img src="https://picgo-img.oss-cn-beijing.aliyuncs.com/md-img/2021-12-02/1638415456.png" alt="image-20211202112416830" style="zoom:50%;" />

3. 以上是横向扩容常规做法，在kubernetes环境，如果外部请求超过了单个pod的处理极限，我们可以增加pod数量来达到横向扩容的目的，如下图：

<img src="https://picgo-img.oss-cn-beijing.aliyuncs.com/md-img/2021-12-02/1638420011.png" alt="image-20211202124011761" style="zoom:50%;" />

- 以上就是背景信息，接下来咱们聊聊elasticweb这个operator的具体功能；

# 需求说明

- 为了说清楚需求，这里虚构一个场景：小欣是个java开发者，就是下图这个妹子：

现在小欣要将springboot应用部署到kubernetes上，她的现状和面临的问题如下：

> 1. springboot应用已做成docker镜像；
>
> 2. 通过压测得出单个pod的QPS为500；
>
> 3. 估算得出上线后的总QPS会在800左右；
>
> 4. 随着运营策略变化，QPS还会有调整；
>
> 5. 总的来说，小欣手里只有三个数据：docker镜像、单个pod的QPS、总QPS，她对kubernetes不了解，需要有个方案来帮她将服务部署好，并且在运行期间能支撑外部的高并发访问；


以上就是小欣的需求了，咱们来小结一下：

>    1. 咱们为小欣开发一个operator（名为elasticweb），对小欣来说，她只要将手里的三个参数（docker镜像、单个pod的QPS、总QPS）告诉elasticweb就完事儿了；
>    
>    2. elasticweb在kubernetes创建pod，至于pod数量当然是自动算出来的，要确保能满足QPS要求，以前面的情况为例，需要两个pod才能满足800的QPS；
>    
>    3. 单个pod的QPS和总QPS都随时可能变化，一旦有变，elasticweb也要自动调整pod数量，以确保服务质量；
>    
>    4. 为了确保服务可以被外部调用，咱们再顺便帮小欣创建好service（她对kubernetes了解不多，这事儿咱们就顺手做了吧）；

# 自保声明

- 看过上述需求后，聪明的您一定会对我投来鄙视的眼光，其实kubernetes早就有现成的QPS调节方案了，例如修改deployment的副本数、单个pod纵向扩容、autoscale等都可以，本次使用operator来实现仅仅是为了展示operator的开发过程，并不是说自定义operator是唯一的解决方案；

- 所以，如果您觉得我这种用operator实现扩容的方式很low，请不要把我骂得太惨，我这也只是为了展示operator开发过程而已，况且咱这个operator也不是一无是处，用了这个operator，您就不用关注pod数量了，只要聚焦单实例QPS和总QPS即可，这两个参数更贴近业务；

- 为了不把事情弄复杂，**假设每个pod所需的CPU和内存是固定的**，直接在operator代码中写死，其实您也可以自己改代码，改成可以在外部配置，就像镜像名称参数那样；

- 把需求都交代清楚了，接下来进入设计环节，先把CRD设计出来，这可是核心的数据结构；

# CRD设计之Spec部分

Spec是用来保存用户的期望值的，也就是小欣手里的三个参数（docker镜像、单个pod的QPS、总QPS），再加上端口号：

```shell
image：业务服务对应的镜像
port：service占用的宿主机端口，外部请求通过此端口访问pod的服务
singlePodQPS：单个pod的QPS上限
totalQPS：当前整个业务的总QPS
```

对小欣来说，输入这四个参数就完事儿了；

# CRD设计之Status部分

- Status用来保存实际值，这里设计成只有一个字段realQPS，表示当前整个operator实际能支持的QPS，这样无论何时，只要小欣用kubectl describe命令就能知道当前系统实际上能支持多少QPS；

# CRD源码

- 把数据结构说明白的最好方法就是看代码：

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
// omitempty 允许为nil
type ElasticWebStatus struct {
	// 当前kubernetes中实际支持的总QPS
	RealQPS *int32 `json:"realQPS"`
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

# 业务逻辑设计

> CRD的完成代表核心数据结构已经确定，接下来是业务逻辑的设计，主要是理清楚controller的Reconcile方法里面做些啥，其实核心逻辑还是非常简单的：算出需要多少个pod，然后通过更新deployment让pod数量达到要求，在此核心的基础上再把创建deployment和service、更新status这些琐碎的事情做好，就完事儿了；

这里将整个业务逻辑的流程图给出来如下所示，用于指导开发：

<img src="https://picgo-img.oss-cn-beijing.aliyuncs.com/md-img/2021-12-02/1638427770.png" alt="image-20211202144930240" style="zoom:50%;" />

- 至此，咱们完成了整个elasticweb的需求和设计，聪明的您肯定已经胸有成竹，而且迫不及待的想启动开发了，好的，下一篇咱们正式开始编码！

