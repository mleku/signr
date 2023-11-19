package signr

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"reflect"
	"unsafe"
)

const (
	PasswordEntryViaTTY = iota
	// other entry types (eg, GUIs) can be added here.
)

func (s *Signr) PasswordEntry(prompt string, entryType int) (pass []byte,
	err error) {

	switch entryType {
	case PasswordEntryViaTTY:
		app := tview.NewApplication()
		inputField := tview.NewInputField().
			SetLabel(prompt).
			SetMaskCharacter('*').
			// SetPlaceholder("E.g. 1234").
			SetFieldWidth(32).
			SetDoneFunc(func(key tcell.Key) {
				app.Stop()
			})
		if err := app.SetRoot(inputField,
			true).EnableMouse(true).Run(); err != nil {
			panic(err)
		}
		passString := inputField.GetText()
		pass = []byte(passString)
		WipeString(&passString)

	default:
		s.Err("password entry type %d not implemented\n", entryType)
	}
	return
}

func WipeString(str *string) {
	stringHeader := (*reflect.StringHeader)(unsafe.Pointer(str))
	data, length := stringHeader.Data, uintptr(stringHeader.Len)
	for i := uintptr(0); i < length; i++ {
		*(*byte)(unsafe.Pointer(data + i)) = ' '
	}
}
