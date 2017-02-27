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
	var (
		split    = bufio.NewScanner(bytes.NewReader(out))
		topLines = []topLine{}
		topl     topLine
		pid      int
		curPid   = os.Getpid()
	)
	split.Scan() // skip headers
	for split.Scan() {
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
	pct := func(i int) float64 {
		return topLines[i].Pct
	}
	sort.Slice(topLines, func(i, j int) bool {
		return pct(i) > pct(j)
	})
	return topLines[:n], nil
}
