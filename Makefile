PREFIX?=/usr
TARGET_DIR=$(DESTDIR)/$(PREFIX)/lib/deepin-daemon
PKG_NAME=file-manager-backend
binary=deepin-file-manager-backend
BUILD_DIR=$(shell pwd)/build
SRC_DIR=$(BUILD_DIR)/src/pkg.deepin.io/service/

GOBUILD=go build


all: build


prepare:
	if [ ! -d $(SRC_DIR) ]; then \
		mkdir -p $(SRC_DIR); \
		ln -sf $(shell dirname `pwd`)/$(shell basename `pwd`) $(SRC_DIR)/$(PKG_NAME); \
	fi


build: prepare
	env GOPATH="${GOPATH}:${BUILD_DIR}" $(GOBUILD) -o $(binary)


install: build
	install -Dm 755 -t $(TARGET_DIR) $(binary)
	install -Dm 644 -t $(DESTDIR)/usr/share/glib-2.0/schemas schema/com.deepin.filemanager.gschema.xml

clean:
	rm -rf $(BUILD_DIR)

distclean: clean
	rm -f $(binary) deepin-file-manager
