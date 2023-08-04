package remote_command

import (
	"fmt"
	"testing"
)

func TestRunRemoteNodeFromConfig(t *testing.T) {
	err := BatchRunRemoteNodeFromConfig("/Users/bytedance/Desktop/code2/shell-extender/config.yaml", "kubectl get pods")
	if err != nil {
		fmt.Println(err)
	}
	err = BatchRunRemoteNodeFromConfigWithTimeout("/Users/bytedance/Desktop/code2/shell-extender/config.yaml", "sleep 3; echo 123", 1)
	if err != nil {
		fmt.Println(err)
	}

	err = RunRemoteNode("root", "", "", 0, "kubectl get pods")
	if err != nil {
		fmt.Println(err)
	}

	err = RunRemoteNodeWithTimeout("root", "", "", 0, "sleep 3; kubectl get pods", 2)
	if err != nil {
		fmt.Println(err)
	}
}
