package kdialog

type Form uint8

const (
	YesNo Form = iota
	YesNoCancel
	WarningYesNo
	WarningContinueCancel
	WarningYesNoCancel
	Sorry
	Error
	MsgBox
	InputBox
	ImgBox
	ImgInputBox
	Password
	NewPassword
	TextBox
	TextInputBox
	ComboBox
	Menu
	Checklist
	Radiolist
	PassivePopup
	OpenFile
	SaveFile
	OpenExistingDirectory
	OpenIcon
	ProgressBar
	PickColor
	Slider
	Calender
)
