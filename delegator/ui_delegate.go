package delegator

import (
	"pkg.deepin.io/service/file-manager-backend/dbusproxy"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/operations"
)

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
func (delegate *UIDelegate) AskSkip(primaryText string, secondaryText string, detailText string, flags operations.UIFlags) operations.Response {
	var response operations.Response
	retry := flags&operations.UIFlagsRetry != 0
	multi := flags&operations.UIFlagsMulti != 0
	delegate.call("AskSkip", primaryText, secondaryText, detailText, retry, multi).Store(&response)
	return response
}

// AskDelete asks user whether delete.
func (delegate *UIDelegate) AskDelete(primaryText string, secondaryText string, detailText string, flags operations.UIFlags) operations.Response {
	var response operations.Response
	retry := flags&operations.UIFlagsRetry != 0
	multi := flags&operations.UIFlagsMulti != 0
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
func (delegate *UIDelegate) ConflictDialog() operations.Response {
	// TODO: impl
	return operations.Response{}
}

// AskRetry asks user whether to retry this operation.
func (delegate *UIDelegate) AskRetry(primaryText string, secondaryText string, detailText string) operations.Response {
	var response operations.Response
	delegate.call("AskRetry", primaryText, secondaryText, detailText).Store(&response)
	return response
}
