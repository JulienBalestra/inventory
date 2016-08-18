SRC_DIR=src/main
BIN_DIR=bin/
TARGET=inventory
PACKAGE=main
CC=go

SRCS = \
machines.go \
main.go \
interfaces.go \
requests.go \
probe.go \
common.go

default: $(TARGET)

$(addprefix $(SRC_DIR), %.go):
	@

assetfs:
	rm -v $(SRC_DIR)/bindata_assetfs.go
	go-bindata-assetfs ./ui/...
	mv -v bindata_assetfs.go $(SRC_DIR)

$(TARGET): $(addprefix $(SRC_DIR), %.go) assetfs
	CGO_ENABLED=0 GOOS=linux $(CC) build -ldflags '-w' -o $(TARGET) $(PACKAGE)

clean:
	rm -v $(TARGET) || true

deploy: $(TARGET)
	swift --insecure upload $(TARGET) $(TARGET)

re: clean $(TARGET)