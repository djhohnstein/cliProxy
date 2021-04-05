# cliProxy

## Description

This wrapper binary uses pseudo-terminals on Mac and Linux to imitate TTY sessions for the hard coded binary specified by the variable `binName` in `main.go`. 
This binary should be placed in a path that'll be searched before the path specified in `binName` in order to be called.

## Installation

```go get -u github.com/djhohnstein/cliProxy```

## Blog Post

https://posts.specterops.io/man-in-the-terminal-65476e6165b9
