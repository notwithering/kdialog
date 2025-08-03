package kdialog

import (
	"fmt"
)

type DontAgain struct {
	File  string
	Entry string
}

func (c DontAgain) String() string {
	if c.File == "" || c.Entry == "" {
		return ""
	}

	return fmt.Sprintf("%s:%s", c.File, c.Entry)
}

func NewDontAgain(file, entry string) DontAgain {
	return DontAgain{
		File:  file,
		Entry: entry,
	}
}
