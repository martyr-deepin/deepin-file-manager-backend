/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

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
