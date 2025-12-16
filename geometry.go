package kdialog

import (
	"fmt"
	"strings"
)

type Geometry struct {
	Width, Height    int
	XOffset, YOffset int
}

func (g Geometry) String() string {
	var b strings.Builder

	if g.Width != 0 && g.Height != 0 {
		fmt.Fprintf(&b, "%dx%d", g.Width, g.Height)
	}
	if g.XOffset != 0 || g.YOffset != 0 {
		fmt.Fprintf(&b, "%+d%+d", g.XOffset, g.YOffset)
	}

	return b.String()
}

func NewGeometry(width, height, xoffset, yoffset int) Geometry {
	return Geometry{
		Width:   width,
		Height:  height,
		XOffset: xoffset,
		YOffset: yoffset,
	}
}
