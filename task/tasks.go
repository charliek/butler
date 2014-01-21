package task

import (
	"bufio"
	"crypto/rand"
	"fmt"
	log "github.com/ngmoco/timber"
	"io"
	"os/exec"
	"strings"
)

func generateId() string {
	buf := make([]byte, 16)
	io.ReadFull(rand.Reader, buf)
	return fmt.Sprintf("%x", buf)
}

type Command struct {
	FailOnError bool
	cmd         []string
}

type Task struct {
	Description string
	cmds        []Command
}

func ExecuteStringTask(cmd string) (string, error) {
	c := strings.Split(cmd, " ")
	return ExecuteTask(c...)
}

func ExecuteTask(arg ...string) (string, error) {
	return executeTask(arg[0], arg[1:]...)
}

func executeTask(name string, arg ...string) (string, error) {
	log.Debug("Executing : %s %s\n", name, strings.Join(arg, " "))
	cmd := exec.Command(name, arg...)
	s := make([]string, 0, 10)
	commandOut := make(chan string)
	go runCommand(cmd, commandOut)
	for line := range commandOut {
		s = append(s, line)
	}
	return strings.Join(s, ""), nil
}

func runCommand(cmd *exec.Cmd, commandOut chan string) error {
	defer close(commandOut)
	outPipe, err := cmd.StdoutPipe()
	// TODO read from stderr as well
	//errPipe, err := cmd.StderrPipe()
	if err != nil {
		log.Warn("Error reading command output %s", err)
		return err
	}
	cmd.Start()
	readPipeOutput(outPipe, commandOut)
	// TODO write status code to output
	// Error on wait could be *ExitError
	return cmd.Wait()
}

func readPipeOutput(pipe io.ReadCloser, commandOut chan string) {
	stdout := bufio.NewReader(pipe)
	for {
		line, err := stdout.ReadString('\n')
		if err == nil || err == io.EOF {
			if len(line) > 0 {
				commandOut <- line
			}
		}
		if err != nil {
			break
		}
	}
}
