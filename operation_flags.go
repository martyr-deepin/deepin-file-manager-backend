package main

import (
	"pkg.deepin.io/lib/dbus"
	"gir/gio-2.0"
	"pkg.deepin.io/service/file-manager-backend/operations"
	"pkg.deepin.io/service/file-manager-backend/delegator"
)

type OperationFlags struct {
	ListJobFlagNone          int32
	ListJobFlagRecusive      int32
	ListJobFlagIncludeHidden int32

	CopyFlagNone             uint32
	CopyFlagNofollowSymlinks uint32
	// CopyFlagOverwrite          uint32
	// CopyFlagBackup             uint32
	// CopyFlagAllMetadata        uint32
	// CopyFlagNoFallbackForMove  uint32
	// CopyFlagTargetDefaultPerms uint32
}

func NewOperationFlags() *OperationFlags {
	flags := new(OperationFlags)

	// list flags
	flags.ListJobFlagNone = int32(operations.ListJobFlagNone)
	flags.ListJobFlagRecusive = int32(operations.ListJobFlagRecusive)
	flags.ListJobFlagIncludeHidden = int32(operations.ListJobFlagIncludeHidden)

	// copy/move flags
	flags.CopyFlagNone = uint32(gio.FileCopyFlagsNone)
	flags.CopyFlagNofollowSymlinks = uint32(gio.FileCopyFlagsNofollowSymlinks)
	// flags.CopyFlagOverwrite = uint32(gio.FileCopyFlagsOverwrite)
	// flags.CopyFlagBackup = uint32(gio.FileCopyFlagsBackup)
	// flags.CopyFlagAllMetadata = uint32(gio.FileCopyFlagsAllMetadata)
	// flags.CopyFlagNoFallbackForMove = uint32(gio.FileCopyFlagsNoFallbackForMove)
	// flags.CopyFlagTargetDefaultPerms = uint32(gio.FileCopyFlagsTargetDefaultPerms)

	return flags
}

func (*OperationFlags) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       delegator.JobDestination,
		ObjectPath: delegator.JobObjectPath,
		Interface:  delegator.JobDestination + ".Flags",
	}
}
