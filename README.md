# operator-demo
hello operator


### 创建项目

```bash
➜  youaij operator-sdk new operator-demo
INFO[0000] Creating new Go operator 'operator-demo'.
INFO[0000] Created go.mod
INFO[0000] Created tools.go
INFO[0000] Created cmd/manager/main.go
INFO[0000] Created build/Dockerfile
INFO[0000] Created build/bin/entrypoint
INFO[0000] Created build/bin/user_setup
INFO[0000] Created deploy/service_account.yaml
INFO[0000] Created deploy/role.yaml
INFO[0000] Created deploy/role_binding.yaml
INFO[0000] Created deploy/operator.yaml
INFO[0000] Created pkg/apis/apis.go
INFO[0000] Created pkg/controller/controller.go
INFO[0000] Created version/version.go
INFO[0000] Created .gitignore
INFO[0000] Validating project
go: downloading github.com/operator-framework/operator-sdk v0.17.0
go: downloading k8s.io/client-go v0.17.4
go: downloading sigs.k8s.io/controller-runtime v0.5.2
go: downloading k8s.io/api v0.17.4
go: downloading k8s.io/apimachinery v0.17.4
go: downloading golang.org/x/crypto v0.0.0-20200220183623-bac4c82f6975
go: downloading golang.org/x/net v0.0.0-20200226121028-0de0cce0169b
go: downloading github.com/coreos/prometheus-operator v0.38.0
go: downloading k8s.io/utils v0.0.0-20191114200735-6ca3b61696b6
go: downloading golang.org/x/sys v0.0.0-20200122134326-e047566fdf82
INFO[0030] Project validation successful.
INFO[0030] Project creation complete.
```

**创建新的api资源**（create crcd）

```bash
➜  youaij cd operator-demo
➜  operator-demo operator-sdk add api --api-version=app.learn.com/v1 --kind=Learn
INFO[0000] Generating api version app.learn.com/v1 for kind Learn.
INFO[0000] Created pkg/apis/app/group.go
INFO[0002] Created pkg/apis/app/v1/learn_types.go
INFO[0002] Created pkg/apis/addtoscheme_app_v1.go
INFO[0002] Created pkg/apis/app/v1/register.go
INFO[0002] Created pkg/apis/app/v1/doc.go
INFO[0002] Created deploy/crds/app.learn.com_v1_learn_cr.yaml
INFO[0002] Running deepcopy code-generation for Custom Resource group versions: [app:[v1], ]
INFO[0009] Code-generation complete.
INFO[0009] Running CRD generator.
INFO[0013] CRD generation complete.
INFO[0013] API generation complete.
➜  operator-demo
```

可以发现pkg/apis/、deploy/crds等目录下有新文件生成。

**创建控制器**

```bash
➜  operator-demo operator-sdk add controller --api-version=app.learn.com/v1 --kind=Learn
INFO[0000] Generating controller version app.learn.com/v1 for kind Learn.
INFO[0000] Created pkg/controller/learn/learn_controller.go
INFO[0000] Created pkg/controller/add_learn.go
INFO[0000] Controller generation complete.
➜  operator-demo
```

**生成可执行文件 && 构建镜像**

执行`operator-sdk build` 命令后，会在`build/_output/bin/` 目录下生成可执行文件：`operator-demo` ，并且会执行Dockerfile将可执行文件构建成镜像。

```bash
➜  operator-demo operator-sdk build ccr.ccs.tencentyun.com/youaijj/operator-demo
INFO[0025] Building OCI image ccr.ccs.tencentyun.com/youaijj/operator-demo 
Sending build context to Docker daemon  43.45MB
Step 1/7 : FROM registry.access.redhat.com/ubi8/ubi-minimal:latest
latest: Pulling from ubi8/ubi-minimal
0fd3b5213a9b: Pull complete 
aebb8c556853: Pull complete 
Digest: sha256:5cfbaf45ca96806917830c183e9f37df2e913b187aadb32e89fd83fa455ebaa6
Status: Downloaded newer image for registry.access.redhat.com/ubi8/ubi-minimal:latest
 ---> 28095021e526
Step 2/7 : ENV OPERATOR=/usr/local/bin/operator-demo     USER_UID=1001     USER_NAME=operator-demo
 ---> Running in 8b12a0d955a6
Removing intermediate container 8b12a0d955a6
 ---> 8bed7d73a4ef
Step 3/7 : COPY build/_output/bin/operator-demo ${OPERATOR}
 ---> c0005e4c8988
Step 4/7 : COPY build/bin /usr/local/bin
 ---> 804eb6a1ebc5
Step 5/7 : RUN  /usr/local/bin/user_setup
 ---> Running in 6df48777bf5a
+ echo 'operator-demo:x:1001:0:operator-demo user:/root:/sbin/nologin'
+ mkdir -p /root
+ chown 1001:0 /root
+ chmod ug+rwx /root
+ rm /usr/local/bin/user_setup
Removing intermediate container 6df48777bf5a
 ---> 7fc2618c58c9
Step 6/7 : ENTRYPOINT ["/usr/local/bin/entrypoint"]
 ---> Running in af7ce56f695e
Removing intermediate container af7ce56f695e
 ---> a38d25c5a971
Step 7/7 : USER ${USER_UID}
 ---> Running in 487165095ddd
Removing intermediate container 487165095ddd
 ---> 79ea87700f15
Successfully built 79ea87700f15
Successfully tagged ccr.ccs.tencentyun.com/youaijj/operator-demo:latest
INFO[0062] Operator build complete.                     
➜  operator-demo
```

**上传镜像**

```bash
➜  ~ docker push ccr.ccs.tencentyun.com/youaijj/operator-demo:latest
The push refers to repository [ccr.ccs.tencentyun.com/youaijj/operator-demo]
acee948057db: Pushed
8ec8f1c519e2: Pushed
4008b7edb24c: Pushed
3485805ce47c: Pushed
b0e2911c99f3: Pushed
latest: digest: sha256:9bf7f007b68dbae7c357e51504dea34d4fa2c7674ed539ee4834167a25501183 size: 1363
➜  ~
```

**RBAC权限创建与绑定**

```bash
➜  cd deploy
➜  deploy kubectl apply -f role.yaml                            
role.rbac.authorization.k8s.io/operator-demo created
➜  deploy kubectl apply -f service_account.yaml 
serviceaccount/operator-demo created
➜  deploy kubectl apply -f role_binding.yaml   
rolebinding.rbac.authorization.k8s.io/operator-demo created
➜  deploy

➜  ~ kubectl get role
NAME            AGE
operator-demo   8m23s
➜  ~ kubectl get serviceaccounts
NAME            SECRETS   AGE
default         1         293d
operator-demo   1         27m

```

**部署应用（启动deployment）**

```bash
➜  deploy kubectl apply -f operator.yaml 
deployment.apps/operator-demo created

## 执行operator.yaml 之后发现容器不断在在重启，查找原因
➜  ~ docker run -it 79ea87700f15
{"level":"info","ts":1603882995.8388538,"logger":"cmd","msg":"Operator Version: 0.0.1"}
{"level":"info","ts":1603882995.8389218,"logger":"cmd","msg":"Go Version: go1.13.3"}
{"level":"info","ts":1603882995.838946,"logger":"cmd","msg":"Go OS/Arch: linux/amd64"}
{"level":"info","ts":1603882995.8390565,"logger":"cmd","msg":"Version of operator-sdk: v0.17.0"}
{"level":"error","ts":1603882995.8390841,"logger":"cmd","msg":"Failed to get watch namespace","error":"WATCH_NAMESPACE must be set","stacktrace":"github.com/go-logr/zapr.(*zapLogger).Error\n\t/Users/chenjunjie/go/pkg/mod/github.com/go-logr/zapr@v0.1.1/zapr.go:128\nmain.main\n\toperator-demo/cmd/manager/main.go:76\nruntime.main\n\t/usr/local/go/src/runtime/proc.go:203"}
➜  ~
## 日志中提示：Failed to get watch namespace","error":"WATCH_NAMESPACE must be set"

➜  deploy cd crds  
➜  crds kubectl apply -f app.learn.com_learns_crd.yaml 
customresourcedefinition.apiextensions.k8s.io/learns.app.learn.com created
➜  crds

## 创建crd后，容器重启成功，不会不断重启，容器运行正常
➜  ~ kubectl get pod
NAME                                     READY   STATUS    RESTARTS   AGE
operator-demo-5f64bd68dc-slmkv           1/1     Running   7          8m13s
➜  ~
```

**执行CR**  

```bash
## 执行CR前查容器状态
➜  ~ kubectl get pod
NAME                                     READY   STATUS    RESTARTS   AGE
operator-demo-5f64bd68dc-slmkv           1/1     Running   7          12m13s
➜  ~

## 执行CR
➜  crds kubectl apply -f app.learn.com_v1_learn_cr.yaml 
learn.app.learn.com/example-learn created
➜  crds

## 执行CR后查容器状态
➜  ~ kubectl get pod
NAME                                     READY   STATUS    RESTARTS   AGE
example-learn-6764b9858-8g578            1/1     Running   0          18s
example-learn-6764b9858-mjq4c            1/1     Running   0          18s
example-learn-6764b9858-x78bw            1/1     Running   0          18s
operator-demo-5f64bd68dc-slmkv           1/1     Running   7          13m46s
```

**CR功能演示**

直接修改 `kind: Learn` && `metadata.name: example-learn` 这个对象

```bash
➜  ~ kubectl edit learn example-learn
```

编辑器中显示如下:

```bash
# Please edit the object below. Lines beginning with a '#' will be ignored,
# and an empty file will abort the edit. If an error occurs while saving this file will be
# reopened with the relevant failures.
#
apiVersion: app.learn.com/v1
kind: Learn
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"app.learn.com/v1","kind":"Learn","metadata":{"annotations":{},"name":"example-learn","namespace":"default"},"spec":{"size":3}}
  creationTimestamp: "2020-10-28T11:26:00Z"
  generation: 1
  name: example-learn
  namespace: default
  resourceVersion: "11989305504"
  selfLink: /apis/app.learn.com/v1/namespaces/default/learns/example-learn
  uid: 6ce9f570-c9e4-4564-9199-6d8943c4e0ac
spec:
  size: 3
```

将`spec.size`修改4，保存查看容器状态：

```bash
➜  ~ kubectl get pod
NAME                                     READY   STATUS    RESTARTS   AGE
example-learn-6764b9858-8g578            1/1     Running   0          22m
example-learn-6764b9858-mjq4c            1/1     Running   0          22m
example-learn-6764b9858-rvpq2            1/1     Running   0          8s
example-learn-6764b9858-x78bw            1/1     Running   0          22m
operator-demo-5f64bd68dc-slmkv           1/1     Running   7          36m
```
