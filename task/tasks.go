package task

import (
	"bufio"
	"fmt"
	log "github.com/ngmoco/timber"
	"io"
	"os/exec"
	"strings"
)

const (
	PENDING  = "PENDING"
	ERROR    = "ERROR"
	COMPLETE = "COMPLETE"
	listKey  = "tasks"
)

type TaskStatus struct {
	Id          string
	Status      string
	Description string
}

func (t *TaskStatus) redisListKey() string {
	return fmt.Sprintf("task:%s:lines", t.Id)
}

func (t *TaskStatus) redisStatusKey() string {
	return statusKeyFromId(t.Id)
}

func (t *TaskStatus) redisLineKey() string {
	return fmt.Sprintf("task:%s:lines", t.Id)
}

func statusKeyFromId(id string) string {
	return fmt.Sprintf("task:%s:status", id)
}

func ExecuteStringTask(cmd string) (string, error) {
	log.Info("Executing task: %s", cmd)
	return executeTask(stringToCmd(cmd))
}

func executeTask(cmd *exec.Cmd) (string, error) {
	s := make([]string, 0, 10)
	commandOut := make(chan string)
	status := &TaskStatus{generateId(), PENDING, ""}
	go runCommand(cmd, commandOut, status)
	for line := range commandOut {
		s = append(s, line)
	}
	return strings.Join(s, ""), nil
}

func runCommand(cmd *exec.Cmd, commandOut chan string, status *TaskStatus) {
	defer close(commandOut)
	status.Status = PENDING
	outPipe, err := cmd.StdoutPipe()
	// TODO read from stderr as well
	//errPipe, err := cmd.StderrPipe()
	if err != nil {
		log.Warn("Error reading command output %s", err)
		status.Status = ERROR
		return
	}
	cmd.Start()
	readPipeOutput(outPipe, commandOut)
	// TODO write status code to output
	// Error on wait could be *ExitError
	err = cmd.Wait()
	if err != nil {
		status.Status = ERROR
	} else {
		status.Status = COMPLETE
	}
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
