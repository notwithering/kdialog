package kdialog

type Form uint8

const (
	FormYesNo Form = iota
	FormYesNoCancel
	FormWarningYesNo
	FormWarningContinueCancel
	FormWarningYesNoCancel
	FormSorry
	FormError
	FormMsgBox
	FormInputBox
	FormImgBox
	FormImgInputBox
	FormPassword
	FormNewPassword
	FormTextBox
	FormTextInputBox
	FormComboBox
	FormMenu
	FormChecklist
	FormRadiolist
	FormPassivePopup
	FormOpenFile
	FormSaveFile
	FormOpenExistingDirectory
	FormOpenIcon
	FormProgressBar
	FormPickColor
	FormSlider
	FormCalendar
)
