package output

import (
	"fmt"
	"os/exec"
	"testing"
)

func TestOut(t *testing.T) {
	stdoutChan := make(chan string, 100)
	incr := 0
	go func() {
		for line := range stdoutChan {
			incr++
			fmt.Println(incr, line)
		}
	}()

	cmd := exec.Command("bash", "-c", "echo 123;sleep 1;echo 456; echo 789")
	stdout := NewOutputStream(stdoutChan)
	cmd.Stdout = stdout
	cmd.Run()

	select {

	}
}

func TestCheckBuffer(t *testing.T) {
	cmd := exec.Command("bash", "-c", "echo 123")
	stdout := NewOutputBuffer()
	cmd.Stdout = stdout
	cmd.Run()

	fmt.Println(stdout.buf.String())


}
