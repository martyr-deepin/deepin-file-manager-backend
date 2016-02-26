/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package monitor

import (
	"pkg.deepin.io/lib/dbus"
	"gir/gio-2.0"
)

type TrashMonitor struct {
	trash   *gio.File
	monitor *gio.FileMonitor

	ItemCountChanged func(uint32)
}

func NewTrashMonitor() (*TrashMonitor, error) {
	trashMonitor := new(TrashMonitor)
	trash := gio.FileNewForUri("trash://")
	trashMonitor.trash = trash
	monitor, err := trash.MonitorDirectory(gio.FileMonitorFlagsNone, nil)
	if err != nil {
		return nil, err
	}

	trashMonitor.monitor = monitor
	monitor.Connect("changed", func(*gio.FileMonitor, *gio.File, *gio.File, gio.FileMonitorEvent) {
		// TODO: seperate trash monitor and dbus into two parts.
		dbus.Emit(trashMonitor, "ItemCountChanged", trashMonitor.ItemCount())
	})

	return trashMonitor, nil
}

func getTrashItemCount(trash *gio.File) uint32 {
	info, err := trash.QueryInfo(gio.FileAttributeTrashItemCount, 0, nil)
	if err != nil {
		// logger.Warning()
		return 0
	}

	return info.GetAttributeUint32(gio.FileAttributeTrashItemCount)
}

func (monitor *TrashMonitor) ItemCount() uint32 {
	return getTrashItemCount(monitor.trash)
}
