// Package preset names common sets of rules that select the files of a
// project written in one language: the sources and docs to keep, and the
// build artifacts and caches to ignore. The rules are patterns for
// [filter.Filter.Keep] and [filter.Filter.Ignore].
package preset

import (
	"fmt"
	"slices"
	"strings"
)

// Preset is a set of keep and ignore patterns.
type Preset struct {
	Keep   []string
	Ignore []string
}

// Option is the preset of one language.
type Option int

// Presets.
const (
	Go Option = iota
	Python
	C
	Cpp
)

// names are the names of the options, indexed by option.
var names = []string{"go", "python", "c", "cpp"}

// all is the name that stands for every option in [Parse].
const all = "all"

// All returns every option.
func All() []Option {
	return []Option{Go, Python, C, Cpp}
}

// String returns the name of the option.
func (o Option) String() string {
	return names[o]
}

// New returns the preset with the rules of every option in opts, sorted
// and without duplicates. With no options, the preset is empty.
func New(opts ...Option) *Preset {
	var p Preset
	for _, opt := range opts {
		r := one(opt)
		p.Keep = append(p.Keep, r.Keep...)
		p.Ignore = append(p.Ignore, r.Ignore...)
	}
	slices.Sort(p.Keep)
	slices.Sort(p.Ignore)
	p.Keep, p.Ignore = slices.Compact(p.Keep), slices.Compact(p.Ignore)
	return &p
}

// one returns the rules of the single option opt.
func one(opt Option) *Preset {
	switch opt {
	case Go:
		return &Preset{
			Keep:   []string{".go", "go.mod", ".md"},
			Ignore: []string{"vendor", "bin"},
		}
	case Python:
		return &Preset{
			Keep:   []string{".py", ".pyi", ".md", ".rst", "pyproject.toml"},
			Ignore: []string{"*.py[cod]", "__pycache__", "*.egg-info", "venv", ".venv", "build", "dist"},
		}
	case C:
		return &Preset{
			Keep:   []string{".c", ".h", "Makefile", "CMakeLists.txt", ".md"},
			Ignore: []string{"build", "cmake-build-*", "*.o", "*.a", "*.so"},
		}
	case Cpp:
		return &Preset{
			Keep:   []string{".cpp", ".cc", ".hpp", ".h", "Makefile", "CMakeLists.txt", ".md"},
			Ignore: []string{"build", "cmake-build-*", "*.o", "*.a", "*.so"},
		}
	}
	return &Preset{}
}

// Names returns the names accepted by [Parse].
func Names() []string {
	return append(slices.Clone(names), all)
}

// Parse returns the options named in value, separated by pipe or comma.
// The name "all" stands for every option. Empty names are skipped.
func Parse(value string) ([]Option, error) {
	var opts []Option
	for item := range strings.FieldsFuncSeq(value, func(r rune) bool { return r == '|' || r == ',' }) {
		name := strings.TrimSpace(item)
		if name == "" {
			continue
		}
		if name == all {
			opts = append(opts, All()...)
			continue
		}
		i := slices.Index(names, name)
		if i < 0 {
			return nil, fmt.Errorf("unknown preset %q: must be one of: %s", name, strings.Join(Names(), ", "))
		}
		opts = append(opts, Option(i))
	}
	return opts, nil
}
