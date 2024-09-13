package kdialog

type Form int

const (
	YesNo Form = iota
	YesNoCancel
	WarningYesNo
	WarningContinueCancel
	WarningYesNoCancel
	Sorry
	DetailedSorry
	Error
	DetailedError
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
	Progress
	PickColor
	Slider
	Calender
)

type DialogBox struct {
	Form Form

	Title string

	Text        string
	InitialText string
	Details     string
	FilePath    string

	Timeout int

	StartDir   string
	FileFilter string
	Group      string
	Context    string

	Default string

	Multiple bool

	Items  []string
	Checks []bool

	DontAgain string
	Geometry  string

	Minimum  int
	Maximum  int
	Interval int

	Ok       string
	Yes      string
	No       string
	Cancel   string
	Continue string
}
