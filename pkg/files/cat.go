package files

import (
	"bufio"
	"bytes"
	"os"
	"sync"

	"github.com/i-zaitsev/gitcat/pkg/internal/utils"
	"github.com/i-zaitsev/gitcat/pkg/log"
)

type concat struct {
	paths []string
	lines map[string][]string
	mu    sync.Mutex
}

// Cat reads files and concatenates their contents.
// If maxLines > 0, only the first maxLines of each file are read.
func Cat(maxLines int, paths ...string) string {
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
				log.Warn("failed to open file", "path", p, "error", err)
				return
			}
			defer utils.SilentClose(f)
			log.Debug("reading file", "path", p)
			scanner := bufio.NewScanner(f)
			lineCount := 0
			for scanner.Scan() {
				if maxLines > 0 && lineCount >= maxLines {
					break
				}
				cc.mu.Lock()
				cc.lines[p] = append(cc.lines[p], scanner.Text()+"\n")
				cc.mu.Unlock()
				lineCount++
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
	return buf.String()
}
