# 安装Telepresence

> 这个[工具](https://www.telepresence.io/reference/install)是在本地运行服务，同时可以将该服务连接到远程Kubernetes群集

```shell
brew cask install osxfuse
brew install datawire/blackbird/telepresence
```

## 运行telepresence例子

>1. 创建在端口8000上公开的名为hello-world的部署和服务,并确认已经就绪

```shell
$ kubectl apply -f https://raw.githubusercontent.com/telepresenceio/telepresence/master/docs/tutorials/hello-world.yaml
deployment.apps/hello-world created
service/hello-world created

$ kubectl get deployments
NAME                          READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/hello-world   1/1     1            1           6s
```

> 2. 通过telepresence运行一个可以访问该服务的Docker容器，即使该进程是本地的，但该服务正在Kubernetes集群中运行：

```shell
$ telepresence --docker-run --rm -it pstauffer/curl curl http://hello-world:8000/
[...]
T: Setup complete. Launching your container.
Hello, world!
T: Your process has exited.
[...]
```



# 安装kube-builder服务

>[kube-builder](https://book.kubebuilder.io/quick-start.html) 脚手架工具，快速创建一个api







