package formatter

import (
	"fmt"
	"log/slog"
	"os/exec"

	"github.com/amyy54/uformat/internal/configloader"
)

func Format(config configloader.Config, directory string, use_git bool) (int, string, error) {
	var output string
	counter := 0

	need_formatting, err := matchFiles(directory, config.ToFormatList(), config.IgnoreToGlob(), use_git)

	if err != nil {
		return 0, "", err
	}
	slog.Debug("----------") // Starting process logs
	for num, need_to_format := range need_formatting {
		counter++
		slog.Info(fmt.Sprintf("%d) running formatter for %v", num+1, need_to_format.File), "format", need_to_format)
		if p_output, err := process_execution(need_to_format); err == nil {
			output += p_output
			slog.Debug(output)
		} else {
			return 0, "", err
		}
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
