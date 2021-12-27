[toc]

# 本篇概览

本文是《kubebuilder实战》系列的第二篇，前文将kubebuilder环境准备完毕，今天咱们在此环境创建CRD和Controller，再部署到kubernetes环境并且验证是否生效，整篇文章由以下内容组成：

- 创建API(CRD和Controller)
- 构建和部署CRD
- 编译和运行controller
- 创建CRD对应的实例
- 删除实例并停止controller
- 将controller制作成docker镜像
- 卸载和清理

# 创建helloworld项目

1. 执行以下命令，创建helloworld项目：

```shell
mkdir -p $GOPATH/src/helloworld
cd $GOPATH/src/helloworld
kubebuilder init --domain com.bolingcavalry
```

2. 控制台输出类似以下内容：

```shell
[root@kubebuilder helloworld]# kubebuilder init --domain com.bolingcavalry
Writing scaffold for you to edit...
Get controller runtime:
$ go get sigs.k8s.io/controller-runtime@v0.5.0
Update go.mod:
$ go mod tidy
Running make:
$ make
/root/gopath/bin/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
go fmt ./...
go vet ./...
go build -o bin/manager main.go
Next: define a resource with:
$ kubebuilder create api
```

3. 等待数分钟后创建完成，在$GOPATH/src/helloworld目录下新增以下内容，可见这是个标准的go module工程：
```shell
[root@kubebuilder ~]# tree $GOPATH/src/helloworld
/root/gopath/src/helloworld
   ├── bin
   │   └── manager
   ├── config
   │   ├── certmanager
   │   │   ├── certificate.yaml
   │   │   ├── kustomization.yaml
   │   │   └── kustomizeconfig.yaml
   │   ├── default
   │   │   ├── kustomization.yaml
   │   │   ├── manager_auth_proxy_patch.yaml
   │   │   ├── manager_webhook_patch.yaml
   │   │   └── webhookcainjection_patch.yaml
   │   ├── manager
   │   │   ├── kustomization.yaml
   │   │   └── manager.yaml
   │   ├── prometheus
   │   │   ├── kustomization.yaml
   │   │   └── monitor.yaml
   │   ├── rbac
   │   │   ├── auth_proxy_client_clusterrole.yaml
   │   │   ├── auth_proxy_role_binding.yaml
   │   │   ├── auth_proxy_role.yaml
   │   │   ├── auth_proxy_service.yaml
   │   │   ├── kustomization.yaml
   │   │   ├── leader_election_role_binding.yaml
   │   │   ├── leader_election_role.yaml
   │   │   └── role_binding.yaml
   │   └── webhook
   │       ├── kustomization.yaml
   │       ├── kustomizeconfig.yaml
   │       └── service.yaml
   ├── Dockerfile
   ├── go.mod
   ├── go.sum
   ├── hack
   │   └── boilerplate.go.txt
   ├── main.go
   ├── Makefile
   └── PROJECT

9 directories, 30 files
```

4. 创建API(CRD和Controller)

```shell
# 1. 接下来要要创建资源相关的内容了，group/version/kind这三部分可以确定资源的唯一身份，命令如下：
cd $GOPATH/src/helloworld
kubebuilder create api \
--group webapp \
--version v1 \
--kind Guestbook
# 2. 控制台会提醒是否创建资源(Create Resource [y/n])，输入y
# 3. 接下来控制台会提醒是否创建控制器(Create Controller [y/n])，输入y
# 4. kubebuilder会根据上述命令新增多个文件，如下图红框所示：
```

<img src="https://picgo-img.oss-cn-beijing.aliyuncs.com/md-img/2021-12-01/1638352465.png" alt="image-20211201175425439" style="zoom:50%;" />

## 构建和部署CRD

1. kubebuilder提供的Makefile将构建和部署工作大幅度简化，执行以下命令会将最新构建的CRD部署在kubernetes上：

```shell
cd $GOPATH/src/helloworld
make install
```

1. 控制台输出如下内容，提示部署成功：

```shell
[root@kubebuilder helloworld]# make install
/root/gopath/bin/controller-gen "crd:trivialVersions=true" rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases
kustomize build config/crd | kubectl apply -f -
Warning: apiextensions.k8s.io/v1beta1 CustomResourceDefinition is deprecated in v1.16+, unavailable in v1.22+; use apiextensions.k8s.io/v1 CustomResourceDefinition
customresourcedefinition.apiextensions.k8s.io/guestbooks.webapp.com.bolingcavalry created
```

## 编译和运行controller

1. kubebuilder自动生成的controller源码地址是：$GOPATH/src/helloworld/controllers/guestbook_controller.go ， 内容如下：

```golang
package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	webappv1 "helloworld/api/v1"
)

// GuestbookReconciler reconciles a Guestbook object
type GuestbookReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=webapp.com.bolingcavalry,resources=guestbooks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=webapp.com.bolingcavalry,resources=guestbooks/status,verbs=get;update;patch

func (r *GuestbookReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("guestbook", req.NamespacedName)

	// your logic here

	return ctrl.Result{}, nil
}

func (r *GuestbookReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webappv1.Guestbook{}).
		Complete(r)
}
```

2. 本文以体验基本流程为主，不深入研究源码，所以对上面的代码仅做少量修改，用于验证是否能生效，改动如下图红框所示：

<img src="https://picgo-img.oss-cn-beijing.aliyuncs.com/md-img/2021-12-01/1638352652.png" alt="image-20211201175732819" style="zoom:50%;" />

3. 执行以下命令，会编译并启动刚才修改的controller：

```shell
cd $GOPATH/src/helloworld
make run
```

4. 此时控制台输出以下内容，这里要注意，controller是在kubebuilder电脑上运行的，一旦使用Ctrl+c中断控制台，就会导致controller停止：

```shell
[root@kubebuilder helloworld]# make run
/root/gopath/bin/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
go fmt ./...
go vet ./...
/root/gopath/bin/controller-gen "crd:trivialVersions=true" rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases
go run ./main.go
2021-01-23T20:58:35.107+0800	INFO	controller-runtime.metrics	metrics server is starting to listen	{"addr": ":8080"}
2021-01-23T20:58:35.108+0800	INFO	setup	starting manager
2021-01-23T20:58:35.108+0800	INFO	controller-runtime.manager	starting metrics server	{"path": "/metrics"}
2021-01-23T20:58:35.108+0800	INFO	controller-runtime.controller	Starting EventSource	{"controller": "guestbook", "source": "kind source: /, Kind="}
2021-01-23T20:58:35.208+0800	INFO	controller-runtime.controller	Starting Controller	{"controller": "guestbook"}
2021-01-23T20:58:35.209+0800	INFO	controller-runtime.controller	Starting workers	{"controller": "guestbook", "worker count": 1}
```

## 创建Guestbook资源的实例

1. 现在kubernetes已经部署了Guestbook类型的CRD，而且对应的controller也已正在运行中，可以尝试创建Guestbook类型的实例了(相当于有了pod的定义后，才可以创建pod)；
2. kubebuilder已经自动创建了一个类型的部署文件：$GOPATH/src/helloworld/config/samples/webapp_v1_guestbook.yaml ，内容如下，很简单，接下来咱们就用这个文件来创建Guestbook实例：

```yaml
apiVersion: webapp.com.bolingcavalry/v1
kind: Guestbook
metadata:
  name: guestbook-sample
spec:
  # Add fields here
  foo: bar
```

3. 重新打开一个控制台，登录kubebuilder电脑，执行以下命令即可创建Guestbook类型的实例：

```shell
cd $GOPATH/src/helloworld
kubectl apply -f config/samples/
```

4. 如下所示，控制台提示资源创建成功：

```shell
[root@kubebuilder helloworld]# kubectl apply -f config/samples/
guestbook.webapp.com.bolingcavalry/guestbook-sample created
```

5. 用kubectl get命令可以看到实例已经创建：

```shell
[root@kubebuilder helloworld]# kubectl get Guestbook
NAME               AGE
guestbook-sample   112s
```

6. 用命令kubectl edit Guestbook guestbook-sample编辑该实例，修改的内容如下图红框所示：

<img src="https://picgo-img.oss-cn-beijing.aliyuncs.com/md-img/2021-12-01/1638354747.png" alt="image-20211201183227490" style="zoom:50%;" />

7. 此时去controller所在控制台，可以看到新增和修改的操作都有日志输出，咱们新增的日志都在里面，代码调用栈一目了然：

```shell
2021-01-24T09:51:50.418+0800	INFO	controllers.Guestbook	1. default/guestbook-sample
2021-01-24T09:51:50.418+0800	INFO	controllers.Guestbook	2. goroutine 188 [running]:
runtime/debug.Stack(0xc0002a1808, 0xc0002fc600, 0x1b)
	/root/go/src/runtime/debug/stack.go:24 +0x9f
helloworld/controllers.(*GuestbookReconciler).Reconcile(0xc0003c9dd0, 0xc0002d02f9, 0x7, 0xc0002d02e0, 0x10, 0x12f449647b, 0xc000456f30, 0xc000456ea8, 0xc000456ea0)
	/root/gopath/src/helloworld/controllers/guestbook_controller.go:49 +0x1a9
sigs.k8s.io/controller-runtime/pkg/internal/controller.(*Controller).reconcileHandler(0xc00022a480, 0x1430e00, 0xc0003e7560, 0x0)
	/root/gopath/pkg/mod/sigs.k8s.io/controller-runtime@v0.5.0/pkg/internal/controller/controller.go:256 +0x166
sigs.k8s.io/controller-runtime/pkg/internal/controller.(*Controller).processNextWorkItem(0xc00022a480, 0xc000469600)
	/root/gopath/pkg/mod/sigs.k8s.io/controller-runtime@v0.5.0/pkg/internal/controller/controller.go:232 +0xb0
sigs.k8s.io/controller-runtime/pkg/internal/controller.(*Controller).worker(0xc00022a480)
	/root/gopath/pkg/mod/sigs.k8s.io/controller-runtime@v0.5.0/pkg/internal/controller/controller.go:211 +0x2b
k8s.io/apimachinery/pkg/util/wait.JitterUntil.func1(0xc000292980)
	/root/gopath/pkg/mod/k8s.io/apimachinery@v0.17.2/pkg/util/wait/wait.go:152 +0x5f
k8s.io/apimachinery/pkg/util/wait.JitterUntil(0xc000292980, 0x3b9aca00, 0x0, 0x1609101, 0xc000102480)
	/root/gopath/pkg/mod/k8s.io/apimachinery@v0.17.2/pkg/util/wait/wait.go:153 +0x105
k8s.io/apimachinery/pkg/util/wait.Until(0xc000292980, 0x3b9aca00, 0xc000102480)
	/root/gopath/pkg/mod/k8s.io/apimachinery@v0.17.2/pkg/util/wait/wait.go:88 +0x4d
created by sigs.k8s.io/controller-runtime/pkg/internal/controller.(*Controller).Start.func1
	/root/gopath/pkg/mod/sigs.k8s.io/controller-runtime@v0.5.0/pkg/internal/controller/controller.go:193 +0x32d

2021-01-24T09:51:50.418+0800	DEBUG	controller-runtime.controller	Successfully Reconciled	{"controller": "guestbook", "request": "default/guestbook-sample"}


2021-01-24T09:52:33.632+0800	INFO	controllers.Guestbook	1. default/guestbook-sample
2021-01-24T09:52:33.633+0800	INFO	controllers.Guestbook	2. goroutine 188 [running]:
runtime/debug.Stack(0xc0002a1808, 0xc0003fa5e0, 0x1b)
	/root/go/src/runtime/debug/stack.go:24 +0x9f
helloworld/controllers.(*GuestbookReconciler).Reconcile(0xc0003c9dd0, 0xc0002d02f9, 0x7, 0xc0002d02e0, 0x10, 0x1d0410fe42, 0xc000456f30, 0xc000456ea8, 0xc000456ea0)
	/root/gopath/src/helloworld/controllers/guestbook_controller.go:49 +0x1a9
sigs.k8s.io/controller-runtime/pkg/internal/controller.(*Controller).reconcileHandler(0xc00022a480, 0x1430e00, 0xc0003d24c0, 0x0)
	/root/gopath/pkg/mod/sigs.k8s.io/controller-runtime@v0.5.0/pkg/internal/controller/controller.go:256 +0x166
sigs.k8s.io/controller-runtime/pkg/internal/controller.(*Controller).processNextWorkItem(0xc00022a480, 0xc000469600)
	/root/gopath/pkg/mod/sigs.k8s.io/controller-runtime@v0.5.0/pkg/internal/controller/controller.go:232 +0xb0
sigs.k8s.io/controller-runtime/pkg/internal/controller.(*Controller).worker(0xc00022a480)
	/root/gopath/pkg/mod/sigs.k8s.io/controller-runtime@v0.5.0/pkg/internal/controller/controller.go:211 +0x2b
k8s.io/apimachinery/pkg/util/wait.JitterUntil.func1(0xc000292980)
	/root/gopath/pkg/mod/k8s.io/apimachinery@v0.17.2/pkg/util/wait/wait.go:152 +0x5f
k8s.io/apimachinery/pkg/util/wait.JitterUntil(0xc000292980, 0x3b9aca00, 0x0, 0x1609101, 0xc000102480)
	/root/gopath/pkg/mod/k8s.io/apimachinery@v0.17.2/pkg/util/wait/wait.go:153 +0x105
k8s.io/apimachinery/pkg/util/wait.Until(0xc000292980, 0x3b9aca00, 0xc000102480)
	/root/gopath/pkg/mod/k8s.io/apimachinery@v0.17.2/pkg/util/wait/wait.go:88 +0x4d
created by sigs.k8s.io/controller-runtime/pkg/internal/controller.(*Controller).Start.func1
	/root/gopath/pkg/mod/sigs.k8s.io/controller-runtime@v0.5.0/pkg/internal/controller/controller.go:193 +0x32d

2021-01-24T09:52:33.633+0800	DEBUG	controller-runtime.controller	Successfully Reconciled	{"controller": "guestbook", "request": "default/guestbook-sample"}
```

## 删除实例并停止controller

1. 不再需要Guestbook实例的时候，执行以下命令即可删除：

```
cd $GOPATH/src/helloworld
kubectl delete -f config/samples/
```

2. 不再需要controller的时候，去它的控制台使用Ctrl+c中断即可；

## 将controller制作成docker镜像,并部署到k8s

```shell
# 制作镜像
cd $GOPATH/src/helloworld
make docker-build docker-push IMG=bolingcavalry/guestbook:002

# 修改镜像，部署到k8s
[root@kubebuilder ~]# cd $GOPATH/src/helloworld
[root@kubebuilder helloworld]# make deploy IMG=bolingcavalry/guestbook:002
/root/gopath/bin/controller-gen "crd:trivialVersions=true" rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases
cd config/manager && kustomize edit set image controller=bolingcavalry/guestbook:002
kustomize build config/default | kubectl apply -f -
namespace/helloworld-system created
Warning: apiextensions.k8s.io/v1beta1 CustomResourceDefinition is deprecated in v1.16+, unavailable in v1.22+; use apiextensions.k8s.io/v1 CustomResourceDefinition
customresourcedefinition.apiextensions.k8s.io/guestbooks.webapp.com.bolingcavalry configured
role.rbac.authorization.k8s.io/helloworld-leader-election-role created
clusterrole.rbac.authorization.k8s.io/helloworld-manager-role created
clusterrole.rbac.authorization.k8s.io/helloworld-proxy-role created
Warning: rbac.authorization.k8s.io/v1beta1 ClusterRole is deprecated in v1.17+, unavailable in v1.22+; use rbac.authorization.k8s.io/v1 ClusterRole
clusterrole.rbac.authorization.k8s.io/helloworld-metrics-reader created
rolebinding.rbac.authorization.k8s.io/helloworld-leader-election-rolebinding created
clusterrolebinding.rbac.authorization.k8s.io/helloworld-manager-rolebinding created
clusterrolebinding.rbac.authorization.k8s.io/helloworld-proxy-rolebinding created
service/helloworld-controller-manager-metrics-service created
deployment.apps/helloworld-controller-manager created
```

## 卸载和清理

- 体验完毕后，如果想把前面创建的资源和CRD全部清理掉，可以执行以下命令：

```shell
cd $GOPATH/src/helloworld
make uninstall
```

 至此，通过kubebuilder创建Operator相关资源的基本流程，咱们已经体验过一遍了，本篇以熟悉工具和流程为主，并未体验到Operator实质性的强大功能，这些都留待后面的章节吧，咱们逐步深入学习实践；