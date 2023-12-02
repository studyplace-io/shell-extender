package pod_exec_command

import (
	"log"
	"testing"
)

func TestHandleCommand(t *testing.T) {

	cmd := NewExecPodContainerCmd("./config1", "test-pod",
		"my-container", "default", true)
	err := cmd.Run([]string{"sh", "-c", "ls -a"})
	if err != nil {
		log.Fatal(err)
	}
}
