package command

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os/exec"
	"time"
)

var (
	ErrTimeout = errors.New("exec command timeout")
)

// ExecShellCommand 执行命令
// 输出：1.命令行结果 2.进程输出code 3.错误
func ExecShellCommand(cmd string) (string, int, error) {
	executor := exec.Command("bash", "-c", cmd)
	outByte, err := executor.CombinedOutput()
	out := string(outByte)
	return out, executor.ProcessState.ExitCode(), err
}

// ExecShellCommandWithResult 执行命令并输出所有结果
// 输出：1.命令行结果 2.命令行错误 3.进程输出code 4.错误
func ExecShellCommandWithResult(cmd string) (string, string, int, error) {

	executor := exec.Command("bash", "-c", cmd)
	var (
		stdout, stderr bytes.Buffer
		err            error
	)
	executor.Stdout = &stdout
	executor.Stderr = &stderr
	err = executor.Start()
	if err != nil {
		return string(stdout.Bytes()), string(stderr.Bytes()), executor.ProcessState.ExitCode(), err
	}

	err = executor.Wait()
	return string(stdout.Bytes()), string(stderr.Bytes()), executor.ProcessState.ExitCode(), err
}

// ExecShellCommandWithTimeout 执行命令并超时时间
// 输出：1.命令行结果 2.命令行错误 3.进程输出code 4.错误
func ExecShellCommandWithTimeout(cmd string, timeout int64) (string, string, int, error) {
	executor := exec.Command("bash", "-c", cmd)
	// executor.Run() 会阻塞，因此开一个goroutine异步执行，
	// 当执行结束时，使用chan通知
	notifyC := make(chan struct{})
	var err error
	execFunc := func() {
		err = executor.Run()
		close(notifyC)
	}
	go execFunc()

	var (
		stdout, stderr bytes.Buffer
	)
	executor.Stdout = &stdout
	executor.Stderr = &stderr

	if err != nil {
		return string(stdout.Bytes()), string(stderr.Bytes()), executor.ProcessState.ExitCode(), err
	}

	// 超时执行返回
	t := time.Duration(timeout) * time.Second
	select {
	case <-notifyC:
		return string(stdout.Bytes()), string(stderr.Bytes()), executor.ProcessState.ExitCode(), err
	case <-time.After(t):
		return string(stdout.Bytes()), string(stderr.Bytes()), executor.ProcessState.ExitCode(), ErrTimeout
	}
}

// ExecShellCommandWithChan 执行命令并使用管道输出
// 输入：chan 输出：错误
func ExecShellCommandWithChan(cmd string, queue chan string) error {
	executor := exec.Command("bash", "-c", cmd)
	stdout, err := executor.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := executor.StderrPipe()
	if err != nil {
		return err
	}

	executor.Start()

	callbackFunc := func(in io.ReadCloser) {
		reader := bufio.NewReader(in)
		for {
			line, _, err := reader.ReadLine()
			if err != nil || io.EOF == err {
				break
			}

			select {
			case queue <- string(line):
			}
		}
	}

	go callbackFunc(stdout)
	go callbackFunc(stderr)

	executor.Wait()
	close(queue)
	return nil
}
