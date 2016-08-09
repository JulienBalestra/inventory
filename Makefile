SRC_DIR=src/main
BIN_DIR=bin/
TARGET=inventory
PACKAGE=main
CC=go

SRCS = \
machines.go \
main.go \
interfaces.go \
requests.go

default: $(TARGET)

$(addprefix $(SRC_DIR), %.go):
	@

$(TARGET): $(addprefix $(SRC_DIR), %.go)
	CGO_ENABLED=0 GOOS=linux $(CC) build -ldflags '-w' -o $(TARGET) $(PACKAGE)

clean:
	rm -v $(TARGET) || true

re: clean $(TARGET)