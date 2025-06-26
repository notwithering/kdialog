package main

import (
	"github.com/notwithering/kdialog"
)

func main() {
	d := kdialog.DialogBox{
		Title: "Love Letter",
		Form:  kdialog.YesNo,
		Text:  "Do you love me?",
	}

	result := d.MustRun()        // or use d.Run() and handle errors
	r := result.(kdialog.Button) // assert the type to a button

	switch r {
	case kdialog.Yes:
		d.Text = "I love you too!"
	case kdialog.No:
		d.Text = "I don't love you!"
	default:
		// most likely kdialog.Cancel which can be done with a YesNo dialog by just closing the window
		return
	}

	// reuse the dialog box
	d.Form = kdialog.MsgBox
	d.MustRun()
}
