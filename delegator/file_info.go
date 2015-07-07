package delegator

import (
	"encoding/json"
	"fmt"
	"pkg.linuxdeepin.com/lib/gio-2.0"
	"pkg.linuxdeepin.com/lib/operations"
	"strings"
)

var _ = fmt.Println

var attributeTypeMap = map[string]gio.FileAttributeType{
	gio.FileAttributeStandardType:                   gio.FileAttributeTypeUint32,
	gio.FileAttributeStandardIsHidden:               gio.FileAttributeTypeBoolean,
	gio.FileAttributeStandardIsBackup:               gio.FileAttributeTypeBoolean,
	gio.FileAttributeStandardIsSymlink:              gio.FileAttributeTypeBoolean,
	gio.FileAttributeStandardIsVirtual:              gio.FileAttributeTypeBoolean,
	gio.FileAttributeStandardName:                   gio.FileAttributeTypeString,
	gio.FileAttributeStandardDisplayName:            gio.FileAttributeTypeString,
	gio.FileAttributeStandardEditName:               gio.FileAttributeTypeString,
	gio.FileAttributeStandardCopyName:               gio.FileAttributeTypeString,
	gio.FileAttributeStandardIcon:                   gio.FileAttributeTypeObject, // GIcon
	gio.FileAttributeStandardSymbolicIcon:           gio.FileAttributeTypeObject, // GIcon
	gio.FileAttributeStandardContentType:            gio.FileAttributeTypeString,
	gio.FileAttributeStandardFastContentType:        gio.FileAttributeTypeString,
	gio.FileAttributeStandardSize:                   gio.FileAttributeTypeUint64,
	gio.FileAttributeStandardAllocatedSize:          gio.FileAttributeTypeUint64,
	gio.FileAttributeStandardSymlinkTarget:          gio.FileAttributeTypeByteString,
	gio.FileAttributeStandardTargetUri:              gio.FileAttributeTypeString,
	gio.FileAttributeStandardSortOrder:              gio.FileAttributeTypeInt32,
	gio.FileAttributeEtagValue:                      gio.FileAttributeTypeString,
	gio.FileAttributeIdFile:                         gio.FileAttributeTypeString,
	gio.FileAttributeIdFilesystem:                   gio.FileAttributeTypeString,
	gio.FileAttributeAccessCanRead:                  gio.FileAttributeTypeBoolean,
	gio.FileAttributeAccessCanWrite:                 gio.FileAttributeTypeBoolean,
	gio.FileAttributeAccessCanExecute:               gio.FileAttributeTypeBoolean,
	gio.FileAttributeAccessCanDelete:                gio.FileAttributeTypeBoolean,
	gio.FileAttributeAccessCanTrash:                 gio.FileAttributeTypeBoolean,
	gio.FileAttributeAccessCanRename:                gio.FileAttributeTypeBoolean,
	gio.FileAttributeMountableCanMount:              gio.FileAttributeTypeBoolean,
	gio.FileAttributeMountableCanUnmount:            gio.FileAttributeTypeBoolean,
	gio.FileAttributeMountableCanEject:              gio.FileAttributeTypeBoolean,
	gio.FileAttributeMountableUnixDevice:            gio.FileAttributeTypeUint32,
	gio.FileAttributeMountableUnixDeviceFile:        gio.FileAttributeTypeString,
	gio.FileAttributeMountableHalUdi:                gio.FileAttributeTypeString,
	gio.FileAttributeMountableCanPoll:               gio.FileAttributeTypeBoolean,
	gio.FileAttributeMountableIsMediaCheckAutomatic: gio.FileAttributeTypeBoolean,
	gio.FileAttributeMountableCanStart:              gio.FileAttributeTypeBoolean,
	gio.FileAttributeMountableCanStartDegraded:      gio.FileAttributeTypeBoolean,
	gio.FileAttributeMountableCanStop:               gio.FileAttributeTypeBoolean,
	gio.FileAttributeMountableStartStopType:         gio.FileAttributeTypeUint32, // GDriveStartStopType
	gio.FileAttributeTimeModified:                   gio.FileAttributeTypeUint64,
	gio.FileAttributeTimeModifiedUsec:               gio.FileAttributeTypeUint32,
	gio.FileAttributeTimeAccess:                     gio.FileAttributeTypeUint64,
	gio.FileAttributeTimeAccessUsec:                 gio.FileAttributeTypeUint32,
	gio.FileAttributeTimeChanged:                    gio.FileAttributeTypeUint64,
	gio.FileAttributeTimeChangedUsec:                gio.FileAttributeTypeUint32,
	gio.FileAttributeTimeCreated:                    gio.FileAttributeTypeUint64,
	gio.FileAttributeTimeCreatedUsec:                gio.FileAttributeTypeUint32,
	gio.FileAttributeUnixDevice:                     gio.FileAttributeTypeUint32,
	gio.FileAttributeUnixInode:                      gio.FileAttributeTypeUint64,
	gio.FileAttributeUnixMode:                       gio.FileAttributeTypeUint32,
	gio.FileAttributeUnixNlink:                      gio.FileAttributeTypeUint32,
	gio.FileAttributeUnixUid:                        gio.FileAttributeTypeUint32,
	gio.FileAttributeUnixGid:                        gio.FileAttributeTypeUint32,
	gio.FileAttributeUnixRdev:                       gio.FileAttributeTypeUint32,
	gio.FileAttributeUnixBlockSize:                  gio.FileAttributeTypeUint32,
	gio.FileAttributeUnixBlocks:                     gio.FileAttributeTypeUint64,
	gio.FileAttributeUnixIsMountpoint:               gio.FileAttributeTypeBoolean,
	gio.FileAttributeDosIsArchive:                   gio.FileAttributeTypeBoolean,
	gio.FileAttributeDosIsSystem:                    gio.FileAttributeTypeBoolean,
	gio.FileAttributeOwnerUser:                      gio.FileAttributeTypeString,
	gio.FileAttributeOwnerUserReal:                  gio.FileAttributeTypeString,
	gio.FileAttributeOwnerGroup:                     gio.FileAttributeTypeString,
	gio.FileAttributeThumbnailPath:                  gio.FileAttributeTypeByteString,
	gio.FileAttributeThumbnailingFailed:             gio.FileAttributeTypeBoolean,
	gio.FileAttributeThumbnailIsValid:               gio.FileAttributeTypeBoolean,
	gio.FileAttributePreviewIcon:                    gio.FileAttributeTypeObject, // GIcon
	gio.FileAttributeFilesystemSize:                 gio.FileAttributeTypeUint64,
	gio.FileAttributeFilesystemFree:                 gio.FileAttributeTypeUint64,
	gio.FileAttributeFilesystemUsed:                 gio.FileAttributeTypeUint64,
	gio.FileAttributeFilesystemType:                 gio.FileAttributeTypeString,
	gio.FileAttributeFilesystemReadonly:             gio.FileAttributeTypeBoolean,
	gio.FileAttributeGvfsBackend:                    gio.FileAttributeTypeString,
	gio.FileAttributeSelinuxContext:                 gio.FileAttributeTypeString,
	gio.FileAttributeTrashItemCount:                 gio.FileAttributeTypeUint32,
	gio.FileAttributeTrashDeletionDate:              gio.FileAttributeTypeString,
	gio.FileAttributeTrashOrigPath:                  gio.FileAttributeTypeByteString,
	gio.FileAttributeFilesystemUsePreview:           gio.FileAttributeTypeUint32, // GFilesystemPreviewType
	gio.FileAttributeStandardDescription:            gio.FileAttributeTypeString,
}

type QueryFileInfoJob struct {
	QueryFlagNofollowSymlinks uint32
	QueryFlagNone             uint32

	FileTypeRegular      uint32
	FileTypeSpecial      uint32
	FileTypeUnknown      uint32
	FileTypeShortcut     uint32
	FileTypeDirectory    uint32
	FileTypeMountable    uint32
	FileTypeSymbolicLink uint32

	DriveStartStopTypePassword  uint32
	DriveStartStopTypeNetwork   uint32
	DriveStartStopTypeUnknown   uint32
	DriveStartStopTypeShutdown  uint32
	DriveStartStopTypeMultidisk uint32

	FilesystemPreviewTypeIfLocal  uint32
	FilesystemPreviewTypeNever    uint32
	FilesystemPreviewTypeIfAlways uint32

	FileAttributeStandardType                   string
	FileAttributeStandardIsHidden               string
	FileAttributeStandardIsBackup               string
	FileAttributeStandardIsSymlink              string
	FileAttributeStandardIsVirtual              string
	FileAttributeStandardName                   string
	FileAttributeStandardDisplayName            string
	FileAttributeStandardEditName               string
	FileAttributeStandardCopyName               string
	FileAttributeStandardIcon                   string
	FileAttributeStandardSymbolicIcon           string
	FileAttributeStandardContentType            string
	FileAttributeStandardFastContentType        string
	FileAttributeStandardSize                   string
	FileAttributeStandardAllocatedSize          string
	FileAttributeStandardSymlinkTarget          string
	FileAttributeStandardTargetUri              string
	FileAttributeStandardSortOrder              string
	FileAttributeEtagValue                      string
	FileAttributeIdFile                         string
	FileAttributeIdFilesystem                   string
	FileAttributeAccessCanRead                  string
	FileAttributeAccessCanWrite                 string
	FileAttributeAccessCanExecute               string
	FileAttributeAccessCanDelete                string
	FileAttributeAccessCanTrash                 string
	FileAttributeAccessCanRename                string
	FileAttributeMountableCanMount              string
	FileAttributeMountableCanUnmount            string
	FileAttributeMountableCanEject              string
	FileAttributeMountableUnixDevice            string
	FileAttributeMountableUnixDeviceFile        string
	FileAttributeMountableHalUdi                string
	FileAttributeMountableCanPoll               string
	FileAttributeMountableIsMediaCheckAutomatic string
	FileAttributeMountableCanStart              string
	FileAttributeMountableCanStartDegraded      string
	FileAttributeMountableCanStop               string
	FileAttributeMountableStartStopType         string
	FileAttributeTimeModified                   string
	FileAttributeTimeModifiedUsec               string
	FileAttributeTimeAccess                     string
	FileAttributeTimeAccessUsec                 string
	FileAttributeTimeChanged                    string
	FileAttributeTimeChangedUsec                string
	FileAttributeTimeCreated                    string
	FileAttributeTimeCreatedUsec                string
	FileAttributeUnixDevice                     string
	FileAttributeUnixInode                      string
	FileAttributeUnixMode                       string
	FileAttributeUnixNlink                      string
	FileAttributeUnixUid                        string
	FileAttributeUnixGid                        string
	FileAttributeUnixRdev                       string
	FileAttributeUnixBlockSize                  string
	FileAttributeUnixBlocks                     string
	FileAttributeUnixIsMountpoint               string
	FileAttributeDosIsArchive                   string
	FileAttributeDosIsSystem                    string
	FileAttributeOwnerUser                      string
	FileAttributeOwnerUserReal                  string
	FileAttributeOwnerGroup                     string
	FileAttributeThumbnailPath                  string
	FileAttributeThumbnailingFailed             string
	FileAttributeThumbnailIsValid               string
	FileAttributePreviewIcon                    string
	FileAttributeFilesystemSize                 string
	FileAttributeFilesystemFree                 string
	FileAttributeFilesystemUsed                 string
	FileAttributeFilesystemType                 string
	FileAttributeFilesystemReadonly             string
	FileAttributeGvfsBackend                    string
	FileAttributeSelinuxContext                 string
	FileAttributeTrashItemCount                 string
	FileAttributeTrashDeletionDate              string
	FileAttributeTrashOrigPath                  string
	FileAttributeFilesystemUsePreview           string
	FileAttributeStandardDescription            string
}

func NewQueryFileInfoJob() *QueryFileInfoJob {
	return &QueryFileInfoJob{
		QueryFlagNofollowSymlinks: uint32(gio.FileQueryInfoFlagsNofollowSymlinks),
		QueryFlagNone:             uint32(gio.FileQueryInfoFlagsNone),

		FileTypeRegular:      uint32(gio.FileTypeRegular),
		FileTypeSpecial:      uint32(gio.FileTypeSpecial),
		FileTypeUnknown:      uint32(gio.FileTypeUnknown),
		FileTypeShortcut:     uint32(gio.FileTypeShortcut),
		FileTypeDirectory:    uint32(gio.FileTypeDirectory),
		FileTypeMountable:    uint32(gio.FileTypeMountable),
		FileTypeSymbolicLink: uint32(gio.FileTypeSymbolicLink),

		DriveStartStopTypePassword:  uint32(gio.DriveStartStopTypePassword),
		DriveStartStopTypeNetwork:   uint32(gio.DriveStartStopTypeNetwork),
		DriveStartStopTypeUnknown:   uint32(gio.DriveStartStopTypeUnknown),
		DriveStartStopTypeShutdown:  uint32(gio.DriveStartStopTypeShutdown),
		DriveStartStopTypeMultidisk: uint32(gio.DriveStartStopTypeMultidisk),

		FilesystemPreviewTypeIfLocal:  uint32(gio.FilesystemPreviewTypeIfLocal),
		FilesystemPreviewTypeNever:    uint32(gio.FilesystemPreviewTypeNever),
		FilesystemPreviewTypeIfAlways: uint32(gio.FilesystemPreviewTypeIfAlways),

		FileAttributeStandardType:                   gio.FileAttributeStandardType,
		FileAttributeStandardIsHidden:               gio.FileAttributeStandardIsHidden,
		FileAttributeStandardIsBackup:               gio.FileAttributeStandardIsBackup,
		FileAttributeStandardIsSymlink:              gio.FileAttributeStandardIsSymlink,
		FileAttributeStandardIsVirtual:              gio.FileAttributeStandardIsVirtual,
		FileAttributeStandardName:                   gio.FileAttributeStandardName,
		FileAttributeStandardDisplayName:            gio.FileAttributeStandardDisplayName,
		FileAttributeStandardEditName:               gio.FileAttributeStandardEditName,
		FileAttributeStandardCopyName:               gio.FileAttributeStandardCopyName,
		FileAttributeStandardIcon:                   gio.FileAttributeStandardIcon,
		FileAttributeStandardSymbolicIcon:           gio.FileAttributeStandardSymbolicIcon,
		FileAttributeStandardContentType:            gio.FileAttributeStandardContentType,
		FileAttributeStandardFastContentType:        gio.FileAttributeStandardFastContentType,
		FileAttributeStandardSize:                   gio.FileAttributeStandardSize,
		FileAttributeStandardAllocatedSize:          gio.FileAttributeStandardAllocatedSize,
		FileAttributeStandardSymlinkTarget:          gio.FileAttributeStandardSymlinkTarget,
		FileAttributeStandardTargetUri:              gio.FileAttributeStandardTargetUri,
		FileAttributeStandardSortOrder:              gio.FileAttributeStandardSortOrder,
		FileAttributeEtagValue:                      gio.FileAttributeEtagValue,
		FileAttributeIdFile:                         gio.FileAttributeIdFile,
		FileAttributeIdFilesystem:                   gio.FileAttributeIdFilesystem,
		FileAttributeAccessCanRead:                  gio.FileAttributeAccessCanRead,
		FileAttributeAccessCanWrite:                 gio.FileAttributeAccessCanWrite,
		FileAttributeAccessCanExecute:               gio.FileAttributeAccessCanExecute,
		FileAttributeAccessCanDelete:                gio.FileAttributeAccessCanDelete,
		FileAttributeAccessCanTrash:                 gio.FileAttributeAccessCanTrash,
		FileAttributeAccessCanRename:                gio.FileAttributeAccessCanRename,
		FileAttributeMountableCanMount:              gio.FileAttributeMountableCanMount,
		FileAttributeMountableCanUnmount:            gio.FileAttributeMountableCanUnmount,
		FileAttributeMountableCanEject:              gio.FileAttributeMountableCanEject,
		FileAttributeMountableUnixDevice:            gio.FileAttributeMountableUnixDevice,
		FileAttributeMountableUnixDeviceFile:        gio.FileAttributeMountableUnixDeviceFile,
		FileAttributeMountableHalUdi:                gio.FileAttributeMountableHalUdi,
		FileAttributeMountableCanPoll:               gio.FileAttributeMountableCanPoll,
		FileAttributeMountableIsMediaCheckAutomatic: gio.FileAttributeMountableIsMediaCheckAutomatic,
		FileAttributeMountableCanStart:              gio.FileAttributeMountableCanStart,
		FileAttributeMountableCanStartDegraded:      gio.FileAttributeMountableCanStartDegraded,
		FileAttributeMountableCanStop:               gio.FileAttributeMountableCanStop,
		FileAttributeMountableStartStopType:         gio.FileAttributeMountableStartStopType,
		FileAttributeTimeModified:                   gio.FileAttributeTimeModified,
		FileAttributeTimeModifiedUsec:               gio.FileAttributeTimeModifiedUsec,
		FileAttributeTimeAccess:                     gio.FileAttributeTimeAccess,
		FileAttributeTimeAccessUsec:                 gio.FileAttributeTimeAccessUsec,
		FileAttributeTimeChanged:                    gio.FileAttributeTimeChanged,
		FileAttributeTimeChangedUsec:                gio.FileAttributeTimeChangedUsec,
		FileAttributeTimeCreated:                    gio.FileAttributeTimeCreated,
		FileAttributeTimeCreatedUsec:                gio.FileAttributeTimeCreatedUsec,
		FileAttributeUnixDevice:                     gio.FileAttributeUnixDevice,
		FileAttributeUnixInode:                      gio.FileAttributeUnixInode,
		FileAttributeUnixMode:                       gio.FileAttributeUnixMode,
		FileAttributeUnixNlink:                      gio.FileAttributeUnixNlink,
		FileAttributeUnixUid:                        gio.FileAttributeUnixUid,
		FileAttributeUnixGid:                        gio.FileAttributeUnixGid,
		FileAttributeUnixRdev:                       gio.FileAttributeUnixRdev,
		FileAttributeUnixBlockSize:                  gio.FileAttributeUnixBlockSize,
		FileAttributeUnixBlocks:                     gio.FileAttributeUnixBlocks,
		FileAttributeUnixIsMountpoint:               gio.FileAttributeUnixIsMountpoint,
		FileAttributeDosIsArchive:                   gio.FileAttributeDosIsArchive,
		FileAttributeDosIsSystem:                    gio.FileAttributeDosIsSystem,
		FileAttributeOwnerUser:                      gio.FileAttributeOwnerUser,
		FileAttributeOwnerUserReal:                  gio.FileAttributeOwnerUserReal,
		FileAttributeOwnerGroup:                     gio.FileAttributeOwnerGroup,
		FileAttributeThumbnailPath:                  gio.FileAttributeThumbnailPath,
		FileAttributeThumbnailingFailed:             gio.FileAttributeThumbnailingFailed,
		FileAttributeThumbnailIsValid:               gio.FileAttributeThumbnailIsValid,
		FileAttributePreviewIcon:                    gio.FileAttributePreviewIcon,
		FileAttributeFilesystemSize:                 gio.FileAttributeFilesystemSize,
		FileAttributeFilesystemFree:                 gio.FileAttributeFilesystemFree,
		FileAttributeFilesystemUsed:                 gio.FileAttributeFilesystemUsed,
		FileAttributeFilesystemType:                 gio.FileAttributeFilesystemType,
		FileAttributeFilesystemReadonly:             gio.FileAttributeFilesystemReadonly,
		FileAttributeGvfsBackend:                    gio.FileAttributeGvfsBackend,
		FileAttributeSelinuxContext:                 gio.FileAttributeSelinuxContext,
		FileAttributeTrashItemCount:                 gio.FileAttributeTrashItemCount,
		FileAttributeTrashDeletionDate:              gio.FileAttributeTrashDeletionDate,
		FileAttributeTrashOrigPath:                  gio.FileAttributeTrashOrigPath,
		FileAttributeFilesystemUsePreview:           gio.FileAttributeFilesystemUsePreview,
		FileAttributeStandardDescription:            gio.FileAttributeStandardDescription,
	}
}

func (job *QueryFileInfoJob) QueryInfo(arg string, attributes string, flags uint32) string {
	file := gio.FileNewForCommandlineArg(arg)
	if file == nil {
		return ""
	}
	defer file.Unref()

	info, err := file.QueryInfo(attributes, gio.FileQueryInfoFlags(flags), nil)
	if err != nil {
		// fmt.Println(err)
		return ""
	}
	defer info.Unref()

	infoJsonMap := map[string]interface{}{}
	queriedAttributes := info.ListAttributes("\x00") // pass NULL to get all attributes.
	for _, attribute := range queriedAttributes {
		attributeType := attributeTypeMap[attribute]
		switch attributeType {
		case gio.FileAttributeTypeBoolean:
			infoJsonMap[attribute] = info.GetAttributeBoolean(attribute)
		case gio.FileAttributeTypeByteString:
			infoJsonMap[attribute] = info.GetAttributeByteString(attribute)
		case gio.FileAttributeTypeInt32:
			infoJsonMap[attribute] = info.GetAttributeInt32(attribute)
		case gio.FileAttributeTypeObject:
			var icon string
			filePath := file.GetPath()

			isApp := strings.HasSuffix(filePath, ".desktop")
			if isApp {
				app := gio.NewDesktopAppInfoFromFilename(filePath)
				defer app.Unref()
				icon = operations.GetIconForApp(app.GetIcon(), 48)
			} else {
				icon = operations.GetIconForFile(info.GetIcon(), 48)
			}

			infoJsonMap[attribute] = icon
		case gio.FileAttributeTypeString:
			infoJsonMap[attribute] = info.GetAttributeString(attribute)
		case gio.FileAttributeTypeUint32:
			infoJsonMap[attribute] = info.GetAttributeUint32(attribute)
		case gio.FileAttributeTypeUint64:
			infoJsonMap[attribute] = info.GetAttributeUint64(attribute)
		}
	}
	infoJsonByteStr, err := json.Marshal(infoJsonMap)
	if err != nil {
		// fmt.Println(err)
		return ""
	}

	return string(infoJsonByteStr)
}

func (job *QueryFileInfoJob) IsNativeFile(arg string) bool {
	f := gio.FileNewForCommandlineArg(arg)
	if f == nil {
		return true // FIXME: is true ok?
	}
	defer f.Unref()

	return f.IsNative()
}
