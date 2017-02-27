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
	err = lsofErr(err)
	if err != nil {
		return
	}
	split := bufio.NewScanner(bytes.NewReader(out))
	split.Scan() // skip headers
	seen := map[int]struct{}{}
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
		seen[port] = struct{}{}
		ports = append(ports, port)
	}
	err = split.Err()
	return
}

func lsofErr(err error) error {
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
