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
