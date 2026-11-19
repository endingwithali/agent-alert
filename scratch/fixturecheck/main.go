package main

import (
	"fmt"
	"os"

	"agent-monitor/internal/sessionpoller"
)

func main() {
	p, err := sessionpoller.New("internal/sessionpoller/testdata")
	if err != nil {
		panic(err)
	}
	s, err := p.Poll()
	if err != nil {
		panic(err)
	}
	fmt.Printf("path=%s bytes=%d\n", s.Path, len(s.Data))
	os.Exit(0)
}
