package pod_exec_command

import (
	"context"
	"errors"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"os"
)

var (
	ErrParseKubeConfig = errors.New("parse kube config file error")
	ErrPrepareClient   = errors.New("prepare kube client error")
	ErrSPDYExecutor    = errors.New("SPDY Exec error")
)

type ExecPodContainerCmd struct {
	// kubeConfig kubeconfig目录地址，默认集群内.kube/config
	kubeConfig    string
	// insecure 是否跳过tls鉴权
	insecure      bool
	podName       string
	containerName string
	namespace     string
}

func NewExecPodContainerCmd(kubeConfig string, podName string, containerName string, namespace string, insecure bool) *ExecPodContainerCmd {
	return &ExecPodContainerCmd{kubeConfig: kubeConfig,
		podName: podName, containerName: containerName,
		namespace: namespace, insecure: insecure}
}

// prepare 准备登入材料
func (epc *ExecPodContainerCmd) prepare() (*rest.Config, *kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", epc.kubeConfig)
	if err != nil {
		return nil, nil, ErrParseKubeConfig
	}
	config.Insecure = epc.insecure
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, ErrPrepareClient
	}
	return config, client, nil
}

// Run 执行远程命令
func (epc *ExecPodContainerCmd) Run(command []string) error {
	option := &v1.PodExecOptions{
		Container: epc.containerName,
		Command:   command,
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
		TTY:       true,
	}
	config, client, err := epc.prepare()
	if err != nil {
		return err
	}

	// 执行pods中 特定container容器的命令
	req := client.CoreV1().RESTClient().Post().Resource("pods").
		Namespace(epc.namespace).
		Name(epc.podName).
		SubResource("exec").
		Param("color", "false").
		VersionedParams(
			option,
			scheme.ParameterCodec,
		)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return ErrSPDYExecutor
	}

	return exec.StreamWithContext(context.Background(), remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Tty:    true,
	})
}
