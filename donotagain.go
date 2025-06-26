package kdialog

import (
	"fmt"
)

type DoNotAgain struct {
	File  string
	Entry string
}

func (c *DoNotAgain) String() string {
	if c.File == "" || c.Entry == "" {
		return ""
	}

	return fmt.Sprintf("%s:%s", c.File, c.Entry)
}
