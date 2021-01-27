# kubebuilder

> 官方：https://book.kubebuilder.io/quick-start.html  
>
> 好博客：https://www.cnblogs.com/alisystemsoftware/p/11580202.html
>
> https://www.huweihuang.com/kubernetes-notes/code-analysis/kube-controller-manager/sharedIndexInformer.html

## 自动生成crd项目

```shell
mkdir $GOPATH/src/example
kubebuilder init --domain tutorial.kubebuilder.io

或者直接创建出来api接口：
kubebuilder create api --group apps --version v1alpha1 --kind Application
```

