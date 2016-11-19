package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Roman2K/tmuxstatus/cpu"
	"github.com/Roman2K/tmuxstatus/ports"
)

var commands commandsMap

type commandsMap map[string]func() error

func init() {
	commands = commandsMap{
		"cpu":   cmdCPU,
		"ports": cmdPorts,
	}
}

func main() {
	switch len(os.Args) {
	case 1:
		usage(os.Stdout)
		os.Exit(0)
	case 2:
		cmd, ok := commands[os.Args[1]]
		if !ok {
			usageErr()
		}
		if err := cmd(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	default:
		usageErr()
	}
}

func usage(w io.Writer) {
	exe := filepath.Base(os.Args[0])
	cmds := []string{}
	for name, _ := range commands {
		cmds = append(cmds, name)
	}
	sort.Strings(cmds)
	fmt.Fprintf(w, "usage: %s %s\n", exe, strings.Join(cmds, "|"))
}

func usageErr() {
	usage(os.Stderr)
	os.Exit(2)
}

func cmdCPU() (err error) {
	top, err := cpu.Top(3)
	if err != nil {
		return
	}

	descs := make([]string, len(top))
	for i, line := range top {
		const maxCmdLen = 7
		cmd := truncate(shortCommand(line.Command), maxCmdLen)
		descs[i] = fmt.Sprintf("%.1f %s", line.Pct, cmd)
	}

	fmt.Println(strings.Join(descs, ", "))
	return
}

func cmdPorts() (err error) {
	ports, err := ports.List([]string{"8000-8999", "3000-3999"})
	if err != nil {
		return
	}

	sort.Ints(ports)

	descs := make([]string, len(ports))
	for i, port := range ports {
		descs[i] = fmt.Sprintf("%d", port)
	}

	fmt.Println(strings.Join(descs, ", "))
	return
}

func shortCommand(cmd string) string {
	if len(cmd) > 0 && cmd[0] != '/' {
		return cmd
	}
	parts := strings.SplitN(cmd, " ", 2)
	exe := filepath.Base(parts[0])
	if len(parts) < 2 {
		return exe
	}
	return exe + " " + parts[1]
}

func truncate(str string, max int) string {
	if len(str) <= max {
		return str
	}
	return str[:max]
}
