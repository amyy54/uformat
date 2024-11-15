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

func FormatText(config configloader.Config, input string, path string, options FormatOptions) (bool, string, error) {
	file, err := os.CreateTemp("", "uformatstdin")
	if err != nil {
		return false, "", err
	}
	_, err = file.Write([]byte(input))
	if err != nil {
		return false, "", err
	}

	tmpPath := file.Name()

	file.Close()

	var diffPath string

	if options.Diff {
		file, err = os.CreateTemp("", "uformatstdindiff")
		if err != nil {
			return false, "", err
		}
		_, err = file.Write([]byte(input))
		if err != nil {
			return false, "", err
		}

		diffPath = file.Name()

		file.Close()
	}

	var formatter FileFormatter
	if len(path) > 0 {
		var success bool
		success, formatter = matchFile(path, config.ToFormatList())
		if !success {
			return false, "", fmt.Errorf("did not find a formatter for the path specified %s", path)
		}
		formatter.File = tmpPath
	} else if len(options.FormatModule) > 0 {
		config_formats, err := config.FilterFormatList(options.FormatModule)
		if err != nil {
			return false, "", err
		}
		formatter = FileFormatter{File: tmpPath, Format: config_formats[0]}
	} else {
		return false, "", fmt.Errorf("Either -file or -module needs to be passed to read standard input")
	}
	slog.Info("running formatter on standard input", "format", formatter.ToLogString())
	_, err = process_execution(formatter)
	if err != nil {
		return false, "", err
	}

	var output string
	show_output := false

	if len(options.OutputFile) > 0 {
		file_output, err := os.ReadFile(tmpPath)
		if err != nil {
			return false, "", err
		}
		if options.OutputFile == "-" {
			output = string(file_output)
			show_output = true
		} else {
			err = os.WriteFile(options.OutputFile, file_output, os.ModePerm)
			if err != nil {
				return false, "", err
			}
		}
	}
	if options.Diff && options.OutputFile != "-" {
		var diffformatters []DiffFormatter

		diffformatters = append(diffformatters, DiffFormatter{FileFormatter: formatter, DiffOriginal: diffPath})

		output, err = generateDiffOutput(filepath.Dir(tmpPath), diffformatters, false)
		if err != nil {
			return false, "", err
		}
		show_output = true
	}

	os.Remove(tmpPath)
	if options.Diff {
		os.Remove(diffPath)
	}

	return show_output, output, nil
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
		if len(options.OutputFile) > 0 {
			if len(diff_need_formatting) > 0 {
				file_output, err := os.ReadFile(diff_need_formatting[0].File)
				if err != nil {
					return 0, "", "", err
				}

				if options.OutputFile == "-" {
					output = string(file_output)
				} else {
					err = os.WriteFile(options.OutputFile, file_output, os.ModePerm)
					if err != nil {
						return 0, "", "", err
					}
				}
			}
		}
		if len(options.OutputFile) == 0 || options.OutputFile != "-" {
			output, err = generateDiffOutput(directory, diff_need_formatting, options.UseGit)
			if err != nil {
				return 0, "", "", err
			}
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
