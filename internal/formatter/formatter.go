package formatter

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"

	"github.com/amyy54/uformat/internal/configloader"
)

func Format(config configloader.Config, directory string, use_git bool, use_diff bool) (int, string, error) {
	var output string
	counter := 0
	var tempdiffdir string
	var diff_need_formatting []DiffFormatter

	need_formatting, err := matchFiles(directory, config.ToFormatList(), config.IgnoreToGlob(), use_git)

	if err != nil {
		return 0, "", err
	}

	if use_diff {
		diff_need_formatting, tempdiffdir, err = substituteDiffPaths(directory, need_formatting)
		if err != nil {
			return 0, "", err
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
			if len(output) > 0 {
				slog.Debug("output", "output", output)
			}
		} else {
			return 0, "", err
		}
	}

	if use_diff {
		output, err = generateDiffOutput(directory, diff_need_formatting, use_git)
		if err != nil {
			return 0, "", err
		}
		os.RemoveAll(tempdiffdir)
	}

	return counter, output, nil
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
