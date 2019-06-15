
#export GOPATH=$(PWD)/../

COMM_PATH=./comm/

EXCELJSON_PATH = ./exceljson 
CEHUA_SVNURL = https://192.168.150.238/svn/GiantCode/wxcat/cehua/$(USER)/

SUB_DIRS = SuperServer UserServer GatewayServer BridgeServer MatchServer RoomServer LoginServer

all: debug 

debug:
	@echo $(GOPATH)
	@echo "Buildding message..."
	@cd tools && ./makemsg.sh
	@echo "Buildding src..."
	@echo ${GOPATH}
	@for dir in $(SUB_DIRS); do \
		go install -x -gcflags "-N -l" ./$$dir || exit 1; \
	done

proto:
	protoc -I=$(COMM_PATH) --proto_path=$(COMM_PATH) --go_out=. $(COMM_PATH)/*.proto

clean:
	@for bdir in $(SUB_DIRS); do \
		rm -rf ../bin/$$bdir; \
	done
	@find -name "*~" | xargs rm -f
	@find -name "*.swp" | xargs rm -f
	@cd test;make clean
	@rm -rf release/*
wc:
	@find . -iname \*.go -exec cat \{\} \; | wc -l

tags:
	@ctags -R

res:
	@rm -rf $(EXCELJSON_PATH)
	@if [ ! -d "$(EXCELJSON_PATH)" ]; then mkdir -pv $(EXCELJSON_PATH); fi
	@svn export --force $(CEHUA_SVNURL)/exceljson $(EXCELJSON_PATH)

.PHONY: all debug clean wc image



