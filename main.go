package main

// #cgo pkg-config: gtk+-3.0
// #include <gtk/gtk.h>
// void GtkInit() { gtk_init(NULL, NULL); }
import "C"
import (
	"pkg.deepin.io/service/file-manager-backend/clipboard"
	"pkg.deepin.io/service/file-manager-backend/desktop"
	"pkg.deepin.io/service/file-manager-backend/fileinfo"
	"pkg.deepin.io/service/file-manager-backend/monitor"
	// "fmt"
	"log"
	"os"
	"pkg.deepin.io/lib"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/glib-2.0"
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
		log.Println("already exists a session bus named", info.Dest)
		os.Exit(1)
	}

	initializer := new(Initializer)

	initializer.Init(func() (dbus.DBusObject, error) {
		return operationBackend, nil
	}).Init(func() (dbus.DBusObject, error) {
		return NewOperationFlags(), nil
	}).Init(func() (dbus.DBusObject, error) {
		return monitor.NewMonitorManager(), nil
	}).Init(func() (dbus.DBusObject, error) {
		return monitor.NewTrashMonitor()
	}).Init(func() (dbus.DBusObject, error) {
		return fileinfo.NewQueryFileInfoJob(), nil
	}).Init(func() (dbus.DBusObject, error) {
		return clipboard.NewClipboard(), nil
	}).Init(func() (dbus.DBusObject, error) {
		return desktop.NewDesktopDaemon()
	})

	if err := initializer.GetError(); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	go glib.StartLoop()
	dbus.DealWithUnhandledMessage()
	if err := dbus.Wait(); err != nil {
		log.Println(err)
		os.Exit(2)
	}
}
