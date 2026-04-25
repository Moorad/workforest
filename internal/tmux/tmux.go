package tmux

import (
	"os"
	"strings"

	"github.com/Moorad/workforest/internal/exec"
)

type Session struct {
	name string
}

func Health() (bool, string) {
	out, err := exec.CmdOutput("tmux", "-V")
	if err != nil {
		return false, ""
	}

	return true, strings.Trim(strings.Split(out, " ")[1], "\n")
}

func NewSession(name string) (*Session, error) {
	session := Session{
		name,
	}

	cmd := []string{"tmux", "new-session", "-d", "-s", name}

	return &session, exec.Cmd(cmd...)
}

func (s *Session) RenameWindowMain(name string) error {
	return exec.Cmd("tmux", "rename-window", "-t", s.name, name)
}

func (s *Session) NewWindow(name string) error {
	return exec.Cmd("tmux", "new-window", "-t", s.name, "-n", name)
}

func (s *Session) SendKeys(windowName string, keys ...string) error {
	cmd := []string{"tmux", "send-keys", "-t", s.name + ":" + windowName}
	cmd = append(cmd, keys...)

	return exec.Cmd(cmd...)
}

func (s *Session) Attach() error {
	return exec.SysCall("tmux", "attach", "-t="+s.name)
}

func (s *Session) Switch() error {
	return exec.SysCall("tmux", "switch", "-t="+s.name)
}

func (s *Session) SwitchOrAttach() error {
	termProgram, isFound := os.LookupEnv("TERM_PROGRAM")

	if isFound && termProgram == "tmux" {
		return s.Switch()
	}

	return s.Attach()
}
