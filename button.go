package kdialog

type Button int8

const (
	ButtonUndefined Button = iota - 1
	ButtonOk
	ButtonYes
	ButtonNo
	ButtonCancel
	ButtonContinue
)
