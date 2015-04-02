package config

import (
	"os"
)

type EnumDebugState int

const (
	DebugStateEnvironment EnumDebugState = iota
	DebugStateForceDebug
	DebugStateForceNotDebug
)

var _DebugState EnumDebugState

func SetDebugState(state EnumDebugState) {
	_DebugState = state
}

func Debug() bool {
	switch _DebugState {
	case DebugStateForceDebug:
		return true
	case DebugStateForceNotDebug:
		return false
	case DebugStateEnvironment:
	default:
	}

	return os.Getenv("DFMB_DEBUG") != ""
}
