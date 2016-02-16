/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package delegator

// TODO: make a dbus undomanager.
type UndoManager struct {
	UndoAvailable        func(bool)
	UndoJobFinished      func()
	UndoTextChanged      func(string)
	JobRecordingStarted  func(int32) // op
	JobRecordingFinished func(int32) //op
}

func (*UndoManager) Undo() {
}
