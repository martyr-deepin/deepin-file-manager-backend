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
	"pkg.deepin.io/lib/initializer"
)

func NewDesktopDaemon() (*Application, error) {
	var app *Application

	initializer := initializer.NewInitializer()
	initializer.InitOnSessionBus(func(v interface{}) (interface{}, error) {
		return NewSettings()
	}).InitOnSessionBus(func(v interface{}) (interface{}, error) {
		app = NewApplication(v.(*Settings))
		return app, nil
	})

	return app, initializer.GetError()
}
