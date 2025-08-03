package kdialog

import (
	"strings"
)

type FileFilter struct {
	Name     string
	Patterns []string
}

func (f FileFilter) String() string {
	return f.Name + " (" + strings.Join(f.Patterns, " ") + ")"
}

func NewFileFilter(name string, patterns ...string) FileFilter {
	return FileFilter{
		Name:     name,
		Patterns: patterns,
	}
}

type FileFilters []FileFilter

func (f FileFilters) String() string {
	var b strings.Builder

	for i, filter := range f {
		if i > 0 {
			b.WriteString(";;")
		}
		b.WriteString(filter.String())
	}

	return b.String()
}

func NewFileFilters(filters ...FileFilter) FileFilters {
	return FileFilters(filters)
}
