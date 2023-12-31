package main

import (
	"fmt"
	"github.com/practice/shell_extender/pkg/command"
	"github.com/practice/shell_extender/pkg/pod_exec_command"
	"github.com/practice/shell_extender/pkg/remote_command"
	"log"
)

func main() {
	fmt.Println("==============ExecShellCommand=================")
	out, i, err := command.ExecShellCommand("kubectl get node")
	if err != nil {
		fmt.Println(err)
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

	fmt.Println("==============RunRemoteNode=================")
	err = remote_command.RunRemoteNode("root", "", "", 0, "kubectl get pods")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("===============================")

	fmt.Println("==============RunRemoteNodeWithTimeout=================")
	err = remote_command.RunRemoteNodeWithTimeout("root", "", "", 0, "sleep 3; kubectl get pods", 2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("===============================")

	fmt.Println("==============ExecPodContainerCmd=================")
	cmd := pod_exec_command.NewExecPodContainerCmd("./pkg/pod_exec_command/config1", "test-pod",
		"my-container", "default", true)
	err = cmd.Run([]string{"sh", "-c", "ls -a"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("===============================")
}
