
PKG_LIST := `go list ./... | grep -v "servercomm" | grep -v "testmsg"`
lint:
	@golint -set_exit_status $(PKG_LIST)

test:
	@go test -short -timeout 10s $(PKG_LIST)

all: msg
msg:
	@cd tools && ./makeservermsg.sh

.PHONY: msg lint test
