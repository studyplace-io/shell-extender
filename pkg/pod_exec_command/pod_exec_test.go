package pod_exec_command

import (
	"testing"
)

func TestHandleCommand(t *testing.T) {

	cmd := NewExecPodContainerCmd("./config", "myinspect-controller-69748dc6bf-84wdp",
		"myinspect-controller", "default", true)
	err := cmd.Run([]string{"sh", "-c", "ls -a"})
	if err != nil {
		return
	}
}
