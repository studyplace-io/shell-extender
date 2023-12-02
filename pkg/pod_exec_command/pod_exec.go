package pod_exec_command

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"os"
)

type ExecPodContainerCmd struct {
	// kubeConfig kubeconfig目录地址，默认集群内.kube/config
	kubeConfig string
	// insecure 是否跳过tls鉴权
	insecure      bool
	podName       string
	containerName string
	namespace     string
	k8sClient     kubernetes.Interface
}

func NewExecPodContainerCmd(kubeConfig string, podName string, containerName string, namespace string, insecure bool) *ExecPodContainerCmd {
	return &ExecPodContainerCmd{kubeConfig: kubeConfig,
		podName: podName, containerName: containerName,
		namespace: namespace, insecure: insecure}
}

// prepare 准备登入材料
func (epc *ExecPodContainerCmd) prepare() (*rest.Config, kubernetes.Interface, error) {
	config, err := clientcmd.BuildConfigFromFlags("", epc.kubeConfig)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "parse kube config file error: %s", err)
	}
	config.Insecure = epc.insecure
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "prepare kube client error: %s", err)
	}
	epc.k8sClient = client
	return config, client, nil
}

// validatePod 验证 pod container
func (epc *ExecPodContainerCmd) validatePod(ctx context.Context, podName, containerName, namespace string) error {
	pod, err := epc.k8sClient.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
	if err != nil && k8serrors.IsNotFound(err) {
		return fmt.Errorf("pod %s/%s not found", namespace, podName)
	}

	if err != nil {
		return err
	}

	if pod.Status.Phase == corev1.PodSucceeded || pod.Status.Phase == corev1.PodFailed {
		return fmt.Errorf("cannot exec into container in a completed pod, current phase %s", pod.Status.Phase)
	}

	for _, cc := range pod.Spec.InitContainers {
		if containerName == cc.Name {
			return fmt.Errorf("can't exec init container %s in pod %s/%s ", containerName, namespace, podName)
		}
	}

	for _, cs := range pod.Status.ContainerStatuses {
		if containerName == cs.Name {
			return nil
		}
	}

	return fmt.Errorf("pod has no container %s", containerName)
}

// Run 执行远程命令
func (epc *ExecPodContainerCmd) Run(command []string) error {
	option := &corev1.PodExecOptions{
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

	// 校验 pod container 是否存在
	err = epc.validatePod(context.Background(), epc.podName, epc.containerName, epc.namespace)
	if err != nil {
		return errors.Wrapf(err, "validate Pod error: %s", err)
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
		return errors.Wrapf(err, "SPDY Exec error: %s", err)
	}

	return exec.StreamWithContext(context.Background(), remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Tty:    true,
	})
}
