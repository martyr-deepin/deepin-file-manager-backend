package desktop

import (
	"pkg.deepin.io/lib/dbus"
)

func showModule(module string) {
	go func() {
		conn, err := dbus.SessionBus()
		if err != nil {
			return
		}

		obj := conn.Object("com.deepin.dde.ControlCenter", "/com/deepin/dde/ControlCenter")
		if obj != nil {
			obj.Call("com.deepin.dde.ControlCenter.ShowModule", 0, module).Store()
		}
	}()
}
