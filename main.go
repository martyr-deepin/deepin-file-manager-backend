package main

// #cgo pkg-config: gtk+-3.0
// #include <gtk/gtk.h>
// void GtkInit() { gtk_init(NULL, NULL); }
import "C"
import (
	"os"
	"pkg.deepin.io/lib"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/glib-2.0"
	"pkg.deepin.io/service/file-manager-backend/clipboard"
	"pkg.deepin.io/service/file-manager-backend/desktop"
	"pkg.deepin.io/service/file-manager-backend/fileinfo"
	. "pkg.deepin.io/service/file-manager-backend/log"
	"pkg.deepin.io/service/file-manager-backend/monitor"
)

type Initializer struct {
	err error
}

func (init *Initializer) Init(fn func() (dbus.DBusObject, error)) *Initializer {
	if init.err != nil {
		return init
	}

	v, err := fn()
	if err != nil {
		init.err = err
		return init
	}

	err = dbus.InstallOnSession(v)
	if err != nil {
		init.err = err
	}

	return init
}

func (init *Initializer) GetError() error {
	return init.err
}

func main() {
	C.GtkInit()
	operationBackend := NewOperationBackend()

	gettext.InitI18n()
	gettext.Textdomain("DFMB")

	info := operationBackend.GetDBusInfo()
	if !lib.UniqueOnSession(info.Dest) {
		Log.Info("already exists a session bus named", info.Dest)
		os.Exit(1)
	}

	initializer := new(Initializer)

	initializer.Init(func() (dbus.DBusObject, error) {
		Log.Info("initialize operation backend...")
		return operationBackend, nil
	}).Init(func() (dbus.DBusObject, error) {
		Log.Info("ok")
		Log.Info("initialize operation flags dbus interface...")
		return NewOperationFlags(), nil
	}).Init(func() (dbus.DBusObject, error) {
		Log.Info("ok")
		Log.Info("initialize monitor manager...")
		return monitor.NewMonitorManager(), nil
	}).Init(func() (dbus.DBusObject, error) {
		Log.Info("ok")
		Log.Info("initialize trash monitor...")
		return monitor.NewTrashMonitor()
	}).Init(func() (dbus.DBusObject, error) {
		Log.Info("ok")
		Log.Info("initialize file info...")
		return fileinfo.NewQueryFileInfoJob(), nil
	}).Init(func() (dbus.DBusObject, error) {
		Log.Info("ok")
		Log.Info("initialize Clipboard...")
		return clipboard.NewClipboard(), nil
	}).Init(func() (dbus.DBusObject, error) {
		Log.Info("ok")
		Log.Info("initialize desktop daemon...")
		return desktop.NewDesktopDaemon()
	})

	if err := initializer.GetError(); err != nil {
		Log.Info("Failed:", err)
		os.Exit(1)
	}

	Log.Info("ok")
	go glib.StartLoop()
	dbus.DealWithUnhandledMessage()
	if err := dbus.Wait(); err != nil {
		Log.Info(err)
		os.Exit(2)
	}
}
