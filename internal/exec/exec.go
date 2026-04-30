package exec

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/Moorad/workforest/internal/utils"
	"github.com/fatih/color"
)

type Command struct {
	args   []string
	Cmd    *exec.Cmd
	stdout *bytes.Buffer
	stderr *bytes.Buffer
}

func NewCommand(args ...string) *Command {
	var stdout, stderr bytes.Buffer
	Cmd := exec.Command(args[0], args[1:]...)

	Cmd.Stdout = &stdout
	Cmd.Stderr = &stderr

	return &Command{
		args,
		Cmd,
		&stdout,
		&stderr,
	}
}

func (c *Command) Output() (string, error) {
	utils.Debug("$ %s\n", strings.Join(c.args, " "))

	err := c.Cmd.Run()
	if err != nil {
		errStr := c.stderr.String()
		utils.DebugError("> %s", errStr)
		return errStr, err
	}

	outStr := c.stdout.String()
	if outStr != "" {
		utils.DebugSuccess("> %s", outStr)
	}

	return outStr, nil
}

func (c *Command) Run() error {
	_, err := c.Output()

	return err
}

func SysCall(args ...string) error {
	gray := color.RGB(170, 170, 170)

	_, err := gray.Printf("$ %s", strings.Join(args, " "))
	if err != nil {
		return err
	}

	path, err := exec.LookPath(args[0])
	if err != nil {
		return err
	}
	return syscall.Exec(path, args, os.Environ())
}
