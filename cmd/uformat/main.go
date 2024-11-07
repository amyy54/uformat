package main

import (
	"flag"
	"fmt"
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

	slog.Debug("flags parsed", "config_location", config_location, "resolved_config_location", resolve_conf_location, "target_dir", target_dir, "show_formats", show_formats, "ignore_git", ignore_git)

	config, err := configloader.LoadConfig(resolve_conf_location)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("config loaded", "config", config)

	if show_formats {
		for name, formats := range config.Formats {
			fmt.Printf("%s, matching \"%s\" with command: %s (%s)\n", name, formats.Glob, formats.Command, strings.Join(formats.Args, " "))
		}
		os.Exit(0)
	} else {
		count, output, paths, err := formatter.Format(config, target_dir, !ignore_git, diff_mode, show_abs)

		if err != nil {
			log.Fatal(err)
		} else {
			if diff_mode {
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
