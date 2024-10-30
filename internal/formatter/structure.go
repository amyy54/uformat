package formatter

import "github.com/amyy54/uformat/internal/configloader"

type FileFormatter struct {
	File   string
	Format configloader.Formatter
}

type DiffFormatter struct {
	FileFormatter
	DiffOriginal string
}
