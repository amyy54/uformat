package formatter

import (
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/gobwas/glob"

	"github.com/amyy54/uformat/internal/configloader"
)

func matchFiles(directory string, formatters []configloader.Formatter, ignores []glob.Glob, use_git bool) ([]FileFormatter, error) {
	var res []FileFormatter

	if file_info, err := os.Stat(directory); err == nil {
		if file_info.IsDir() {
			_, git_err := findRepository(directory)
			if git_err != nil {
				slog.Warn("not a git directory, not checking ignores", "error", git_err)
			}
			gave_git_warning := false
			err := filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
				slog.Debug("testing path", "path", path, "name", d.Name())
				is_ignored, ignore_err := isIgnored(directory, path, d.Name(), ignores, use_git)
				if ignore_err == nil || gave_git_warning {
					if is_ignored {
						if d.IsDir() {
							slog.Info("ignoring directory", "directory", path)
							return filepath.SkipDir
						} else {
							slog.Info("ignoring file", "file", path)
						}
					} else {
						if !d.IsDir() {
							if found, formatter := matchFile(path, formatters); found {
								res = append(res, formatter)
							} else {
								slog.Info("did not find a formatter for the path specified", "path", path)
							}
						}
					}
				} else {
					slog.Warn("file failed to be checked for ignorance. this is likely because of git. execution will continue regardless.", "error", ignore_err)
					gave_git_warning = true
				}
				return nil
			})
			if err != nil {
				return []FileFormatter{}, err
			}
			return res, nil
		} else {
			return []FileFormatter{}, fmt.Errorf("path %v is not a directory", directory)
		}
	} else {
		return []FileFormatter{}, err
	}
}

func matchFile(path string, formatters []configloader.Formatter) (bool, FileFormatter) {
	for _, formatter := range formatters {
		g := glob.MustCompile(formatter.Glob)
		if g.Match(path) {
			return true, FileFormatter{File: path, Format: formatter}
		}
	}
	return false, FileFormatter{}
}
