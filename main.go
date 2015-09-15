package main

// #cgo pkg-config: gtk+-3.0
// #include <gtk/gtk.h>
// void GtkInit() { gtk_init(NULL, NULL); }
import "C"
import (
	"os"
	"pkg.deepin.io/lib"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/glib-2.0"
	"pkg.deepin.io/service/file-manager-backend/clipboard"
	"pkg.deepin.io/service/file-manager-backend/desktop"
	"pkg.deepin.io/service/file-manager-backend/fileinfo"
	"pkg.deepin.io/service/file-manager-backend/log"
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

	info := operationBackend.GetDBusInfo()
	if !lib.UniqueOnSession(info.Dest) {
		log.Info("already exists a session bus named", info.Dest)
		os.Exit(1)
	}

	initializer := new(Initializer)

	initializer.Init(func() (dbus.DBusObject, error) {
		log.Info("initialize operation backend...")
		return operationBackend, nil
	}).Init(func() (dbus.DBusObject, error) {
		log.Info("ok")
		log.Info("initialize operation flags dbus interface...")
		return NewOperationFlags(), nil
	}).Init(func() (dbus.DBusObject, error) {
		log.Info("ok")
		log.Info("initialize monitor manager...")
		return monitor.NewMonitorManager(), nil
	}).Init(func() (dbus.DBusObject, error) {
		log.Info("ok")
		log.Info("initialize trash monitor...")
		return monitor.NewTrashMonitor()
	}).Init(func() (dbus.DBusObject, error) {
		log.Info("ok")
		log.Info("initialize file info...")
		return fileinfo.NewQueryFileInfoJob(), nil
	}).Init(func() (dbus.DBusObject, error) {
		log.Info("ok")
		log.Info("initialize Clipboard...")
		return clipboard.NewClipboard(), nil
	}).Init(func() (dbus.DBusObject, error) {
		log.Info("ok")
		log.Info("initialize desktop daemon...")
		return desktop.NewDesktopDaemon()
	})

	if err := initializer.GetError(); err != nil {
		log.Info("Failed:", err)
		os.Exit(1)
	}

	log.Info("ok")
	go glib.StartLoop()
	dbus.DealWithUnhandledMessage()
	if err := dbus.Wait(); err != nil {
		log.Info(err)
		os.Exit(2)
	}
}
