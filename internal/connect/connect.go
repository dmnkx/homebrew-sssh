package connect

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"sssh/internal/debuglog"
)

func SSHPath() (string, error) {
	p, err := exec.LookPath("ssh")
	if err != nil {
		return "", fmt.Errorf("ssh not found in PATH")
	}
	return p, nil
}

func PrintCmd(alias string) (string, error) {
	p, err := SSHPath()
	if err != nil {
		return "", err
	}
	return p + " " + alias, nil
}

func Exec(alias string) error {
	p, err := SSHPath()
	// #region agent log
	debuglog.Log("H1", "connect.go:Exec", "resolved ssh", map[string]any{"path": p, "alias": alias, "err": errString(err)})
	// #endregion
	if err != nil {
		return err
	}
	argv := []string{"ssh", alias}
	env := os.Environ()
	// #region agent log
	debuglog.Log("H2", "connect.go:Exec", "before syscall.Exec", map[string]any{"argv": argv})
	// #endregion
	return syscall.Exec(p, argv, env)
}

func errString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
