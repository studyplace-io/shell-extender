package command

import (
	"fmt"
	"testing"
)

func TestExecShellCommand(t *testing.T) {
	out, i, err := ExecShellCommand("echo TestExecShellCommand")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("i: %v, out: %s", i, out)

}

func TestExecShellCommandWithResult(t *testing.T) {
	stdout, stderr, code, err := ExecShellCommandWithResult("echo TestExecShellCommandWithResult")
	fmt.Printf("stdout: %v, stderr: %s, code: %v, err: %v\n", stdout, stderr, code, err)
}

func TestExecShellCommandWithTimeout(t *testing.T) {
	stdout, stderr, code, err := ExecShellCommandWithTimeout("sleep 15; echo TestExecShellCommandWithTimeout; kubectl get pods", 20)
	fmt.Printf("stdout: %v, stderr: %s, code: %v, err: %v\n", stdout, stderr, code, err)
}

func TestExecShellCommandWithChan(t *testing.T) {
	outputC := make(chan string, 10)

	go func() {
		for i := range outputC {
			fmt.Println("output line: ", i)
		}
	}()

	err := ExecShellCommandWithChan("echo TestExecShellCommandWithChan ;sleep 1;kubectl get node", outputC)
	if err != nil {
		fmt.Println(err)
		return
	}
}
