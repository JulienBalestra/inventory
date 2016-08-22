SRC_DIR=src/main
BIN_DIR=bin/
TARGET=inventory
ASSETFS=bindata_assetfs.go
PACKAGE=main
CC=go

PROXY_EXISTS=$(shell [ $PROXY ] && echo 1 || echo 0 )

ifeq ($(PROXY_EXISTS), 1)
	GET=HTTP_PROXY=$(PROXY) go get
else
	GET=go get
endif

DEPS=src/github.com

UI=ui

default: $(TARGET)

$(DEPS):
	@GOPATH=$(shell pwd) $(GET) -u github.com/jteeuwen/go-bindata/...
	@GOPATH=$(shell pwd) $(GET) -u github.com/elazarl/go-bindata-assetfs/...

js:
	make -C $(UI)

assetfs: $(DEPS) js
	GOPATH=$(shell pwd) PATH=$(PATH):$(shell echo "./bin") go-bindata-assetfs $(UI)/...
	@mv -v $(ASSETFS) $(SRC_DIR)

$(TARGET): $(DEPS)
	GOPATH=$(shell pwd) $(CC) build -o $(TARGET) $(PACKAGE)

static: assetfs
	GOPATH=$(shell pwd) CGO_ENABLED=0 GOOS=linux $(CC) build -ldflags '-w' -o $(TARGET) $(PACKAGE)

clean:
	rm -v $(TARGET) || true

fclean: clean
	make -C $(UI) fclean
	rm -Rf $(DEPS) $(BIN_DIR) pkg

run: clean $(TARGET)
	./$(TARGET)

deploy: static
	swift --insecure upload $(TARGET) $(TARGET)

re: clean $(TARGET)