package configloader

import (
	"fmt"
	"strings"

	"github.com/gobwas/glob"
)

type Config struct {
	Version int                  `json:"version"`
	Formats map[string]Formatter `json:"formats"`
	Ignore  []string             `json:"ignore"`
}

type Formatter struct {
	Glob    string   `json:"glob"`
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

func (c *Config) ToFormatList() []Formatter {
	var res []Formatter
	for _, formats := range c.Formats {
		res = append(res, formats)
	}
	return res
}

func (c *Config) FilterFormatList(filter string) ([]Formatter, error) {
	var res []Formatter
	for name, formats := range c.Formats {
		if strings.EqualFold(name, filter) {
			res = append(res, formats)
		}
	}
	if len(res) == 0 {
		return []Formatter{}, fmt.Errorf("formatter does not exist in loaded configuration file")
	}
	return res, nil
}

func (c *Config) IgnoreToGlob() []glob.Glob {
	var res []glob.Glob
	for _, ignore := range c.Ignore {
		res = append(res, glob.MustCompile(ignore))
	}
	return res
}
