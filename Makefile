SRC_DIR=src/main
BIN_DIR=bin/
TARGET=inventory
ASSETFS=bindata_assetfs.go
PACKAGE=main
CC=go

UI=ui

default: $(TARGET)

js:
	make -C $(UI)

assetfs: js
	@go-bindata-assetfs $(UI)/...
	@mv -v $(ASSETFS) $(SRC_DIR)

$(TARGET): assetfs
	CGO_ENABLED=0 GOOS=linux $(CC) build -ldflags '-w' -o $(TARGET) $(PACKAGE)

clean:
	rm -v $(TARGET) || true

fclean: clean
	make -C $(UI) fclean

deploy: $(TARGET)
	swift --insecure upload $(TARGET) $(TARGET)

re: clean $(TARGET)