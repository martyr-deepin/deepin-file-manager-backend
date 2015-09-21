package clipboard

import (
	"pkg.deepin.io/lib/dbus"
	. "pkg.deepin.io/service/file-manager-backend/log"
	"pkg.deepin.io/service/file-manager-backend/operations"
)

type Clipboard struct {
	RequestPaste func(string, []string, string)
}

func NewClipboard() *Clipboard {
	c := &Clipboard{}
	return c
}

func (c *Clipboard) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       "com.deepin.filemanager.Backend.Clipboard",
		ObjectPath: "/com/deepin/filemanager/Backend/Clipboard",
		Interface:  "com.deepin.filemanager.Backend.Clipboard",
	}
}

func (c *Clipboard) EmitPaste(file string) {
	contents := operations.GetClipboardContents()
	if len(contents) < 2 {
		Log.Warning("invalid content or empty content in clipboard")
		return
	}

	op := contents[0]
	files := contents[1:]

	switch op {
	case operations.OpCut:
		fallthrough
	case operations.OpCopy:
		dbus.Emit(c, "RequestPaste", op, files, file)
	default:
		Log.Warning("not valid operation")
	}
}

func (c *Clipboard) CutToClipboard(files []string) {
	operations.CutToClipboard(files)
}

func (c *Clipboard) CopyToClipboard(files []string) {
	operations.CopyToClipboard(files)
}

func (c *Clipboard) ClearClipboard() {
	operations.ClearClipboard()
}

// func (c *Clipboard) GetContent() []string {
// 	return operations.GetClipboardContents()
// }
