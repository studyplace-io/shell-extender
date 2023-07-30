package main

import (
	"fmt"
	"github.com/practice/shell_extender/pkg/command"
	"github.com/practice/shell_extender/pkg/remote_command"
)

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

}
