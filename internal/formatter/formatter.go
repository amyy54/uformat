package formatter

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"

	"github.com/amyy54/uformat/internal/configloader"
)

func MatchSingle(config configloader.Config, path string) (FileFormatter, error) {
	if _, err := os.Stat(path); err == nil {
		matched, formatter := matchFile(path, config.ToFormatList())

		if matched {
			return formatter, nil
		} else {
			return FileFormatter{}, fmt.Errorf("did not find a formatter for the path specified %s", path)
		}
	} else {
		return FileFormatter{}, fmt.Errorf("file does not exist, exiting")
	}
}

func Format(config configloader.Config, directory string, use_git bool, use_diff bool, abs_path bool, file_formatters []FileFormatter) (int, string, string, error) {
	var err error
	var output string
	counter := 0
	var tempdiffdir string
	var diff_need_formatting []DiffFormatter
	var paths []string

	need_formatting := file_formatters

	if len(file_formatters) == 0 {
		need_formatting, err = matchFiles(directory, config.ToFormatList(), config.IgnoreToGlob(), use_git)

		if err != nil {
			return 0, "", "", err
		}
	}

	if use_diff {
		diff_need_formatting, tempdiffdir, err = substituteDiffPaths(directory, need_formatting)
		if err != nil {
			return 0, "", "", err
		}

		need_formatting = []FileFormatter{}
		for _, format := range diff_need_formatting {
			need_formatting = append(need_formatting, format.FileFormatter)
		}
	}

	slog.Debug("----------") // Starting process logs
	logdir := directory
	if use_diff {
		logdir = tempdiffdir
	}
	for num, need_to_format := range need_formatting {
		counter++
		slog.Info(fmt.Sprintf("%d) running formatter for file %s", num+1, getRelativePath(logdir, need_to_format.File)), "format", need_to_format.ToLogString())
		if p_output, err := process_execution(need_to_format); err == nil {
			output += p_output
			if abs_path {
				paths = append(paths, need_to_format.File)
			} else {
				paths = append(paths, getRelativePath(logdir, need_to_format.File))
			}
			if len(output) > 0 {
				slog.Debug("output", "output", output)
			}
		} else {
			return 0, "", "", err
		}
	}

	if use_diff {
		output, err = generateDiffOutput(directory, diff_need_formatting, use_git)
		if err != nil {
			return 0, "", "", err
		}
		os.RemoveAll(tempdiffdir)
	}

	return counter, output, strings.Join(paths, "\n"), nil
}

func process_execution(formatter FileFormatter) (string, error) {
	_, err := exec.LookPath(formatter.Format.Command)

	if err != nil {
		return "", fmt.Errorf("could not find \"%v\" in the path. cannot execute formatter", formatter.Format.Command)
	}

	var modified_args []string
	// TODO There is absolutely a better way of doing this.
	for _, arg := range formatter.Format.Args {
		if arg == "<file>" {
			modified_args = append(modified_args, formatter.File)
		} else {
			modified_args = append(modified_args, arg)
		}
	}

	slog.Debug("spawning process for formatter", "command", formatter.Format.Command, "args", modified_args)

	format_exec := exec.Command(formatter.Format.Command, modified_args...)

	output, err := format_exec.CombinedOutput()

	if err != nil {
		return string(output), fmt.Errorf("error running formatter: %v", err)
	} else {
		return string(output), nil
	}
}
