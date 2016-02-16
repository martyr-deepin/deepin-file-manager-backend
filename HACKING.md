# Deepin File Manager Backend Contribute

## Project layout

### Coding layout

| log/        | the logger used for this project                           |
|-------------|------------------------------------------------------------|
| operations/ | the file operations implimentation                         |
| delegator/  | dbus layer wrapper for file operations                     |
| clipboard/  | dbus layer wrapper for clipboard                           |
| fileinfo/   | dbus layer file information operation, like get theme icon |
| locale/     | contains translation template and files                    |
| monitor/    | dbus layer wrapper for file monitor                        |
| schema/     | setting's schema                                           |
| setting/    | dbus layer setting related operation                       |
| services/   | dbus service files                                         |
| dbusproxy/  | a dbus proxy to avoid generating codes                     |
| desktop/    | the desktop backed                                         |


### Others

The most directory will contains a directory named testdata which stores data for testing.

## Design

The base rule is splitting operation into dbus layer and implementation layer. The dbus
layer is exported API, the implementation layer can be used internally.

### File operations

Because some file operations need some UI response, these operations will need extra three
arguments on dbus layer. These three arguments are **dbus name**, **dbus path**, **dbus interface**,
a delegator will be created on dbus layer according to these arguments, and then pass it to
implementation layer.

In order not to block the UI, the UI response should be pass to backend through **response** signal.

Because operation can be co-existed, a new dbus object path will be created for a new operation
job, and then dbus name, dbus path and dbus interface will be return to caller.


## TODO

[ ] undo manager
[x] mount/unmount job. Use daemon/mounts
[ ] add version to dbus interface or path of operation job, API user should decide which version to use
[ ] writing much more tests
[ ] finish TODOs listed in codes
[x] command line
