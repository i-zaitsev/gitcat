package gitpath

import (
	"fmt"
	"strconv"
)

// Size represents a file size in bytes, parsed from KB values.
type Size int64

func (s *Size) String() string {
	if s == nil {
		return "0"
	}
	return fmt.Sprintf("%d", *s)
}

func (s *Size) Set(value string) error {
	if value == "" {
		*s = 0
		return nil
	}

	kb, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid size: must be integer (KB)")
	}

	*s = Size(kb * 1024)
	return nil
}

// InBytes returns the size as int64 bytes.
func (s *Size) InBytes() int64 {
	if s == nil {
		return 0
	}
	return int64(*s)
}
