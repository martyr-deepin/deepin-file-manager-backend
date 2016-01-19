# Deepin File Manager Backend

**Description**: Deepin File Manager Backend includes the most UI non-related operations and informations, like settings, file info.

**Tags**: backend


## Dependencies

### Build dependencies

- gio-2.0
- gtk+-3.0
- gdk-3.0
- gdk-pixbuf-xlib-2.0
- x11
- xi
- libcanberra
- cairo-ft
- poppler-glib
- libdeepin-metacity-private
- librsvg-2.0


### Runtime dependencies

- [deepin-nautilus-properties](https://github.com/linuxdeepin/deepin-nautilus-properties)


## Installation

Build:
```
$ make GOPATH=/usr/share/gocode
```

Or, build through gccgo
```
$ make GOPATH=/usr/share/gocode USE_GCCGO=1
```

Install:
```
sudo make install
```


## Usage

This program will be started automatically if needed, but you can still start it manually.

    $ $(PREFIX)/lib/deepin-daemon/deepin-file-manager-backend

The default $(PREFIX) equals to /usr.


## Getting help

Any usage issues can ask for help via

- [Gitter](https://gitter.im/orgs/linuxdeepin/rooms)
- [IRC channel](https://webchat.freenode.net/?channels=deepin)
- [Forum](https://bbs.deepin.org)
- [Wiki](http://wiki.deepin.org/)


## Getting involved

We encourage you to report issuses and contribute changes.

- [Contribution guide for users](http://wiki.deepin.org/index.php?title=Contribution\_Guidelinex\_for\_Users)
- [Contribution guide for developers](http://wiki.deepin.org/index.php?title=-Contribution\_Guidelines\_for\_Developers)


## License

Deepin file manager backend is licensed under [GPLv3](LICENSE).
