package exec

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/fatih/color"
)

func CmdOutput(args ...string) (string, error) {
	gray := color.RGB(170, 170, 170)

	_, err := gray.Printf("$ %s\n", strings.Join(args, " "))
	if err != nil {
		return "", err
	}
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		errStr := stderr.String()
		color.Red("> %s", errStr)
		return errStr, err
	}

	outStr := stdout.String()
	if outStr != "" {
		color.Green("> %s", outStr)
	}

	return outStr, nil
}

func Cmd(args ...string) error {
	_, err := CmdOutput(args...)

	return err
}

func SysCall(args ...string) error {
	gray := color.RGB(170, 170, 170)

	_, err := gray.Printf("$ %s\n", strings.Join(args, " "))
	if err != nil {
		return err
	}

	path, err := exec.LookPath(args[0])
	if err != nil {
		return err
	}
	return syscall.Exec(path, args, os.Environ())
}
