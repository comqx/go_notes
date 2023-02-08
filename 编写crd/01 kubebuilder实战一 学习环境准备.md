[toc]

版权声明：本文为CSDN博主「程序员欣宸」的原创文章，遵循CC 4.0 BY-SA版权协议，转载请附上原文出处链接及本声明。
原文链接：https://blog.csdn.net/boling_cavalry/article/details/113035349

# 关于kubebuilder

1. 在实际工作中，对kubernetes的资源执行各种个性化配置和控制是很常见的需求，例如自定义镜像的pod如何控制副本数、主从关系，以及各种自定义资源的控制等；
2. 对于上述需求，很适合使用Operator 模式来解决，这里有官方对Operator的介绍：https://kubernetes.io/zh/docs/concepts/extend-kubernetes/operator/ ，Operator模式的执行流程如下图所示：

![image-20211130133644516](https://picgo-img.oss-cn-beijing.aliyuncs.com/md-img/2021-11-30/1638250604.png)

3. 为了简化Operator开发，我们可以选用一些已有的开源工具，kubebuilder就是其中之一，《kubebuilder实战》系列就是对此工具由浅入深的实践；

# 本篇概览

作为《kubebuilder实战》系列的开篇，除了前面对kubebuilder的简单说明，还会列出整个实战的通用环境信息，以及涉及到的软件版本，然后再搭建好kubebuilder开发环境，总的来说需要做好以下准备工作，才能顺利开始kubebuilder的开发工作：

kubectl安装和配置，这样可以在kubebuilder电脑上操作kubernetes环境；

- 安装golang
- 安装docker
- 安装kustomize
- 安装kubebuilder

# 环境信息

如下图，整个实战环境一共由两台电脑组成：



kubernetes：上面运行着1.20版本的kubernetes，关于kubernetes的部署不是本文重点，请参考其他教程完成，需要确保kubernetes正常可用；
kubebuilder电脑：操作系统是CentoOS-7.9.2009，hostname是kubebuilder，咱们的实战就在这台电脑上操作；
kubebuilder版本：2.3.1
go版本：1.15.6
docker版本：19.03.13
为了省事儿，所有操作都是用root帐号执行的；

## kubectl安装和配置

>  自行安装，忽略

## 安装golang

>  自行安装，忽略

## 安装docker

> 自行安装，忽略

##  安装kustomize

```shell
mkdir -p $GOPATH/bin
cd $GOPATH/bin
GOBIN=$(pwd)/ GO111MODULE=on go get sigs.k8s.io/kustomize/kustomize/v3
```

## 安装kubebuilder

1. 以下脚本通过go命令确定当前系统和CPU架构，再去服务器下载对应的kubebuilder文件，然后设置环境变量：

```shell
os=$(go env GOOS)
arch=$(go env GOARCH)
curl -L https://go.kubebuilder.io/dl/2.3.1/${os}/${arch} | tar -xz -C /tmp/
mv /tmp/kubebuilder_2.3.1_${os}_${arch} /usr/local/kubebuilder
export PATH=$PATH:/usr/local/kubebuilder/bin
```

1. 执行命令确认安装成功：

```shell
[~/tmp/dev]$ kubebuilder version
Version: version.Version{KubeBuilderVersion:"2.3.1", KubernetesVendor:"1.16.4", GitCommit:"8b53abeb4280186e494b726edf8f54ca7aa64a49", BuildDate:"2020-03-26T16:42:00Z", GoOs:"unknown", GoArch:"unknown"}
```

- 至此，kubebuilder开发环境的准备工作就完成了，接下来的章节咱们正式进入开发实战，一起去学习精彩的kubernetes operator；
- 