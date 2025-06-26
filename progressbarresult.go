package kdialog

import (
	"github.com/godbus/dbus/v5"
)

type ProgressBarResult struct {
	obj dbus.BusObject
}

const (
	errNoExist = "The name is not activatable"
)

func (p ProgressBarResult) SetProgress(prog int) error {
	call := p.obj.Call("org.freedesktop.DBus.Properties.Set", 0, "org.kde.kdialog.ProgressDialog", "value", dbus.MakeVariant(prog))
	if call.Err != nil {
		if call.Err.Error() == errNoExist {
			return nil
		}
		return call.Err
	}

	return nil
}
func (p ProgressBarResult) Cancelled() (bool, error) {
	var cancelled bool

	call := p.obj.Call("org.kde.kdialog.ProgressDialog.wasCancelled", 0)
	if call.Err != nil {
		if call.Err.Error() == errNoExist {
			return true, nil
		}
		return false, call.Err
	}

	err := call.Store(&cancelled)
	if err != nil {
		return false, err
	}

	return cancelled, err
}
func (p ProgressBarResult) Quit() error {
	err := p.obj.Call("org.kde.kdialog.ProgressDialog.close", 0).Err
	if err != nil {
		if err.Error() == errNoExist {
			return nil
		}
		return err
	}

	return nil
}
