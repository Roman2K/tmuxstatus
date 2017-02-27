package ports

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"syscall"
)

var listenRe *regexp.Regexp

func init() {
	listenRe = regexp.MustCompile(`.*:(\d+) \(LISTEN\)`)
}

func List(ranges []string) (ports []int, err error) {
	args := []string{"-l", "-P", "-n", "-s", "TCP:LISTEN"}

	for _, r := range ranges {
		args = append(args, "-i", "TCP:"+r)
	}
	if len(ranges) == 0 {
		args = append(args, "-i", "TCP")
	}

	cmd := exec.Command("lsof", args...)
	out, err := cmd.Output()
	if err != nil {
		err = makeLsofErr(err)
		if err != nil {
			return
		}
	}

	seen := map[int]bool{}

	split := bufio.NewScanner(bytes.NewReader(out))
	split.Scan() // skip headers
	for split.Scan() {
		m := listenRe.FindStringSubmatch(split.Text())
		if m == nil {
			continue
		}
		port, err := strconv.Atoi(m[1])
		if err != nil {
			continue
		}
		if _, ok := seen[port]; ok {
			continue
		}
		seen[port] = true
		ports = append(ports, port)
	}

	return
}

func makeLsofErr(err error) error {
	exit, ok := err.(*exec.ExitError)
	if !ok {
		return err
	}
	status := exit.Sys().(syscall.WaitStatus).ExitStatus()
	if status == 1 && len(exit.Stderr) == 0 {
		return nil
	}
	return fmt.Errorf("%v: %s", exit, exit.Stderr)
}
