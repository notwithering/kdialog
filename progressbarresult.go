package kdialog

import (
	"github.com/godbus/dbus/v5"
)

type ProgressBarResult struct {
	obj dbus.BusObject
}

func (p *ProgressBarResult) SetProgress(prog int) error {
	call := p.obj.Call("org.freedesktop.DBus.Properties.Set", 0, "org.kde.kdialog.ProgressDialog", "value", dbus.MakeVariant(prog))
	if call.Err != nil && call.Err.Error() != "The name is not activatable" {
		return call.Err
	}

	return nil
}
func (p *ProgressBarResult) Cancelled() (bool, error) {
	var cancelled bool

	call := p.obj.Call("org.kde.kdialog.ProgressDialog.wasCancelled", 0)
	if call.Err != nil && call.Err.Error() != "The name is not activatable" {
		return false, call.Err
	}

	err := call.Store(&cancelled)
	if err != nil {
		return true, nil
	}

	return cancelled, err
}
func (p *ProgressBarResult) Quit() error {
	return p.obj.Call("org.kde.kdialog.ProgressDialog.close", 0).Err
}
