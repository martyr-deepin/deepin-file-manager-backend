package delegator

import (
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/service/file-manager-backend/dbusproxy"
	. "pkg.deepin.io/service/file-manager-backend/log"
	"pkg.deepin.io/service/file-manager-backend/operations"
	"sync"
)

// Response will be used to store dbus result,
// because private fields cannot be accessed.
type Response struct {
	Code       int32
	ApplyToAll bool
	UserData   string
}

func toResponse(r Response) operations.Response {
	return operations.NewResponseWithUserData(operations.ResponseCode(r.Code), r.ApplyToAll, r.UserData)
}

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
	var response Response
	retry := flags&operations.UIFlagsRetry != 0
	multi := flags&operations.UIFlagsMulti != 0
	response = delegate.getResponseFor("AskSkip", primaryText, secondaryText, detailText, retry, multi)
	return toResponse(response)
}

// AskDelete asks user whether delete.
func (delegate *UIDelegate) AskDelete(primaryText string, secondaryText string, detailText string, flags operations.UIFlags) operations.Response {
	var response Response
	retry := flags&operations.UIFlagsRetry != 0
	multi := flags&operations.UIFlagsMulti != 0
	response = delegate.getResponseFor("AskDelete", primaryText, secondaryText, detailText, retry, multi)
	return toResponse(response)
}

// AskDeleteConfirmation asks for the confirm for delete.
func (delegate *UIDelegate) AskDeleteConfirmation(primaryText string, secondaryText string, detailText string) bool {
	confirm := false
	delegate.call("AskDeleteConfirmation", primaryText, secondaryText, detailText).Store(&confirm)
	return confirm
}

// ConflictDialog is used for the conflict situaction like copy.
func (delegate *UIDelegate) ConflictDialog() operations.Response {
	var response Response
	response = delegate.getResponseFor("ConflictDialog")
	return toResponse(response)
}

// AskRetry asks user whether to retry this operation.
func (delegate *UIDelegate) AskRetry(primaryText string, secondaryText string, detailText string) operations.Response {
	var response Response
	response = delegate.getResponseFor("AskRetry", primaryText, secondaryText, detailText)
	return toResponse(response)
}

func (delegate *UIDelegate) getResponseFor(name string, args ...interface{}) Response {
	var response Response
	var wg sync.WaitGroup
	wg.Add(1)
	var unsubscriber func()
	unsubscriber = delegate.proxy.Subscribe("response", func(value []interface{}) {
		response.Code = value[0].(int32)
		response.ApplyToAll = value[1].(bool)
		response.UserData = value[2].(string)
		wg.Done()
		unsubscriber()
	})
	go delegate.proxy.Call(name, args...)
	wg.Wait()
	Log.Debug("response for", name, response)
	return response
}
