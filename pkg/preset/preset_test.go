package preset

import (
	"slices"
	"testing"

	"github.com/i-zaitsev/gitcat/pkg/filter"
	"github.com/i-zaitsev/pkg/assert"
)

func TestNew_One(t *testing.T) {
	tests := []struct {
		opt        Option
		wantKeep   string
		wantIgnore string
	}{
		{
			Go,
			".go",
			"vendor",
		},
		{
			Python,
			".py",
			"__pycache__",
		},
		{
			C,
			".c",
			"*.o",
		},
		{
			Cpp,
			".cpp",
			"*.o",
		},
	}
	for _, tt := range tests {
		t.Run(tt.opt.String(), func(t *testing.T) {
			p := New(tt.opt)

			assert.True(t, slices.Contains(p.Keep, tt.wantKeep))
			assert.True(t, slices.Contains(p.Ignore, tt.wantIgnore))
		})
	}
}

func TestNew_Many_IsUnionWithoutDuplicates(t *testing.T) {
	p := New(Go, C, Go)

	assert.True(t, slices.Contains(p.Keep, ".go"))
	assert.True(t, slices.Contains(p.Keep, ".c"))
	assert.True(t, slices.Contains(p.Ignore, "vendor"))
	assert.True(t, slices.Contains(p.Ignore, "*.o"))
	assert.Eq(t, slices.Compact(slices.Clone(p.Keep)), p.Keep)
	assert.Eq(t, slices.Compact(slices.Clone(p.Ignore)), p.Ignore)
}

func TestNew_All(t *testing.T) {
	all := New(All()...)

	for _, opt := range All() {
		p := New(opt)
		for _, pattern := range p.Keep {
			assert.True(t, slices.Contains(all.Keep, pattern))
		}
		for _, pattern := range p.Ignore {
			assert.True(t, slices.Contains(all.Ignore, pattern))
		}
	}
}

func TestNew_None(t *testing.T) {
	p := New()

	assert.Len(t, p.Keep, 0)
	assert.Len(t, p.Ignore, 0)
}

func TestNew_ValidPatterns(t *testing.T) {
	p := New(All()...)

	for _, pattern := range slices.Concat(p.Keep, p.Ignore) {
		assert.NotErr(t, filter.CheckPattern(pattern))
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		value string
		want  []Option
	}{
		{
			"go",
			[]Option{Go},
		},
		{
			"c|cpp|go",
			[]Option{C, Cpp, Go},
		},
		{
			" python, go ,",
			[]Option{Python, Go},
		},
		{
			"all",
			All(),
		},
		{
			"",
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			got, err := Parse(tt.value)

			assert.NotErr(t, err)
			assert.Eq(t, got, tt.want)
		})
	}
}

func TestParse_Unknown(t *testing.T) {
	_, err := Parse("go|cobol")

	assert.ErrContains(t, err, `unknown preset "cobol": must be one of: go, python, c, cpp, all`)
}

func TestOption_String(t *testing.T) {
	assert.Eq(t, Go.String(), "go")
	assert.Eq(t, Cpp.String(), "cpp")
}
