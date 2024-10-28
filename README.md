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
    	Formatter configuration file. (default "./.uformat.json")
  -directory string
    	Target directory. (default ".")
  -ignore-git
    	Ignores git and all related functions (checking gitignore, etc).
  -list
    	List available formats.
  -v	Shows logs with "Info" or higher.
  -vv
    	Shows logs with "Debug" or higher.
```

## uformat.json

The format of the file is fairly self explanatory, and can be seen in this
project. The name is purely decorative and helps for organization. The `ignore`
section requires some weird formatting (as of now), as it uses absolute path
names, as opposed to relative.
