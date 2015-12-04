package main

// #cgo pkg-config: gtk+-3.0
// #include <gtk/gtk.h>
// void GtkInit() { gtk_init(NULL, NULL); }
import "C"
import (
	"flag"
	"os"
	"time"

	"pkg.deepin.io/lib"
	"pkg.deepin.io/lib/dbus"
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
	cpuprof := flag.String("cpuprof", "", "-cpuprof=profile_path")

	flag.Parse()
	if *cpuprof != "" {
		startCPUProfile(*cpuprof)
	}

	startTime := time.Now()

	C.GtkInit()

	Log.Info("initialize i18n...")
	InitI18n()
	Log.Info("initialize i18n done, cost", time.Now().Sub(startTime))

	moduleStartTime := time.Now()
	Log.Info("initialize operation backend...")
	operationBackend := NewOperationBackend()
	info := operationBackend.GetDBusInfo()
	if !lib.UniqueOnSession(info.Dest) {
		Log.Info("already exists a session bus named", info.Dest)
		os.Exit(1)
	}

	initializer := new(Initializer)

	var moduleEndTime time.Time
	initializer.Init(func() (dbus.DBusObject, error) {
		moduleStartTime = time.Now()
		return operationBackend, nil
	}).Init(func() (dbus.DBusObject, error) {
		moduleEndTime = time.Now()
		Log.Info("ok, cost", moduleEndTime.Sub(moduleStartTime))
		Log.Info("initialize operation flags dbus interface...")
		moduleStartTime = moduleEndTime
		return NewOperationFlags(), nil
	}).Init(func() (dbus.DBusObject, error) {
		moduleEndTime = time.Now()
		Log.Info("ok, cost", moduleEndTime.Sub(moduleStartTime))
		Log.Info("initialize monitor manager...")
		moduleStartTime = moduleEndTime
		return monitor.NewMonitorManager(), nil
	}).Init(func() (dbus.DBusObject, error) {
		moduleEndTime = time.Now()
		Log.Info("ok, cost", moduleEndTime.Sub(moduleStartTime))
		Log.Info("initialize trash monitor...")
		moduleStartTime = moduleEndTime
		return monitor.NewTrashMonitor()
	}).Init(func() (dbus.DBusObject, error) {
		moduleEndTime = time.Now()
		Log.Info("ok, cost", moduleEndTime.Sub(moduleStartTime))
		Log.Info("initialize file info...")
		moduleStartTime = moduleEndTime
		return fileinfo.NewQueryFileInfoJob(), nil
	}).Init(func() (dbus.DBusObject, error) {
		moduleEndTime = time.Now()
		Log.Info("ok, cost", moduleEndTime.Sub(moduleStartTime))
		Log.Info("initialize Clipboard...")
		moduleStartTime = moduleEndTime
		return clipboard.NewClipboard(), nil
	}).Init(func() (dbus.DBusObject, error) {
		moduleEndTime = time.Now()
		Log.Info("ok, cost", moduleEndTime.Sub(moduleStartTime))
		Log.Info("initialize desktop daemon...")
		moduleStartTime = moduleEndTime
		return desktop.NewDesktopDaemon()
	})

	ddeSessionRegister()

	if err := initializer.GetError(); err != nil {
		Log.Info("Failed:", err)
		os.Exit(1)
	}

	moduleEndTime = time.Now()
	Log.Info("ok, cost", moduleEndTime.Sub(moduleStartTime))

	go glib.StartLoop()
	dbus.DealWithUnhandledMessage()

	Log.Info("Total cost", time.Now().Sub(startTime))

	if err := dbus.Wait(); err != nil {
		Log.Info(err)
		os.Exit(2)
	}
}
