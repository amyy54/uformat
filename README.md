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
  -show
        List the files formatted using their relative path.
  -show-abs
        List the files formatted using their absolute path. Overrides -show.
  -v    Print logs tagged "Info" or higher.
  -version
        Print the version and exit.
  -vv
        Print logs tagged "Debug" or higher.
```

## uformat.json

The format of the file is fairly self explanatory, and can be seen in this
project. The name is purely decorative and helps for organization. The `ignore`
section requires some weird formatting (as of now), as it uses absolute path
names, as opposed to relative.
