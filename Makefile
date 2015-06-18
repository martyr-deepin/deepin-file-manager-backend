PREFIX?=/usr
TARGET_DIR=$(PREFIX)/lib/deepin-daemon
PKG_NAME=deepin-file-manager
TARGET=deepin-file-manager-backend
GOPATH?=$$(GOPATH)
PWD=$(shell pwd)
BUILD_DIR=${PWD}/build

GOBUILD=go build


all: build


prepare:
	@if [ ! -d $(BUILD_DIR)/src ]; then \
		mkdir -p $(BUILD_DIR)/src; \
		ln -sf $(PWD)/../$(PKG_NAME) $(BUILD_DIR)/src; \
	fi


build: prepare
	env GOPATH=$(GOPATH):$(BUILD_DIR) $(GOBUILD) -o $(TARGET)


install: build
	install -m 755 -t $(TARGET_DIR) $(TARGET)
	install -m 755 -t /usr/share/glib-2.0/schemas $(PWD)/schema/com.deepin.filemanager.gschema.xml


clean:
	rm -rf $(BUILD_DIR)


distclean: clean
	rm $(TARGET) deepin-file-manager
