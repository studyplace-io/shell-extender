## shell-extender
### 项目简介：对os/exec包进行扩展，增加其功能。
- ExecShellCommand 执行命令
- ExecShellCommandWithResult 执行命令并输出所有结果
- ExecShellCommandWithTimeout 执行命令并超时时间
- ExecShellCommandWithChan 执行命令并使用管道输出 

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
}


```