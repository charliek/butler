package task

import (
	"fmt"
	log "github.com/ngmoco/timber"
	"os/exec"
	"strings"
)

func NewTask(desc string) *Task {
	return &Task{
		Description: desc,
		Cmds:        make([]Command, 0, 10),
	}
}

type Task struct {
	Description string
	Cmds        []Command
}

func (t *Task) AddCommand(cmd Command) {
	t.Cmds = append(t.Cmds, cmd)
}

func (t *Task) AddStringCommand(cmd string, fail bool) {
	c := Command{
		FailOnError: fail,
		Cmd:         strings.Split(cmd, " "),
	}
	t.Cmds = append(t.Cmds, c)
}

type Command struct {
	FailOnError bool
	Cmd         []string
}

func ExecuteRecordedTask(task *Task) string {
	id := generateId()
	log.Info("Beginning to run task '%s' with id '%s'", task.Description, id)
	go executeWithRedis(id, task)
	return id
}

func executeWithRedis(id string, task *Task) {
	status := &TaskStatus{id, PENDING, task.Description}
	initializeStatus(status)
	for _, cmd := range task.Cmds {
		recordLine(status, fmt.Sprintf("Running command: %s", strings.Join(cmd.Cmd, " ")))
		recordLine(status, "***************************************")
		err := executeRecordedTask(status, stringsToCmd(cmd.Cmd...))
		if err != nil && cmd.FailOnError {
			status.Status = ERROR
			break
		}
	}
	if status.Status == PENDING {
		status.Status = COMPLETE
	}
	updateStatus(status)
}

func executeRecordedTask(status *TaskStatus, cmd *exec.Cmd) error {
	commandOut := make(chan string)
	go runCommand(cmd, commandOut, status)
	for line := range commandOut {
		// TODO remove info logging once things are working better
		log.Info(line)
		recordLine(status, line)
	}
	return nil
}
