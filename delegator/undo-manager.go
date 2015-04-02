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
