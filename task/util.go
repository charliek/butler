package task

import (
	"crypto/rand"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

func generateId() string {
	buf := make([]byte, 16)
	io.ReadFull(rand.Reader, buf)
	return fmt.Sprintf("%x", buf)
}

func stringsToCmd(arg ...string) *exec.Cmd {
	return exec.Command(arg[0], arg[1:]...)
}

func stringToCmd(cmd string) *exec.Cmd {
	c := strings.Split(cmd, " ")
	return stringsToCmd(c...)
}
