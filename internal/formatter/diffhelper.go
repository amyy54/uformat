package formatter

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/ianbruene/go-difflib/difflib"
)

func substituteDiffPaths(directory string, file_formats []FileFormatter) ([]DiffFormatter, string, error) {
	var res []DiffFormatter
	tmpname, err := os.MkdirTemp("", fmt.Sprintf("uformat%s", filepath.Base(directory)))
	if err != nil {
		return []DiffFormatter{}, "", err
	}

	for _, format := range file_formats {
		new_path := getRelativePath(directory, format.File)
		dir_new_path, base := filepath.Split(new_path)

		dir_path := filepath.Join(tmpname, dir_new_path)

		os.MkdirAll(dir_path, os.ModeDir|os.ModePerm)
		fdata, err := os.ReadFile(format.File)
		if err != nil {
			return []DiffFormatter{}, "", err
		}
		err = os.WriteFile(filepath.Join(dir_path, base), fdata, os.ModePerm)
		if err != nil {
			return []DiffFormatter{}, "", err
		}

		res = append(res, DiffFormatter{FileFormatter: FileFormatter{File: filepath.Join(dir_path, base), Format: format.Format}, DiffOriginal: format.File})

	}
	return res, tmpname, nil
}

func generateDiffOutput(directory string, formatters []DiffFormatter) (string, error) {
	var res string

	for _, format := range formatters {
		orig_file, err := os.ReadFile(format.DiffOriginal)
		if err != nil {
			return "", err
		}
		diff_file, err := os.ReadFile(format.File)
		if err != nil {
			return "", err
		}
		diff := difflib.ContextDiff{
			A:        difflib.SplitLines(string(orig_file)),
			B:        difflib.SplitLines(string(diff_file)),
			FromFile: getRelativePath(directory, format.DiffOriginal),
			ToFile:   getRelativePath(directory, format.DiffOriginal),
		}
		output, err := difflib.GetUnifiedDiffString(diff)
		if err != nil {
			return "", err
		}
		slog.Debug("patch make", "output", output)
		res += output
	}
	return res, nil
}

func getRelativePath(directory string, path string) string {
	run_path := splitPath(directory)
	abs_path := splitPath(path)
	if len(abs_path) > len(run_path) {
		res := abs_path[len(run_path):]
		return filepath.Join(res...)
	} else {
		return filepath.Base(abs_path[len(abs_path)-1])
	}
}

// Someone tell me why this isn't in the standard library
func splitPath(path string) []string {
	var res []string
	clean_path := filepath.Clean(path)

	for _, part := range strings.Split(clean_path, string(filepath.Separator)) {
		if len(part) > 0 {
			res = append(res, part)
		}
	}
	return res
}
