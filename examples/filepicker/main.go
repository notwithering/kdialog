package main

import (
	"fmt"

	"github.com/notwithering/kdialog"
)

func main() {
	d := kdialog.DialogBox{
		Title: "Open File",
		Form:  kdialog.FormOpenFile,
		FileFilters: kdialog.FileFilters{
			{Name: "Images", Patterns: []string{"*.png", "*.jpg", "*.jpeg"}},
			{Name: "Text Files", Patterns: []string{"*.txt", "*.md"}},
			{Name: "All Files", Patterns: []string{"*"}},
		},
	}

	result, err := d.Run()
	if err != nil {
		fmt.Println(err)
		return
	}

	path := result.(string)
	fmt.Println("Selected file:", path)
}
