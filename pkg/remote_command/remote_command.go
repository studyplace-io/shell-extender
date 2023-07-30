package remote_command

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/practice/shell_extender/pkg/waitgroup_timeout"
	"golang.org/x/crypto/ssh"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

// BatchRunRemoteNodeFromConfig 远端执行
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
	ErrTimeout = errors.New("exec command timeout")
	ErrConfigParse = errors.New("config parse error")
)

// BatchRunRemoteNodeFromConfigWithTimeout 远端执行shell命令并提供超时时间
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
func runRemoteNode(user, password, host string, port int, cmd string) {
	session, err := sSHConnect(user, password, host, port)
	if err != nil {
		log.Fatal(err)
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
}