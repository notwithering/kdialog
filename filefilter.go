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
