package util

import (
	"os"
	"regexp"
	"strings"
)

// Sanitize converts branch name to directory-safe string
func Sanitize(branch string) string {
	s := strings.TrimSpace(branch)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, string(os.PathSeparator), "-")
	s = strings.ReplaceAll(s, ":", "-")
	s = strings.ReplaceAll(s, "@", "-")
	re := regexp.MustCompile(`[^A-Za-z0-9._-]+`)
	s = re.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		s = "worktree"
	}
	return s
}

// ShellJoin joins command arguments with proper quoting
func ShellJoin(args []string) string {
	var b []string
	for _, x := range args {
		if strings.ContainsAny(x, " \t\n\"'\\$") {
			b = append(b, Quote(x))
		} else {
			b = append(b, x)
		}
	}
	return strings.Join(b, " ")
}

// Quote quotes a string for shell
func Quote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}
