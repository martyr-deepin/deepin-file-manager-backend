PREFIX?=/usr
TARGET_DIR=$(DESTDIR)/$(PREFIX)/lib/deepin-daemon
PKG_NAME?=file-manager-backend
PKG_VERSION?=unknown
binary=deepin-file-manager-backend
BUILD_DIR=$(shell pwd)/build
SRC_DIR=$(BUILD_DIR)/src/pkg.deepin.io/service/

ifndef USE_GCCGO
	ifndef GOLANG_DEBUG
		LDFLAGS = -ldflags '-s -w'
	endif

	GOBUILD = go build ${LDFLAGS}
else
	ifndef GOLANG_DEBUG
		LDFLAGS += -s -w  -Os -O2
	endif
	LDFLAGS += $(shell pkg-config --libs gio-2.0 gtk+-3.0 gdk-3.0 gdk-pixbuf-xlib-2.0 x11 xi libcanberra cairo-ft poppler-glib libmetacity-private librsvg-2.0)
	GOBUILD = go build -compiler gccgo -gccgoflags "${LDFLAGS}"
endif


all: build


prepare:
	if [ ! -d $(SRC_DIR) ]; then \
		mkdir -p $(SRC_DIR); \
		ln -sf $(shell dirname `pwd`)/$(shell basename `pwd`) $(SRC_DIR)/$(PKG_NAME); \
	fi


build: prepare
	env GOPATH="${GOPATH}:${BUILD_DIR}" $(GOBUILD) -o $(binary)

install-mo:
	make -C locale -f Makefile install -e DESTDIR=$(DESTDIR)

do-install: install-mo
	install -Dm 755 -t $(TARGET_DIR) $(binary)
	install -Dm 644 -t $(DESTDIR)/usr/share/glib-2.0/schemas schema/com.deepin.filemanager.gschema.xml
	mkdir -p $(DESTDIR)/usr/share/dbus-1/services
	cp services/* $(DESTDIR)/usr/share/dbus-1/services

install: build do-install

clean:
	rm -rf $(BUILD_DIR)

distclean: clean
	rm -f $(binary)

pot:
	make -C locale -f Makefile -e PACKAGE_NAME=$(PKG_NAME) -e PACKAGE_VERSION=$(PKG_VERSION)

mo:
	make -C locale -f Makefile mo
