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
	"gir/gio-2.0"
	"pkg.deepin.io/lib/dbus"
	// . "pkg.deepin.io/service/file-manager-backend/log"
	"pkg.deepin.io/service/file-manager-backend/operations"
)

var (
	_MoveJobCount uint64
)

// MoveJob exports to dbus.
type MoveJob struct {
	dbusInfo dbus.DBusInfo

	uris []string
	op   *operations.CopyMoveJob

	Done             func(string)
	TotalAmount      func(int64, uint16)
	ProcessedAmount  func(int64, uint16)
	ProcessedPercent func(int64)
	Moving           func(string)
	Aborted          func()
}

// GetDBusInfo returns dbus information.
func (job *MoveJob) GetDBusInfo() dbus.DBusInfo {
	return job.dbusInfo
}

func (job *MoveJob) listenSignals() {
	job.op.ListenTotalAmount(func(amount int64, unit operations.AmountUnit) {
		dbus.Emit(job, "TotalAmount", amount, uint16(unit))
	})
	job.op.ListenProcessedAmount(func(amount int64, unit operations.AmountUnit) {
		dbus.Emit(job, "ProcessedAmount", amount, uint16(unit))
	})
	job.op.ListenPercent(func(percent int64) {
		dbus.Emit(job, "ProcessedPercent", percent)
	})
	job.op.ListenMoving(func(srcURL string) {
		dbus.Emit(job, "Moving", srcURL)
	})
	job.op.ListenCopyingMovingDone(func(srcURL string, destURL string) {
		// dbus.Emit(job, "MovingDone", srcURL, destURL)
	})
	job.op.ListenCreatingDir(func(dirURL string) {
		// dbus.Emit(job, "CreatingDir", dirURL)
	})
	job.op.ListenDone(func(err error) {
		defer dbus.UnInstallObject(job)
		errMsg := ""
		if errMsg != "" {
			errMsg = err.Error()
		}
		dbus.Emit(job, "Done", errMsg)
	})
	job.op.ListenAborted(func() {
		defer dbus.UnInstallObject(job)
		dbus.Emit(job, "Aborted")
	})
}

// Execute move job.
func (job *MoveJob) Execute() {
	job.listenSignals()
	job.op.Execute()
}

// Abort the job.
func (job *MoveJob) Abort() {
	job.op.Abort()
	dbus.Emit(job, "Aborted")
}

// NewMoveJob creates a new move job for dbus.
func NewMoveJob(srcUrls []string, destDirURL string, targetName string, flags uint32, uiDelegate IUIDelegate) *MoveJob {
	job := &MoveJob{
		dbusInfo: genDBusInfo("MoveJob", &_MoveJobCount),
		op:       operations.NewMoveJob(srcUrls, destDirURL, targetName, gio.FileCopyFlags(flags), uiDelegate),
	}
	return job
}
