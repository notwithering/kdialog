package main

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/notwithering/kdialog"
)

func main() {
	a := rand.Intn(10) + 1
	b := rand.Intn(10) + 1
	answer := strconv.Itoa(a + b)

	d := kdialog.DialogBox{
		Title: "Math Quiz",
		Form:  kdialog.InputBox,
		Text:  fmt.Sprintf("What is %d + %d?", a, b),
	}

	result := d.MustRun()
	userAnswer := result.(string)

	resp := kdialog.DialogBox{
		Title: "Result",
		Form:  kdialog.MsgBox,
	}

	if userAnswer == answer {
		resp.Text = "Correct!"
		resp.OkLabel = "Yay!"
	} else {
		resp.Text = fmt.Sprintf("Wrong. The answer was %s.", answer)
		resp.OkLabel = "Dang!"
	}

	resp.MustRun()
}
