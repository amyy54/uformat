package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/amyy54/uformat/internal/configloader"
	"github.com/amyy54/uformat/internal/formatter"
)

func main() {
	var config_location string
	var target_dir string
	var show_formats bool
	var ignore_git bool

	var v bool
	var vv bool

	log.SetFlags(log.Lshortfile)

	flag.StringVar(&config_location, "config", "./.uformat.json", "Formatter configuration file.")
	flag.StringVar(&target_dir, "directory", ".", "Target directory.")
	flag.BoolVar(&show_formats, "list", false, "List available formats.")
	flag.BoolVar(&ignore_git, "ignore-git", false, "Ignores git and all related functions (checking gitignore, etc).")

	flag.BoolVar(&v, "v", false, "Shows logs with \"Info\" or higher.")
	flag.BoolVar(&vv, "vv", false, "Shows logs with \"Debug\" or higher.")
	flag.Parse()

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
			fmt.Printf("%v, matching \"%v\" | Command: %v\n", name, formats.Glob, formats.Command)
		}
		os.Exit(0)
	} else {
		count, _, err := formatter.Format(config, target_dir, !ignore_git)

		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Printf("✨ Formatted %d files\n", count)
		}
	}

	os.Exit(0)
}
