package formatter

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/gobwas/glob"
)

func isIgnored(directory string, path string, name string, ignores []glob.Glob, use_git bool) (bool, error) {
	if isDefaultIgnored(name) || isConfigIgnored(path, ignores) {
		return true, nil
	}
	if use_git {
		is_ignore, err := isGitIgnored(directory, path)
		if err != nil {
			return false, err
		}
		return is_ignore, nil
	} else {
		return false, nil
	}
}

func isDefaultIgnored(name string) bool {
	defaults := []glob.Glob{
		glob.MustCompile(".git"),
		glob.MustCompile(".gitignore"),
		glob.MustCompile(".uformat.json"),
		glob.MustCompile("*LICENSE*"),
	}
	for _, defaults_glob := range defaults {
		if defaults_glob.Match(name) {
			return true
		}
	}
	return false
}

func isConfigIgnored(path string, ignores []glob.Glob) bool {
	for _, ignore := range ignores {
		if ignore.Match(path) {
			return true
		}
	}
	return false
}

func isGitIgnored(directory string, path string) (bool, error) {
	if git_binary, err := exec.LookPath("git"); err == nil {
		err := os.Chdir(directory)
		if err != nil {
			return false, err
		}

		cmd := exec.Command(git_binary, "check-ignore", path)

		_, err = cmd.Output()

		if err != nil {
			switch t := err.(type) {
			case *exec.ExitError:
				if t.ExitCode() == 1 {
					return false, nil
				} else {
					return false, fmt.Errorf("git exit code something other than 1. %v", t.ExitCode())
				}
			default:
				return false, err
			}
		} else {
			return true, nil
		}
	} else {
		return false, err
	}
}

// Returns the path for the top level repository, if it exists.
func findRepository(directory string) (string, error) {
	if git_binary, err := exec.LookPath("git"); err == nil {
		err := os.Chdir(directory)
		if err != nil {
			return "", err
		}
		cmd := exec.Command(git_binary, "rev-parse", "--show-toplevel")

		output, err := cmd.Output()

		if err == nil {
			return string(output), nil
		} else {
			return "", err
		}
	} else {
		return "", err
	}
}
