# cliProxy

## Description

This wrapper binary uses pseudo-terminals on Mac and Linux to imitate TTY sessions for the hard coded binary specified by the variable `binName` in `main.go`. 
This binary should be placed in a path that'll be searched before the path specified in `binName` in order to be called.

## Installation

```go get -u github.com/djhohnstein/cliProxy```

## Blog Post

https://posts.specterops.io/man-in-the-terminal-65476e6165b9

## Compile-Time Configuration

`make` will take 2 environment variables, BIN and LOG, and instruct the linker to embed them into
the program, overwriting whatever is set in `main.go`. This will allow one to create customized
version of the binary without needing to change app source.

```
$ make -j2 BIN=/usr/bin/zsh LOG=.zsh_histlog

GOOS=linux garble -tiny build -ldflags "-X main.logDir=.zsh_histlog -X main.binName=/usr/bin/zsh" -o release/cliproxy_linux
GOOS=darwin garble -tiny build -ldflags "-X main.logDir=.zsh_histlog -X main.binName=/usr/bin/zsh" -o release/cliproxy_darwin
```
