package cpu

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
)

type topLine struct {
	Pct     float64
	Command string
}

func Top(n int) ([]topLine, error) {
	uid := fmt.Sprintf("%d", os.Getuid())
	cmd := exec.Command("ps", "-U", uid, "-e", "-o", "pid,%cpu,comm")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	topLines := []topLine{}
	curPid := os.Getpid()
	split := bufio.NewScanner(bytes.NewReader(out))
	split.Scan() // skip headers
	for split.Scan() {
		var (
			topl topLine
			pid  int
		)
		_, err := fmt.Sscanf(split.Text(), "%d %f %s",
			&pid, &topl.Pct, &topl.Command,
		)
		if err != nil {
			continue
		}
		if pid == curPid {
			continue
		}
		topLines = append(topLines, topl)
	}
	if l := len(topLines); n > l {
		n = l
	}
	sort.Sort(byPct(topLines))
	return topLines[:n], nil
}

type byPct []topLine

func (s byPct) Len() int           { return len(s) }
func (s byPct) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s byPct) Less(i, j int) bool { return s[j].Pct < s[i].Pct }
