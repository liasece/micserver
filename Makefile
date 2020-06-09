
PKG_LIST := `go list ./... | grep -v "servercomm" | grep -v "testmsg"`

all: msg test

lint:
	@golint -set_exit_status $(PKG_LIST)

test:
	@go test -count=1 -short -timeout 10s $(PKG_LIST)

msg:
	@cd tools && ./makeservermsg.sh

.PHONY: all msg lint test
