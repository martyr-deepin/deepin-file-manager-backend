package operations

import (
	"deepin-file-manager/dbusproxy"
	"fmt"
	"pkg.linuxdeepin.com/lib/dbus"
)

// type DeletionType int
//
// const (
// 	DELETE DeletionType = iota
// 	TRASH
// 	EMPTY_TRASH
// )
//
// type ConfirmationType int
//
// const (
// 	DefaultConfirmation ConfirmationType = iota
// 	ForceConfirmation
// )

// ResponseCode is a type for the response of UIDelegate.
type ResponseCode int32

// the code for response of UIDelegate.
const (
	ResponseCancel ResponseCode = 1 << iota
	ResponseSkip
	ResponseRetry
	ResponseDelete
	ResponseOverwrite
	ResponseAutoRename // auto rename the conflict file/directory
	ResponseYes
)

// String returns a human readable string for ResponseCode.
func (code ResponseCode) String() string {
	switch code {
	case ResponseCancel:
		return "Cancel"
	case ResponseSkip:
		return "Skip"
	case ResponseRetry:
		return "Retry"
	case ResponseDelete:
		return "Delete"
	case ResponseYes:
		return "Yes"
	case ResponseOverwrite:
		return "Overwrite"
	case ResponseAutoRename:
		return "AutoRename"
	}

	return fmt.Sprintf("Unknow code: %d", int32(code))
}

// Response stores the response relavant information like ResponseCode.
type Response struct {
	code       int32
	applyToAll bool
	userData   string
}

// NewResponse creates a new Response from response code and apply to all.
func NewResponse(code ResponseCode, applyToAll bool) Response {
	return Response{
		code:       int32(code),
		applyToAll: applyToAll,
	}
}

// String returns a human readable string for debuging or something like it.
func (response Response) String() string {
	str := ResponseCode(response.code).String()
	if response.applyToAll {
		str += " to all"
	}

	return str
}

// Code returns response code.
func (response Response) Code() ResponseCode {
	return ResponseCode(response.code)
}

// ApplyToAll returns whether apply to all.
func (response Response) ApplyToAll() bool {
	return response.applyToAll
}

// UserData returns some extra data.
func (response Response) UserData() string {
	return response.userData
}

// UIFlags
const (
	UIFlagsNone int32 = iota << 0
	UIFlagsRetry
	UIFlagsMulti
)

// IUIDelegate is the interface for ui delegate.
type IUIDelegate interface {
	// TODO: using this internally, give a simpler interface, like kio,
	// a.k.a: AskDeleteConfirmation(urls, deleteType, confirmationType))
	// if necessary, ask user to confirm whether to delete or trash files.
	AskDeleteConfirmation(primaryText string, secondaryText string, detailText string) bool

	AskDelete(string, string, string, int32) Response
	AskSkip(primaryText string, secondaryText string, detailText string, uiFlags int32) Response
	AskRetry(primaryText string, secondaryText string, detailText string) Response

	// TODO: decide arguments
	ConflictDialog() Response
}

type _DefaultUIDelegate struct{}

func (*_DefaultUIDelegate) AskRetry(primaryText string, secondaryText string, detailText string) Response {
	return NewResponse(ResponseCancel, false)
}

func (*_DefaultUIDelegate) AskDeleteConfirmation(primaryText string, secondaryText string, detailText string) bool {
	return true
}

func (*_DefaultUIDelegate) AskDelete(primaryText string, secondaryText string, detailText string, flags int32) Response {
	return NewResponse(ResponseSkip, true)
}

func (*_DefaultUIDelegate) AskSkip(primaryText string, secondaryText string, detailText string, flags int32) Response {
	// TODO:
	return NewResponse(ResponseCancel, true)
}

func (*_DefaultUIDelegate) ConflictDialog() Response {
	return NewResponse(ResponseSkip, true)
}

var _deafultUIDelegate = &_DefaultUIDelegate{}

// UIDelegate is a proxy for dbus UIDelegate.
type UIDelegate struct {
	proxy *dbusproxy.DBusProxy
}

// NewUIDelegate creates a new UIDelegate for dbus.
func NewUIDelegate(conn *dbus.Conn, dest string, objPath string, iface string) (IUIDelegate, error) {
	proxy, err := dbusproxy.NewDBusProxy(conn, dest, objPath, iface, 0)
	if err != nil {
		return nil, err
	}
	return &UIDelegate{
		proxy: proxy,
	}, nil
}

func (delegate *UIDelegate) call(name string, args ...interface{}) *dbus.Call {
	return delegate.proxy.Call(name, args...)
}

// AskSkip asks user whether skip this error.
func (delegate *UIDelegate) AskSkip(primaryText string, secondaryText string, detailText string, flags int32) Response {
	var response Response
	retry := flags&UIFlagsRetry != 0
	multi := flags&UIFlagsMulti != 0
	delegate.call("AskSkip", primaryText, secondaryText, detailText, retry, multi).Store(&response)
	return response
}

// AskDelete asks user whether delete.
func (delegate *UIDelegate) AskDelete(primaryText string, secondaryText string, detailText string, flags int32) Response {
	var response Response
	retry := flags&UIFlagsRetry != 0
	multi := flags&UIFlagsMulti != 0
	delegate.call("AskDelete", primaryText, secondaryText, detailText, retry, multi).Store(&response)
	return response
}

// AskDeleteConfirmation asks for the confirm for delete.
func (delegate *UIDelegate) AskDeleteConfirmation(primaryText string, secondaryText string, detailText string) bool {
	confirm := false
	delegate.call("AskDeleteConfirmation", primaryText, secondaryText, detailText).Store(&confirm)
	return confirm
}

// ConflictDialog is used for the conflict situaction like copy.
func (delegate *UIDelegate) ConflictDialog() Response {
	// TODO: impl
	return Response{}
}

// AskRetry asks user whether to retry this operation.
func (delegate *UIDelegate) AskRetry(primaryText string, secondaryText string, detailText string) Response {
	var response Response
	delegate.call("AskRetry", primaryText, secondaryText, detailText).Store(&response)
	return response
}
