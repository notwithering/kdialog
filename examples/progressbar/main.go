package main

import (
	"time"

	"github.com/notwithering/kdialog"
)

func main() {
	db := kdialog.DialogBox{
		Form:    kdialog.FormProgressBar,
		Text:    "Loading...",
		Minimum: 0,
		Maximum: 1000,
	}
	bar := db.MustRun().(*kdialog.ProgressBar)

	for i := range db.Maximum + 1 {
		if c, _ := bar.Cancelled(); c {
			return
		}

		bar.SetProgress(i)
		time.Sleep(2 * time.Second / time.Duration(db.Maximum))
	}

	time.Sleep(time.Second)
	bar.Quit()
}
