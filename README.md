# uformat

(U)niversal (format)ter. For projects that utilize multiple formatters, and need
consistency.

## About

This project was made for two reasons.

1. I wanted an excuse to learn Go. Here we are!
1. Because memorizing formats across projects gets inconvenient.

I felt like creating a simple tool to keep formats universal with just one
program is handy, and this allows for multiple format programs to exist in a
project and run with just one simple tool.

## Install

Pre-compiled binaries can be found on the Releases page.

Support:

- macOS (universal binary Intel/M1).
- Linux (amd64/arm64). Pre-packaged for Debian/RHEL.
- Windows (amd64/arm64).

The macOS binary can be installed using brew:
`brew install amyy54/taps/uformat`.

## Usage

```
Usage of uformat:
  -config string
        Configuration file to load. (default "./.uformat.json")
  -diff
        Instead of formatting, print the difference. This acts as a universal dry-run.
  -directory string
        Target directory to format. (default ".")
  -file string
        Instead of formatting a directory, format the specified file.
  -ignore-git
        Ignore git and all related functions, such as checking gitignore.
  -list
        List available formats in the loaded configuration file.
  -module string
        Format using only the specified module.
  -output string
        When using -file, specify the output for the formatted file. - is stdout.
  -show
        List the files formatted using their relative path.
  -show-abs
        List the files formatted using their absolute path. Overrides -show.
  -stdin
        Read from standard input to format file. -file or -module required.
  -v    Print logs tagged "Info" or higher.
  -version
        Print the version and exit.
  -vv
        Print logs tagged "Debug" or higher.
```

## Configuration

The JSON schema for the configuration file (.uformat.json) can be found in
[the dist folder](dist/.uformat.schema.json). Light documentation is available
through the schema and its descriptions of values.
