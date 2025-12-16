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

	DontAgain DontAgain
	Geometry  Geometry
	Attach    string

	Minimum  int
	Maximum  int
	Interval int

	OkLabel       string
	YesLabel      string
	NoLabel       string
	CancelLabel   string
	ContinueLabel string
}

func (db DialogBox) Run() (result any, err error) {
	var args []string

	// appends a list of values to args
	appendArgs := func(a ...any) {
		for _, v := range a {
			args = append(args, fmt.Sprint(v))
		}
	}

	// add a, b to args if b is not empty or true
	addFlagIf := func(a string, b any) {
		switch v := b.(type) {
		case string:
			if v != "" {
				appendArgs(a, v)
			}
		case bool:
			if v {
				appendArgs(a)
			}
		}
	}

	// run kdialog with options and return stdout, exit code, and error
	runDialog := func(a ...any) (string, int, error) {
		appendArgs(a...)
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

	// get the n index of opts if out of range it will return Undefined
	getButton := func(n int, opts ...Button) Button {
		if n < 0 || n >= len(opts) {
			return Undefined
		}
		return opts[n]
	}

	// runs as a listed dialog box
	// if multiple is true, it will return a runListDialog of selected items as a []int
	// if multiple is false, it will return the selected item as an int
	runListDialog := func(tagged, checks, multiple bool) (any, error) {
		var tags []string

		for i, item := range db.Items {
			if tagged {
				tag := fmt.Sprint(i)
				tags = append(tags, tag)
				appendArgs(tag)
			} else {
				tags = append(tags, item)
			}
			appendArgs(item)

			if checks {
				if i < len(db.Checks) {
					if db.Checks[i] {
						appendArgs("on")
					} else {
						appendArgs("off")
					}
				} else {
					appendArgs("off")
				}
			}
		}

		if multiple {
			var checkedTags []int

			msg, _, err := runDialog("--separate-output")
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

		msg, _, err := runDialog()
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

	addFlagIf("--title", db.Title)
	addFlagIf("--default", db.Default)
	addFlagIf("--multiple", db.Multiple)
	addFlagIf("--dontagain", db.DontAgain)
	addFlagIf("--geometry", db.Geometry.String())
	addFlagIf("--attach", db.Attach)
	addFlagIf("--ok-label", db.OkLabel)
	addFlagIf("--yes-label", db.YesLabel)
	addFlagIf("--no-label", db.NoLabel)
	addFlagIf("--cancel-label", db.CancelLabel)
	addFlagIf("--continue-label", db.ContinueLabel)

	switch db.Form {
	case YesNo:
		_, code, err := runDialog("--yesno", db.Text, db.Details)
		if err != nil {
			return nil, err
		}
		return getButton(code, Yes, No, Cancel), nil
	case YesNoCancel:
		_, code, err := runDialog("--yesnocancel", db.Text, db.Details)
		if err != nil {
			return nil, err
		}
		return getButton(code, Yes, No, Cancel), nil
	case WarningYesNo:
		_, code, err := runDialog("--warningyesno", db.Text, db.Details)
		if err != nil {
			return nil, err
		}
		return getButton(code, Yes, No, Cancel), nil
	case WarningContinueCancel:
		_, code, err := runDialog("--warningcontinuecancel", db.Text, db.Details)
		if err != nil {
			return nil, err
		}
		return getButton(code, Continue, Cancel), nil
	case WarningYesNoCancel:
		_, code, err := runDialog("--warningyesnocancel", db.Text, db.Details)
		if err != nil {
			return nil, err
		}
		return getButton(code, Yes, No, Cancel), nil
	case Sorry:
		_, _, err := runDialog("--sorry", db.Text, db.Details)
		return nil, err
	case Error:
		_, _, err := runDialog("--error", db.Text, db.Details)
		return nil, err
	case MsgBox:
		_, _, err := runDialog("--msgbox", db.Text, db.Details)
		return nil, err
	case InputBox:
		msg, _, err := runDialog("--inputbox", db.Text, db.InitialText)
		return msg, err
	case ImgBox:
		_, _, err := runDialog("--imgbox", db.FilePath)
		return nil, err
	case ImgInputBox:
		msg, _, err := runDialog("--imginputbox", db.FilePath, db.Text)
		return msg, err
	case Password:
		msg, _, err := runDialog("--password", db.Text)
		return msg, err
	case NewPassword:
		msg, _, err := runDialog("--newpassword", db.Text)
		return msg, err
	case TextBox:
		msg, _, err := runDialog("--textbox", db.FilePath)
		return msg, err
	case TextInputBox:
		msg, _, err := runDialog("--textinputbox", db.Text, db.InitialText)
		return msg, err
	case ComboBox:
		appendArgs("--combobox", db.Text)
		tagged := false
		checks := false
		multiple := false
		return runListDialog(tagged, checks, multiple)
	case Menu:
		appendArgs("--menu", db.Text)
		tagged := true
		checks := false
		multiple := false
		return runListDialog(tagged, checks, multiple)
	case Checklist:
		appendArgs("--checklist", db.Text)
		tagged := true
		checks := true
		multiple := true
		return runListDialog(tagged, checks, multiple)
	case Radiolist:
		appendArgs("--radiolist", db.Text)
		tagged := true
		checks := true
		multiple := false
		return runListDialog(tagged, checks, multiple)
	case PassivePopup:
		_, _, err := runDialog("--passivepopup", db.Text, db.Timeout)
		return nil, err
	case OpenFile:
		msg, _, err := runDialog("--getopenfilename", db.StartDir, db.FileFilters.String())
		return msg, err
	case SaveFile:
		msg, _, err := runDialog("--getsavefilename", db.StartDir, db.FileFilters.String())
		return msg, err
	case OpenExistingDirectory:
		msg, _, err := runDialog("--getexistingdirectory", db.StartDir)
		return msg, err
	case OpenIcon:
		msg, _, err := runDialog("--geticon", db.Group, db.Context)
		return msg, err
	case ProgressBar:
		msg, _, err := runDialog("--progressbar", db.Text, db.Maximum)
		if err != nil {
			return nil, err
		}

		conn, err := dbus.SessionBus()
		if err != nil {
			return nil, err
		}

		obj := conn.Object(strings.Split(msg, " ")[0], "/ProgressDialog")
		return &ProgressBarResult{obj: obj}, nil
	case PickColor:
		msg, _, err := runDialog("--getcolor")
		if err != nil {
			return nil, err
		}

		if msg == "" {
			return Cancel, nil
		}

		var c color.RGBA
		c.A = 0xff

		if _, err := fmt.Sscanf(msg, "#%02x%02x%02x", &c.R, &c.G, &c.B); err != nil {
			return nil, err
		}

		return c, nil
	case Slider:
		msg, _, err := runDialog("--slider", db.Text, db.Minimum, db.Maximum, db.Interval)
		if err != nil {
			return nil, err
		}

		if msg == "" {
			return Cancel, nil
		}

		n, err := strconv.Atoi(msg)
		if err != nil {
			return nil, err
		}

		return n, nil
	case Calender:
		msg, _, err := runDialog("--calendar", db.Text, "--dateformat", "yyyy-MM-dd")
		if err != nil {
			return nil, err
		}

		if msg == "" {
			return Cancel, nil
		}

		date, err := time.Parse("2006-01-02", msg)
		if err != nil {
			return nil, err
		}

		return date, nil
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
