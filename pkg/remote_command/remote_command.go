package remote_command

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/practice/shell_extender/pkg/waitgroup_timeout"
	"golang.org/x/crypto/ssh"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

// BatchRunRemoteNodeFromConfig 使用配置文件远端执行shell命令，不返回执行结果
func BatchRunRemoteNodeFromConfig(config string, cmd string) error {

	remoteNode, err := loadConfig(config)
	if err != nil {
		return ErrConfigParse
	}

	for _, v := range remoteNode.RemoteNodes {
		batchWg.Add(1)
		if v.Port == "" {
			v.Port = defaultPort
		}
		p, _ := strconv.Atoi(v.Port)
		go runRemoteNodeFromConfig(v.User, v.Password, v.Host, p, cmd)
	}

	batchWg.Wait()
	return nil

}

var (
	ErrTimeout     = errors.New("exec command timeout")
	ErrConfigParse = errors.New("config parse error")
)

// BatchRunRemoteNodeFromConfigWithTimeout 使用配置文件远端执行shell命令并提供超时时间
func BatchRunRemoteNodeFromConfigWithTimeout(config string, cmd string, timeout int64) error {

	remoteNode, err := loadConfig(config)
	if err != nil {
		return ErrConfigParse
	}

	wg := waitgroup_timeout.NewWaitGroupWithTimeout(time.Duration(timeout) * time.Second)

	for _, v := range remoteNode.RemoteNodes {
		wg.Add(1)
		if v.Port == "" {
			v.Port = defaultPort
		}
		p, _ := strconv.Atoi(v.Port)
		go runRemoteNodeFromConfigWithTimeout(v.User, v.Password, v.Host, p, cmd, wg)
	}

	if wg.WaitTimeout() {
		return ErrTimeout
	}

	return nil

}

// RunRemoteNode 远端执行shell命令并提供超时时间
func RunRemoteNode(user, password, host string, port int, cmd string) error {
	return runRemoteNode(user, password, host, port, cmd)
}

// RunRemoteNodeWithTimeout 远端执行shell命令
func RunRemoteNodeWithTimeout(user, password, host string, port int, cmd string, timeout int64) error {
	notifyC := make(chan struct{})
	errC := make(chan error, 1)
	var err error
	execFunc := func() {
		err = runRemoteNode(user, password, host, port, cmd)
		if err != nil {
			errC <- err
		}
		close(notifyC)
	}
	go execFunc()

	// 超时执行返回
	t := time.Duration(timeout) * time.Second
	select {
	case <-notifyC:
		if err != nil {
			return <-errC
		}
		return nil
	case <-time.After(t):
		return ErrTimeout
	}
}

// sSHConnect 使用ssh登入
func sSHConnect(user, password, host string, port int) (*ssh.Session, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		session      *ssh.Session
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	hostKeyCallback := func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		return nil
	}

	clientConfig = &ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		Timeout:         30 * time.Second,
		HostKeyCallback: hostKeyCallback,
	}

	// connect to ssh
	addr = fmt.Sprintf("%s:%d", host, port)

	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create session
	if session, err = client.NewSession(); err != nil {
		return nil, err
	}

	return session, nil
}

var batchWg sync.WaitGroup

const defaultPort = "22"

func runRemoteNodeFromConfig(user, password, host string, port int, cmd string) {
	defer batchWg.Done()
	runRemoteNode(user, password, host, port, cmd)
}

func runRemoteNodeFromConfigWithTimeout(user, password, host string, port int, cmd string, wg *waitgroup_timeout.WaitGroupWithTimeout) {
	defer wg.Done()
	runRemoteNode(user, password, host, port, cmd)
}

// runRemoteNode 对远程节点执行命令
func runRemoteNode(user, password, host string, port int, cmd string) error {
	session, err := sSHConnect(user, password, host, port)
	if err != nil {
		return err
	}
	defer session.Close()

	var stdOut, stdErr bytes.Buffer
	session.Stdout = &stdOut
	session.Stderr = &stdErr

	session.Start(cmd)
	session.Wait()
	fmt.Println("host ip exec result: ", host)
	if stdErr.String() != "" {
		fmt.Println("exec result get error")
		fmt.Println(string(stdErr.Bytes()))
	} else {
		fmt.Println(string(stdOut.Bytes()))
	}

	return nil
}

// RunRemoteCommandLine 登入远程节点命令行
func RunRemoteCommandLine(ip string, port int, user, password string) {
	addr := fmt.Sprintf("%s:%s", ip, strconv.Itoa(port))
	// 建立SSH客户端连接
	client, err := ssh.Dial("tcp", addr, &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		fmt.Printf("SSH dial error: %s", err.Error())
		return
	}

	// 建立新会话
	session, err := client.NewSession()
	if err != nil {
		fmt.Printf("new session error: %s", err.Error())
		return
	}
	defer session.Close()
	session.Stdout = os.Stdout // 会话输出关联到系统标准输出设备
	session.Stderr = os.Stderr // 会话错误输出关联到系统标准错误输出设备
	session.Stdin = os.Stdin   // 会话输入关联到系统标准输入设备
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // 禁用回显（0禁用，1启动）
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, //output speed = 14.4kbaud
	}
	if err = session.RequestPty("linux", 32, 160, modes); err != nil {
		fmt.Printf("request pty error: %s", err.Error())
	}
	if err = session.Shell(); err != nil {
		fmt.Printf("start shell error: %s", err.Error())
	}
	if err = session.Wait(); err != nil {
		fmt.Printf("return error: %s", err.Error())
	}
}
