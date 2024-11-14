package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/amyy54/uformat/internal/configloader"
	"github.com/amyy54/uformat/internal/formatter"
)

var (
	Version     string
	VersionLong string
)

func main() {
	var config_location string
	var target_dir string
	var show_formats bool
	var ignore_git bool
	var diff_mode bool
	var show_files bool
	var show_abs bool
	var single_file string
	var format_module string
	var output_file string
	var stdin bool

	var v bool
	var vv bool

	var version bool

	log.SetFlags(log.Lshortfile)

	flag.StringVar(&config_location, "config", "./.uformat.json", "Configuration file to load.")
	flag.StringVar(&target_dir, "directory", ".", "Target directory to format.")
	flag.BoolVar(&show_formats, "list", false, "List available formats in the loaded configuration file.")
	flag.BoolVar(&ignore_git, "ignore-git", false, "Ignore git and all related functions, such as checking gitignore.")
	flag.BoolVar(&diff_mode, "diff", false, "Instead of formatting, print the difference. This acts as a universal dry-run.")
	flag.BoolVar(&show_files, "show", false, "List the files formatted using their relative path.")
	flag.BoolVar(&show_abs, "show-abs", false, "List the files formatted using their absolute path. Overrides -show.")
	flag.StringVar(&single_file, "file", "", "Instead of formatting a directory, format the specified file.")
	flag.StringVar(&format_module, "module", "", "Format using only the specified module.")
	flag.StringVar(&output_file, "output", "", "When using -file, specify the output for the formatted file. - is stdout.")
	flag.BoolVar(&stdin, "stdin", false, "Read from standard input to format file. -file or -module required.")

	flag.BoolVar(&v, "v", false, "Print logs tagged \"Info\" or higher.")
	flag.BoolVar(&vv, "vv", false, "Print logs tagged \"Debug\" or higher.")

	flag.BoolVar(&version, "version", false, "Print the version and exit.")
	flag.Parse()

	if version {
		if v || vv {
			fmt.Printf("uformat: %s\n", VersionLong)
		} else {
			fmt.Printf("uformat: %s\n", Version)
		}
		os.Exit(0)
	}

	if vv {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	} else if v {
		slog.SetLogLoggerLevel(slog.LevelInfo)
	} else {
		slog.SetLogLoggerLevel(slog.LevelWarn)
	}

	resolve_conf_location, err := filepath.Abs(config_location)
	if err != nil {
		log.Fatal("could not resolve config_location")
	}

	if config_location == "./.uformat.json" {
		if _, err := os.Stat(resolve_conf_location); err != nil {
			slog.Info("file not found in current directory, trying home folder", "config_location", config_location, "resolve_conf_location", resolve_conf_location)
			homedir, err := os.UserHomeDir()
			if err != nil {
				log.Fatal("could not resolve home folder")
			}
			path := filepath.Join(homedir, ".uformat.json")
			if _, err := os.Stat(path); err == nil {
				resolve_conf_location = path
			} else {
				log.Fatal("no format file found in home folder or current folder, exiting")
			}
		}
	}

	target_dir, err = filepath.Abs(target_dir)
	if err != nil {
		log.Fatal("could not resolve target_dir")
	}

	resolve_output_file := output_file
	if len(output_file) > 0 && output_file != "-" {
		resolve_output_file, err = filepath.Abs(resolve_output_file)
		if err != nil {
			log.Fatal("could not resolve output_file")
		}
	}

	resolve_single_file := single_file
	if len(single_file) > 0 {
		resolve_single_file, err = filepath.Abs(resolve_single_file)
		if err != nil {
			log.Fatal("could not resolve single_file")
		}
	}

	slog.Debug("flags parsed", "config_location", config_location, "resolved_config_location", resolve_conf_location, "target_dir", target_dir, "show_formats", show_formats, "ignore_git", ignore_git, "diff_mode", diff_mode, "show_files", show_files, "show_abs", show_abs, "single_file", single_file, "format_module", format_module, "resolve_output_file", resolve_output_file)

	config, err := configloader.LoadConfig(resolve_conf_location)
	if err != nil {
		log.Fatal(err)
	}

	slog.Debug("config loaded", "config", config)

	if show_formats {
		slog.Info("config location", "resolve_conf_location", resolve_conf_location)
		for name, formats := range config.Formats {
			fmt.Printf("%s, matching \"%s\" with command: %s (%s)\n", name, formats.Glob, formats.Command, strings.Join(formats.Args, " "))
		}
	} else if stdin {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			log.Fatal("no standard input provided")
		}
		inputBytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		output, err := formatter.FormatText(config, string(inputBytes), resolve_single_file, formatter.FormatOptions{
			UseGit:         !ignore_git,
			Diff:           diff_mode,
			AbsolutePath:   show_abs,
			FileFormatters: []formatter.FileFormatter{},
			FormatModule:   format_module,
			OutputFile:     resolve_output_file,
		})
		if err != nil {
			log.Fatal(err)
		} else {
			if diff_mode || resolve_output_file == "-" {
				fmt.Print(output)
			} else {
				fmt.Println("✨ Formatted 1 files")
			}
		}
	} else {
		var parsed_files []formatter.FileFormatter
		if len(resolve_single_file) > 0 {
			single_formatter, err := formatter.MatchSingle(config, resolve_single_file, format_module)
			if err != nil {
				log.Fatal(err)
			} else {
				parsed_files = append(parsed_files, single_formatter)
			}
		} else {
			resolve_output_file = ""
		}
		count, output, paths, err := formatter.Format(config, target_dir, formatter.FormatOptions{UseGit: !ignore_git, Diff: (diff_mode || len(resolve_output_file) > 0), AbsolutePath: show_abs, FileFormatters: parsed_files, FormatModule: format_module, OutputFile: resolve_output_file})

		if err != nil {
			log.Fatal(err)
		} else {
			if diff_mode || resolve_output_file == "-" {
				fmt.Print(output)
			} else if show_files || show_abs {
				fmt.Println(paths)
			} else {
				fmt.Printf("✨ Formatted %d files\n", count)
			}
		}
	}

	os.Exit(0)
}
