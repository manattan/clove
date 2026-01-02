package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

// Git executes a git command and returns stdout/stderr
func Git(repoRoot string, args ...string) (string, error) {
	var cmd *exec.Cmd
	if repoRoot != "" {
		a := append([]string{"-C", repoRoot}, args...)
		cmd = exec.Command("git", a...)
	} else {
		cmd = exec.Command("git", args...)
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%w: %s", err, out.String())
	}
	return out.String(), nil
}

// GitOk checks if git command succeeds
func GitOk(repoRoot string, args ...string) bool {
	var cmd *exec.Cmd
	if repoRoot != "" {
		a := append([]string{"-C", repoRoot}, args...)
		cmd = exec.Command("git", a...)
	} else {
		cmd = exec.Command("git", args...)
	}
	return cmd.Run() == nil
}

// Run executes a command with stdout/stderr attached
func Run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout // エラー出力も標準出力に出す
	return cmd.Run()
}
