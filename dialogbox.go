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

type DialogBox struct {
	Form Form

	Title string

	Text        string
	InitialText string
	Details     string
	FilePath    string

	Timeout int

	StartDir    string
	FileFilters FileFilters
	Group       string
	Context     string

	Default string

	Multiple bool

	Items  []string
	Checks []bool

	DontAgain DoNotAgain
	Geometry  Geometry
	Attach    string

	Minimum  int
	Maximum  int
	Interval int

	Ok       string
	Yes      string
	No       string
	Cancel   string
	Continue string
}

func (db DialogBox) Run() (result any, err error) {
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

	// run kdialog with options and return stdout, exit code, and error
	run := func(s ...string) (string, int, error) {
		add(s...)
		var stdout strings.Builder
		cmd := exec.Command("kdialog", args...)
		cmd.Stdout = &stdout

		err := cmd.Run()
		_, isExitError := err.(*exec.ExitError)
		if err != nil && !isExitError {
			return "", 0, err
		}

		return strings.TrimRight(stdout.String(), "\n"), cmd.ProcessState.ExitCode(), nil
	}

	// get the n index of opts if out of range it will return -1
	get := func(n int, opts ...Button) Button {
		if n < 0 || n >= len(opts) {
			return -1
		}
		return opts[n]
	}

	// runs as a listed dialog box
	// if multiple is true, it will return a list of selected items as a []int
	// if multiple is false, it will return the selected item as an int
	list := func(tagged, checks, multiple bool) (any, error) {
		var tags []string

		for i, item := range db.Items {
			if tagged {
				tag := fmt.Sprint(i)
				tags = append(tags, tag)
				add(tag)
			} else {
				tags = append(tags, item)
			}
			add(item)

			if checks {
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
		}

		if multiple {
			var checkedTags []int

			msg, _, err := run("--separate-output")
			if err != nil {
				return nil, err
			}

			for _, msgTag := range strings.Split(msg, "\n") {
				for i, tag := range tags {
					if msgTag == tag {
						checkedTags = append(checkedTags, i)
					}
				}
			}

			return checkedTags, nil
		}

		msg, _, err := run()
		if err != nil {
			return nil, err
		}

		msgTag := strings.TrimRight(msg, "\n")

		for i, tag := range tags {
			if msgTag == tag {
				return i, nil
			}
		}

		return -1, nil
	}

	flag("--title", db.Title)
	flag("--default", db.Default)
	flag("--multiple", db.Multiple)
	flag("--dontagain", db.DontAgain)
	flag("--geometry", db.Geometry.String())
	flag("--attach", db.Attach)
	flag("--ok-label", db.Ok)
	flag("--yes-label", db.Yes)
	flag("--no-label", db.No)
	flag("--cancel-label", db.Cancel)
	flag("--continue-label", db.Continue)

	switch db.Form {
	case YesNo:
		_, code, err := run("--yesno", db.Text, db.Details)
		if err != nil {
			return nil, err
		}
		return get(code, Yes, No, Cancel), nil
	case YesNoCancel:
		_, code, err := run("--yesnocancel", db.Text, db.Details)
		if err != nil {
			return nil, err
		}
		return get(code, Yes, No, Cancel), nil
	case WarningYesNo:
		_, code, err := run("--warningyesno", db.Text, db.Details)
		if err != nil {
			return nil, err
		}
		return get(code, Yes, No), nil
	case WarningContinueCancel:
		_, code, err := run("--warningcontinuecancel", db.Text, db.Details)
		if err != nil {
			return nil, err
		}
		return get(code, Continue, Cancel), nil
	case WarningYesNoCancel:
		_, code, err := run("--warningyesnocancel", db.Text, db.Details)
		if err != nil {
			return nil, err
		}
		return get(code, Yes, No, Cancel), nil
	case Sorry:
		_, _, err := run("--sorry", db.Text, db.Details)
		return nil, err
	case Error:
		_, _, err := run("--error", db.Text, db.Details)
		return nil, err
	case MsgBox:
		_, _, err := run("--msgbox", db.Text, db.Details)
		return nil, err
	case InputBox:
		msg, _, err := run("--inputbox", db.Text, db.InitialText)
		return msg, err
	case ImgBox:
		_, _, err := run("--imgbox", db.FilePath)
		return nil, err
	case ImgInputBox:
		msg, _, err := run("--imginputbox", db.FilePath, db.Text)
		return msg, err
	case Password:
		msg, _, err := run("--password", db.Text)
		return msg, err
	case NewPassword:
		msg, _, err := run("--newpassword", db.Text)
		return msg, err
	case TextBox:
		msg, _, err := run("--textbox", db.FilePath)
		return msg, err
	case TextInputBox:
		msg, _, err := run("--textinputbox", db.Text, db.InitialText)
		return msg, err
	case ComboBox:
		add("--combobox", db.Text)
		return list(false, false, false)
	case Menu:
		add("--menu", db.Text)
		return list(true, false, false)
	case Checklist:
		add("--checklist", db.Text)
		return list(true, true, true)
	case Radiolist:
		add("--radiolist", db.Text)
		return list(true, true, false)
	case PassivePopup:
		_, _, err := run("--passivepopup", db.Text, fmt.Sprint(db.Timeout))
		return nil, err
	case OpenFile:
		msg, _, err := run("--getopenfilename", db.StartDir, db.FileFilters.String())
		return msg, err
	case SaveFile:
		msg, _, err := run("--getsavefilename", db.StartDir, db.FileFilters.String())
		return msg, err
	case OpenExistingDirectory:
		msg, _, err := run("--getexistingdirectory", db.StartDir)
		return msg, err
	case OpenIcon:
		msg, _, err := run("--geticon", db.Group, db.Context)
		return msg, err
	case ProgressBar:
		msg, _, err := run("--progressbar", db.Text, fmt.Sprint(db.Maximum))
		if err != nil {
			return nil, err
		}

		conn, err := dbus.SessionBus()
		if err != nil {
			return nil, err
		}

		obj := conn.Object(strings.Split(msg, " ")[0], "/ProgressDialog")
		return ProgressBarResult{obj: obj}, nil
	case PickColor:
		msg, _, err := run("--getcolor")
		if err != nil {
			return nil, err
		}

		var c color.RGBA
		c.A = 0xff
		fmt.Sscanf(msg, "#%1x%1x%1x", &c.R, &c.G, &c.B)

		c.R *= 17
		c.G *= 17
		c.B *= 17

		return c, nil
	case Slider:
		msg, _, err := run("--slider", db.Text, fmt.Sprint(db.Minimum), fmt.Sprint(db.Maximum), fmt.Sprint(db.Interval))
		if err != nil {
			return nil, err
		}
		n, _ := strconv.Atoi(msg)
		return n, nil
	case Calender:
		msg, _, err := run("--calendar", db.Text, "--dateformat", "yyyy-MM-dd")
		if err != nil {
			return nil, err
		}
		date, _ := time.Parse("2006-01-02", msg)
		return date, err
	}

	return
}

func (db DialogBox) MustRun() (result any) {
	result, err := db.Run()
	if err != nil {
		panic(err)
	}
	return result
}
