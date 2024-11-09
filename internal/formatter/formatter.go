package formatter

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/amyy54/uformat/internal/configloader"
)

func MatchSingle(config configloader.Config, path string, format_module string) (FileFormatter, error) {
	if _, err := os.Stat(path); err == nil {
		config_formats := config.ToFormatList()
		if len(format_module) > 0 {
			config_formats, err = config.FilterFormatList(format_module)
			if err != nil {
				return FileFormatter{}, err
			}
		}
		matched, formatter := matchFile(path, config_formats)

		if matched {
			return formatter, nil
		} else {
			return FileFormatter{}, fmt.Errorf("did not find a formatter for the path specified %s", path)
		}
	} else {
		return FileFormatter{}, fmt.Errorf("file does not exist, exiting")
	}
}

func Format(config configloader.Config, directory string, options FormatOptions) (int, string, string, error) {
	var err error
	var output string
	counter := 0
	var tempdiffdir string
	var diff_need_formatting []DiffFormatter
	var paths []string

	need_formatting := options.FileFormatters

	if len(options.FileFormatters) == 0 {
		config_formats := config.ToFormatList()
		if len(options.FormatModule) > 0 {
			config_formats, err = config.FilterFormatList(options.FormatModule)
			if err != nil {
				return 0, "", "", err
			}
		}
		need_formatting, err = matchFiles(directory, config_formats, config.IgnoreToGlob(), options.UseGit, len(options.FormatModule) > 0)

		if err != nil {
			return 0, "", "", err
		}
	}

	if options.Diff {
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
	if options.Diff {
		logdir = tempdiffdir
	}
	for num, need_to_format := range need_formatting {
		counter++
		slog.Info(fmt.Sprintf("%d) running formatter for file %s", num+1, getRelativePath(logdir, need_to_format.File)), "format", need_to_format.ToLogString())
		if p_output, err := process_execution(need_to_format); err == nil {
			output += p_output
			if options.AbsolutePath {
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

	if options.Diff {
		output, err = generateDiffOutput(directory, diff_need_formatting, options.UseGit)
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
	for _, arg := range formatter.Format.Args {
		if strings.Contains(arg, "<file>") {
			modified_args = append(modified_args, strings.ReplaceAll(arg, "<file>", formatter.File))
		} else if strings.Contains(arg, "<fileName>") {
			modified_args = append(modified_args, strings.ReplaceAll(arg, "<fileName>", strings.TrimSuffix(formatter.File, filepath.Ext(formatter.File))))
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
