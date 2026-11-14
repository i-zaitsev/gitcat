// Package utils provides small filesystem and cleanup helpers shared by the
// gitcat packages.
package utils

import (
	"io"
	"os"
)

// DirExists reports whether dirname exists and is a directory.
func DirExists(dirname string) bool {
	info, err := os.Stat(dirname)
	return err == nil && info.IsDir()
}

// SilentClose closes c and discards the error. It is meant for deferred
// cleanup where a close error cannot be acted upon.
func SilentClose(c io.Closer) {
	_ = c.Close()
}
