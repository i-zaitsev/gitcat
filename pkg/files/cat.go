package files

import (
	"bufio"
	"bytes"
	"os"
	"strings"
	"sync"

	"github.com/i-zaitsev/gitcat/pkg/internal/utils"
)

type concat struct {
	paths []string
	lines map[string][]string
	mu    sync.Mutex
}

// Cat reads files and concatenates their contents.
func Cat(paths []string) []string {
	cc := concat{
		paths: paths,
		lines: make(map[string][]string, len(paths)),
	}
	var wg sync.WaitGroup
	for _, path := range paths {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			f, err := os.Open(p)
			if err != nil {
				return
			}
			defer utils.SilentClose(f)
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				cc.mu.Lock()
				cc.lines[p] = append(cc.lines[p], scanner.Text()+"\n")
				cc.mu.Unlock()
			}
		}(path)
	}
	wg.Wait()
	var buf bytes.Buffer
	for _, path := range paths {
		for _, line := range cc.lines[path] {
			buf.WriteString(line)
		}
	}
	return strings.Split(buf.String(), "\n")
}
