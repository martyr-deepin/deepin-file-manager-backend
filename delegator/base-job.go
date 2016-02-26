/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package delegator

import (
	"fmt"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/service/file-manager-backend/operations"
	"sync/atomic"
)

// IUIDelegate is the interface for UIDelegate, a alias for operations.IUIDelegate.
type IUIDelegate operations.IUIDelegate

// TODO: make a real interface.
type IUndoManager interface{}

const (
	// JobDestination is dbus destination for backend and operations
	JobDestination = "com.deepin.filemanager.Backend.Operations"
	// JobObjectPath is dbus object path for backend and operations.
	JobObjectPath = "/com/deepin/filemanager/Backend/Operations"
)

// _BaseJob causes some problems. NOT using it now.
type _BaseJob struct {
	name     string
	dbusInfo dbus.DBusInfo
}

func (job *_BaseJob) GetDBusInfo() dbus.DBusInfo {
	return job.dbusInfo
}

func newBaseJob(name string, count *uint64) *_BaseJob {
	job := &_BaseJob{
		name: name,
	}
	job.dbusInfo = genDBusInfo(name, count)
	return job
}

func genObjectPath(name string, count uint64) string {
	return fmt.Sprintf("%s/%s%d", JobObjectPath, name, count)
}

func genInterface(name string) string {
	return JobDestination + "." + name
}

func genDBusInfo(name string, count *uint64) dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       JobDestination,
		ObjectPath: genObjectPath(name, atomic.AddUint64(count, 1)),
		Interface:  genInterface(name),
	}
}
