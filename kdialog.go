// Package kdialog is the first and most complete wrapper for the KDialog command in Go.
package kdialog

import (
	"fmt"
	"image/color"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/godbus/dbus/v5"
)

type Button int
type ProgressResult struct {
	SetValue chan int
	Quit     chan struct{}
}

const (
	Ok Button = iota
	Yes
	No
	Cancel
	Continue
)

func RunDialog(db DialogBox) (result any, code int) {
	var args []string

	// add list of strings to args
	add := func(s ...string) {
		args = append(args, s...)
	}

	// add a, b to args if b is not empty or true
	flag := func(a string, b any) {
		switch v := b.(type) {
		case string:
			if v != "" {
				add(a, v)
			}
		case bool:
			if v {
				add(a)
			}
		}
	}

	// run kdialog with options and return exit code and stdout
	run := func(s ...string) (int, string) {
		add(s...)
		var stdout strings.Builder
		cmd := exec.Command("kdialog", args...)
		cmd.Stdout = &stdout

		cmd.Run()

		return cmd.ProcessState.ExitCode(), strings.TrimRight(stdout.String(), "\n")
	}

	// get the codeth index of opts if out of range it will return 0
	get := func(code int, opts ...Button) Button {
		if code < len(opts) && code >= 0 {
			return opts[code]
		}
		return 0
	}

	list := func() (any, int) {
		var tags []string

		for i, item := range db.Items {
			tag := fmt.Sprint(i)

			tags = append(tags, tag)
			add(tag, item)

			if i < len(db.Checks) {
				if db.Checks[i] {
					add("on")
				} else {
					add("off")
				}
			} else {
				add("off")
			}
		}

		var checkedTags []int

		fmt.Println(args)

		code, msg := run("--separate-output")

		for _, msgTag := range strings.Split(msg, "\n") {
			for i, tag := range tags {
				if msgTag == tag {
					checkedTags = append(checkedTags, i)
				}
			}
		}

		return checkedTags, code
	}

	flag("--ok-label", db.Ok)
	flag("--yes-label", db.Yes)
	flag("--no-label", db.No)
	flag("--cancel-label", db.Cancel)
	flag("--continue-label", db.Continue)
	flag("--multiple", db.Multiple)
	flag("--title", db.Title)
	flag("--default", db.Default)
	flag("--geometry", db.Geometry)

	switch db.Form {
	case YesNo:
		code, _ := run("--yesno", db.Text)
		return get(code, Yes, No), code
	case YesNoCancel:
		code, _ := run("--yesnocancel", db.Text)
		return get(code, Yes, No, Cancel), code
	case WarningYesNo:
		code, _ := run("--warningyesno", db.Text)
		return get(code, Yes, No), code
	case WarningContinueCancel:
		code, _ := run("--warningcontinuecancel", db.Text)
		return get(code, Continue, Cancel), code
	case WarningYesNoCancel:
		code, _ := run("--warningyesnocancel", db.Text)
		return get(code, Yes, No, Cancel), code
	case Sorry:
		code, _ := run("--sorry", db.Text)
		return nil, code
	case DetailedSorry:
		code, _ := run("--detailedsorry", db.Text, db.Details)
		return nil, code
	case MsgBox:
		code, _ := run("--msgbox", db.Text, db.Details)
		return nil, code
	case InputBox:
		code, msg := run("--inputbox", db.Text, db.InitialText)
		return msg, code
	case ImgBox:
		code, _ := run("--imgbox", db.FilePath)
		return nil, code
	case ImgInputBox:
		code, msg := run("--imginputbox", db.FilePath, db.Text)
		return msg, code
	case Password:
		code, msg := run("--password", db.Text)
		return msg, code
	case NewPassword:
		code, msg := run("--newpassword", db.Text)
		return msg, code
	case TextBox:
		code, _ := run("--textbox", db.FilePath)
		return nil, code
	case TextInputBox:
		code, msg := run("--textinputbox", db.Text, db.InitialText)
		return msg, code
	case ComboBox:
		add("--combobox", db.Text)
		code, msg := run(db.Items...)
		for i, item := range db.Items {
			if item == msg {
				return i, code
			}
		}
		return nil, code
	case Menu:
		code, _ := run("--menu", db.Text)
		return nil, code
	case Checklist:
		add("--checklist", db.Text)
		list()
	case Radiolist:
		add("--radiolist", db.Text)
		list()
	case PassivePopup:
		code, _ := run("--passivepopup", db.Text, fmt.Sprint(db.Timeout))
		return nil, code
	case OpenFile:
		code, msg := run("--getopenfilename", db.StartDir, db.FileFilter)
		return msg, code
	case SaveFile:
		code, msg := run("--getsavefilename", db.StartDir, db.FileFilter)
		return msg, code
	case OpenExistingDirectory:
		code, msg := run("--getexistingdirectory", db.StartDir)
		return msg, code
	case OpenIcon:
		code, msg := run("--geticon", db.Group, db.Context)
		return msg, code
	case Progress: // FIXME: what the fuck is this shit
		code, msg := run("--progressbar", db.Text, fmt.Sprint(db.Maximum))

		conn, err := dbus.SessionBus()
		if err != nil {
			fmt.Println(err)
			return nil, code
		}

		channel := make(chan int)
		quit := make(chan struct{})

		obj := conn.Object(strings.Split(msg, " ")[0], "/ProgressDialog")

		go func() {
			for {
				select {
				case <-quit:
					obj.Call("org.kde.kdialog.ProgressDialog.close", 0)
					return
				case prog, ok := <-channel:
					if !ok {
						close(quit)
						return
					}

					obj.Call("org.freedesktop.DBus.Properties.Set", 0, "org.kde.kdialog.ProgressDialog", "value", dbus.MakeVariant(prog))
				default:
					var wasCancelled bool
					call := obj.Call("org.kde.kdialog.ProgressDialog.wasCancelled", 0)
					if call.Err != nil {
						close(quit)
						return
					}
					call.Store(&wasCancelled)

					if wasCancelled {
						close(quit)
						return
					}
				}
			}
		}()

		return ProgressResult{SetValue: channel, Quit: quit}, code
	case PickColor:
		code, msg := run("--getcolor")

		var c color.RGBA
		c.A = 0xff
		fmt.Sscanf(msg, "#%1x%1x%1x", &c.R, &c.G, &c.B)

		c.R *= 17
		c.G *= 17
		c.B *= 17

		return c, code
	case Slider:
		code, msg := run("--slider", db.Text, fmt.Sprint(db.Minimum), fmt.Sprint(db.Maximum), fmt.Sprint(db.Interval))
		n, _ := strconv.Atoi(msg)
		return n, code
	case Calender:
		code, msg := run("--calendar", db.Text, "--dateformat", "yyyy-MM-dd")
		date, _ := time.Parse("2006-01-02", msg)
		return date, code
	}

	return nil, 0
}
