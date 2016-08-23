/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package main

// #cgo pkg-config: gtk+-3.0
// #include <gtk/gtk.h>
// void GtkInit() { gtk_init(NULL, NULL); }
import "C"
import (
	"os"

	"gir/glib-2.0"
	"pkg.deepin.io/dde/api/session"
	"pkg.deepin.io/lib"
	dapp "pkg.deepin.io/lib/app"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/initializer/v2"
	"pkg.deepin.io/lib/profile"
	"pkg.deepin.io/lib/proxy"

	"pkg.deepin.io/service/file-manager-backend/clipboard"
	"pkg.deepin.io/service/file-manager-backend/desktop"
	"pkg.deepin.io/service/file-manager-backend/fileinfo"
	. "pkg.deepin.io/service/file-manager-backend/log"
	"pkg.deepin.io/service/file-manager-backend/monitor"
)

func main() {
	// change current working directory to desktop.
	os.Chdir(desktop.GetDesktopDir())

	timer := profile.NewTimer()

	app := dapp.New("deepin-file-manager-backend", "the backend of deepin file manager", "version "+__VERSION__)
	app.ParseCommandLine(os.Args[1:])
	Log.SetLogLevel(app.LogLevel())
	app.StartProfile()

	Log.Debug("Parse command line...ok, cost", timer.Elapsed())

	proxy.SetupProxy()

	Log.Info("initialize operation backend...")
	operationBackend := NewOperationBackend()
	info := operationBackend.GetDBusInfo()
	if !lib.UniqueOnSession(info.Dest) {
		Log.Info("already exists a session bus named", info.Dest)
		os.Exit(1)
	}

	err := initializer.DoWithSessionBus(func() (dbus.DBusObject, error) {
		Log.Info("ok, cost", timer.Elapsed())

		return operationBackend, nil
	}).Do(func() error {
		Log.Info("initialize gtk...")
		C.GtkInit()
		Log.Info("ok, cost", timer.Elapsed())

		Log.Info("initialize i18n...")
		InitI18n()
		Log.Info("ok, cost", timer.Elapsed())

		return nil
	}).DoWithSessionBus(func() (dbus.DBusObject, error) {
		Log.Info("initialize operation flags dbus interface...")
		return NewOperationFlags(), nil
	}).DoWithSessionBus(func() (dbus.DBusObject, error) {
		Log.Info("ok, cost", timer.Elapsed())

		Log.Info("initialize monitor manager...")
		return monitor.NewMonitorManager(), nil
	}).DoWithSessionBus(func() (dbus.DBusObject, error) {
		Log.Info("ok, cost", timer.Elapsed())

		Log.Info("initialize trash monitor...")
		return monitor.NewTrashMonitor()
	}).DoWithSessionBus(func() (dbus.DBusObject, error) {
		Log.Info("ok, cost", timer.Elapsed())

		Log.Info("initialize file info...")
		return fileinfo.NewQueryFileInfoJob(), nil
	}).DoWithSessionBus(func() (dbus.DBusObject, error) {
		Log.Info("ok, cost", timer.Elapsed())

		Log.Info("initialize Clipboard...")
		return clipboard.NewClipboard(), nil
	}).DoWithSessionBus(func() (dbus.DBusObject, error) {
		Log.Info("ok, cost", timer.Elapsed())

		Log.Info("initialize desktop daemon...")
		return desktop.NewDesktopDaemon()
	}).GetError()

	Log.Debug("register session...")
	session.Register()
	Log.Info("ok, cost", timer.Elapsed())

	if err != nil {
		Log.Info("Failed:", err)
		os.Exit(1)
	}

	go glib.StartLoop()
	dbus.DealWithUnhandledMessage()

	Log.Info("Total cost", timer.TotalCost())

	if err := dbus.Wait(); err != nil {
		Log.Info(err)
		os.Exit(2)
	}
}
