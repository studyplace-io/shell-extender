package remote_command

import (
	"fmt"
	"testing"
)

func TestRunRemoteNodeFromConfig(t *testing.T) {
	err := BatchRunRemoteNodeFromConfig("./shell_extender/config.yaml", "kubect get pods")
	if err != nil {
		fmt.Println(err)
	}
	err = BatchRunRemoteNodeFromConfigWithTimeout("./shell_extender/config.yaml", "sleep 3; echo 123", 1)
	if err != nil {
		fmt.Println(err)
	}
}
