module github.com/amyy54/uformat

go 1.23.2

replace github.com/amyy54/uformat/internal/configloader => ./internal/configloader

replace github.com/amyy54/uformat/internal/formatter => ./internal/formatter

require github.com/gobwas/glob v0.2.3

require github.com/ianbruene/go-difflib v1.2.0
