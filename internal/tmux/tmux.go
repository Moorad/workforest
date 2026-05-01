package tmux

import (
	"os"
	"strconv"
	"strings"

	"github.com/Moorad/workforest/internal/exec"
)

type Session struct {
	name string
	dir  string
}

func Health() (bool, string) {
	out, err := exec.NewCommand("tmux", "-V").Output()
	if err != nil {
		return false, ""
	}

	return true, strings.Trim(strings.Split(out, " ")[1], "\n")
}

func NewSession(name string, dir string) (*Session, error) {
	session := Session{
		name,
		dir,
	}

	cmd := exec.NewCommand("tmux", "new-session", "-d", "-s", name)
	cmd.Cmd.Dir = session.dir
	err := cmd.Run()

	return &session, err
}

func (s *Session) RenameWindowMain(name string) error {
	cmd := exec.NewCommand("tmux", "rename-window", "-t", s.name, name)
	cmd.Cmd.Dir = s.dir
	return cmd.Run()
}

func (s *Session) NewWindow(name string) error {
	cmd := exec.NewCommand("tmux", "new-window", "-t", s.name, "-n", name)
	cmd.Cmd.Dir = s.dir
	return cmd.Run()
}

func (s *Session) SendKeys(windowName string, keys ...string) error {
	cmd := []string{"tmux", "send-keys", "-t", s.name + ":" + windowName}
	cmd = append(cmd, keys...)

	command := exec.NewCommand(cmd...)
	command.Cmd.Dir = s.dir
	return command.Run()
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

func DoesSessionExist(name string) bool {
	cmd := exec.NewCommand("tmux", "has-session", "-t", name)

	return cmd.Run() == nil
}

func DirectSwitchOrAttach(name string) error {
	termProgram, isFound := os.LookupEnv("TERM_PROGRAM")

	if isFound && termProgram == "tmux" {
		return exec.SysCall("tmux", "switch", "-t="+name)
	}

	return exec.SysCall("tmux", "attach", "-t="+name)
}

type SessionInfo struct {
	Name        string
	WindowCount int
	IsAttached  bool
}

func ListSessions() ([]SessionInfo, error) {
	cmd := exec.NewCommand("tmux", "list-sessions", "-F", "#{session_name}:#{session_windows}:#{session_attached}")

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.Trim(out, "\n"), "\n")

	sessions := []SessionInfo{}

	for _, line := range lines {
		info := strings.Split(line, ":")

		windowCount, err := strconv.Atoi(info[1])
		if err != nil {
			return nil, err
		}

		sessions = append(sessions, SessionInfo{Name: info[0], WindowCount: windowCount, IsAttached: info[2] == "1"})
	}

	return sessions, nil
}
