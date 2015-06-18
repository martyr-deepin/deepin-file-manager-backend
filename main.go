package main

// #cgo pkg-config: gtk+-3.0
// #include <gtk/gtk.h>
// void GtkInit() { gtk_init(NULL, NULL); }
import "C"
import (
	"deepin-file-manager/delegator"
	"deepin-file-manager/monitor"
	"fmt"
	"log"
	"os"
	"pkg.linuxdeepin.com/lib"
	"pkg.linuxdeepin.com/lib/dbus"
	"pkg.linuxdeepin.com/lib/glib-2.0"
)

type Initializer struct {
	err error
}

func (init *Initializer) Init(fn func() (dbus.DBusObject, error)) *Initializer {
	if init.err != nil {
		return init
	}

	v, err := fn()
	dbusInfo := v.GetDBusInfo()
	fmt.Println(dbusInfo.Dest, dbusInfo.ObjectPath, dbusInfo.Interface, err)
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
		return delegator.NewQueryFileInfoJob(), nil
	})

	if err := initializer.GetError(); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	glib.StartLoop()
	dbus.DealWithUnhandledMessage()
	if err := dbus.Wait(); err != nil {
		log.Println(err)
		os.Exit(2)
	}
}
