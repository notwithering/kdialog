package main

import (
	"strconv"

	"github.com/notwithering/kdialog"
)

func main() {
	_, err := strconv.Atoi("hello")

	kdialog.DialogBox{
		Form:    kdialog.Error,
		Text:    "An error occured.",
		Details: err.Error(),
	}.MustRun()
}
