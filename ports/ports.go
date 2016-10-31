package ports

import (
	"bufio"
	"bytes"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
)

var listenRe *regexp.Regexp

func init() {
	listenRe = regexp.MustCompile(`.*:(\d+) \(LISTEN\)`)
}

type Filter interface {
	Match(int) bool
}

func List(f Filter) ([]int, error) {
	cmd := exec.Command("lsof", "-l", "-P", "-n", "-i", "TCP", "-s", "TCP:LISTEN")

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	ports := []int{}

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
		if !f.Match(port) {
			continue
		}
		ports = append(ports, port)
	}

	sort.Ints(ports)

	return ports, nil
}

type Between struct {
	Min, Max int
}

func (f Between) Match(port int) bool {
	return port >= f.Min && port <= f.Max
}

type Or []Filter

func (f Or) Match(port int) bool {
	for _, filter := range f {
		if filter.Match(port) {
			return true
		}
	}
	return false
}
