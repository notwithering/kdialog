package kdialog

func newTextButtonForm(form Form, text string, d DialogBox) (Button, error) {
	d.Form = form
	d.Text = text

	n, err := RunDialog(d)
	if err != nil {
		return 0, err
	}
	return n.(Button), err
}

func NewYesNo(text string, d DialogBox) (Button, error) {
	return newTextButtonForm(YesNo, text, d)
}

func NewYesNoCancel(text string, d DialogBox) (Button, error) {
	return newTextButtonForm(YesNoCancel, text, d)
}

func NewWarningYesNo(text string, d DialogBox) (Button, error) {
	return newTextButtonForm(WarningYesNo, text, d)
}

func NewWarningContinueCancel(text string, d DialogBox) (Button, error) {
	return newTextButtonForm(WarningContinueCancel, text, d)
}

func NewWarningYesNoCancel(text string, d DialogBox) (Button, error) {
	return newTextButtonForm(WarningYesNoCancel, text, d)
}

func newTextForm(form Form, text string, d DialogBox) error {
	d.Form = form
	d.Text = text

	_, err := RunDialog(d)
	return err
}

func NewSorry(text string, d DialogBox) error {
	return newTextForm(Sorry, text, d)
}

func newDetailedTextForm(form Form, text, details string, d DialogBox) error {
	d.Form = form
	d.Text = text
	d.Details = details

	_, err := RunDialog(d)
	return err
}

func NewDetailedSorry(text, details string, d DialogBox) error {
	return newDetailedTextForm(DetailedSorry, text, details, d)
}

func NewError(text string, d DialogBox) error {
	return newTextForm(Error, text, d)
}

func NewDetailedError(text, details string, d DialogBox) error {
	return newDetailedTextForm(DetailedError, text, details, d)
}

// MsgBox
// InputBox
// ImgBox
// ImgInputBox
// Password
// NewPassword
// TextBox
// TextInputBox
// ComboBox
// Menu
// Checklist
// Radiolist
// PassivePopup
// OpenFile
// SaveFile
// OpenExistingDirectory
// OpenIcon
// Progress
// PickColor
// Slider
// Calender
