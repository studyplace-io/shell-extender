## shell-extender
### 项目简介：对os/exec包进行扩展，增加其功能。
1. 执行shell命令
- ExecShellCommand 执行命令
- ExecShellCommandWithResult 执行命令并输出所有结果
- ExecShellCommandWithTimeout 执行命令并超时时间
- ExecShellCommandWithChan 执行命令并使用管道输出
2. 对远端节点执行shell命令
- RunRemoteNode
- RunRemoteNodeWithTimeout
3. 对批量远端节点执行shell命令
- BatchRunRemoteNodeFromConfig
- BatchRunRemoteNodeFromConfigWithTimeout
4. 支持命令行执行远程运维动作
- 可以将项目二进制编译更方便使用
5. 支持命令行登入远程节点
6. 支持对集群中特定pod的特定container执行命令
```yaml
remoteNodes:
  - host: 127.0.0.1 
    user: <user名>
    password: <password>
    port: 22  # 可选填，不填默认是ssh端口：22
  - host: 127.0.0.1
    user:
    password:
```

### 示例1 
**使用方法**
```go
func main() {
    fmt.Println("==============ExecShellCommand=================")
    out, i, err := command.ExecShellCommand("kubectl get node")
    if err != nil {
        fmt.Println(err)
        return
    }
    
    fmt.Printf("i: %v, out: %s", i, out)
    fmt.Println("===============================")
    
    fmt.Println("==============ExecShellCommandWithResult=================")
    stdout, stderr, code, err := command.ExecShellCommandWithResult("kubectl get node")
    fmt.Printf("stdout: %v, stderr: %s, code: %v, err: %v\n", stdout, stderr, code, err)
    fmt.Println("===============================")
    
    fmt.Println("==============ExecShellCommandWithTimeout=================")
    stdout, stderr, code, err = command.ExecShellCommandWithTimeout("sleep 10; kubectl get node", 3)
    fmt.Printf("stdout: %v, stderr: %s, code: %v, err: %v\n", stdout, stderr, code, err)
    fmt.Println("===============================")
    
    fmt.Println("==============ExecShellCommandWithChan=================")
    outputC := make(chan string, 10)
    
    go func() {
    for i := range outputC {
            fmt.Println("output line: ", i)
        }
    }()
    
    err = command.ExecShellCommandWithChan("kubectl get node", outputC)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println("===============================")

    fmt.Println("==============BatchRunRemoteNodeFromConfig=================")
    err = remote_command.BatchRunRemoteNodeFromConfig("./config.yaml", "kubectl get node")
    if err != nil {
        fmt.Println(err)
        return
    }
    
    fmt.Println("===============================")
    
    fmt.Println("==============ExecShellCommand=================")
    err = remote_command.BatchRunRemoteNodeFromConfigWithTimeout("./config.yaml", "kubectl get node", 10)
    if err != nil {
        fmt.Println(err)
        return
    }
    
    fmt.Println("===============================")
    
    
    cmd := NewExecPodContainerCmd("./config1", "myinspect-controller-69748dc6bf-84wdp",
                    "myinspect-controller", "default", true)
    err := cmd.Run([]string{"sh", "-c", "ls -a"})
    if err != nil {
        fmt.Println("aaaa: ", err)
        return
    }
    
}
```

### 示例2
**使用方法**
```bash
➜  shell_extender git:(main) ✗ go run main.go remoteExec --user=root --password=<password> --host=<host> --port=22 --script=./script.sh
host ip exec result:  xxxxxx
NAME                                  READY   STATUS    RESTARTS   AGE
ddd-55c668c8ff-v45lw                  1/1     Running   0          197d
example-deployment-658789c5cd-l4kwp   1/1     Running   0          56d
example-deployment-658789c5cd-qkn8d   1/1     Running   1          104d
example-pod                           1/1     Running   3          104d
jiangjiang-76fb44d88-pdzqv            1/1     Running   1          261d
k8splay1-59c7f5b4cb-mbw5v             2/2     Running   6          263d
k8splay1-59c7f5b4cb-x964b             2/2     Running   6          263d
kkkk-58cb7984db-mmffd                 1/1     Running   1          197d
my-deployment-5966cb4d75-cgm8f        1/1     Running   0          56d
my-deployment-5966cb4d75-df58x        1/1     Running   1          113d
my-deployment-5966cb4d75-gcg5w        1/1     Running   1          113d
myapp-rs-ftv7x                        1/1     Running   0          56d
myapp-rs-kr487                        1/1     Running   1          289d
mycrd-controller-78b98dcd7-hls78      1/1     Running   1          199d
mycrd-controller-78b98dcd7-n7ctv      1/1     Running   1          197d
mycsi-nginx-64c7d9cb77-mj7q7          1/1     Running   0          61d
mypod                                 1/1     Running   0          113d
myredis-0                             1/1     Running   0          56d
myredis-1                             1/1     Running   2          233d
nfscsi-nginx-768ff5bf55-jjfwx         1/1     Running   0          56d
ss-web-0                              1/1     Running   0          56d
ss-web-1                              1/1     Running   1          289d

➜  shell_extender git:(main) ✗ go run main.go remoteCommandLine --user=root --password=<password> --host=<host> --port=22
Last failed login: Fri Aug  4 21:59:12 CST 2023 from  on ssh:notty
There were 11 failed login attempts since the last successful login.
Last login: Fri Aug  4 21:57:08 2023 from 101.207.203.124
[root@VM-0-16-centos ~]# 
```