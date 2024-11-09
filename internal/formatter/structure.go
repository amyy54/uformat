package formatter

import (
	"fmt"
	"strings"

	"github.com/amyy54/uformat/internal/configloader"
)

type FileFormatter struct {
	File   string
	Format configloader.Formatter
}

func (f *FileFormatter) ToLogString() string {
	return fmt.Sprintf("%s (%s) matching %s", f.Format.Command, strings.Join(f.Format.Args, " "), f.Format.Glob)
}

type DiffFormatter struct {
	FileFormatter
	DiffOriginal string
}

type FormatOptions struct {
	UseGit         bool
	Diff           bool
	AbsolutePath   bool
	FileFormatters []FileFormatter
	FormatModule   string
}
