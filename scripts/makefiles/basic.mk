GOARCH  = amd64
GOOS    = $(OS)

ifeq ($(GOOS),)
  ifeq ($(shell  uname -s), Darwin)
    GOOS	= darwin
  else
    GOOS	= linux
  endif
endif

TAG			= make
MAIN		= main.go
BIN_DIR		= bin
CMD_DIR		= cmd
VENDOR_DIR	= vendor

GOCMD		= go
GODEP		= dep
GOVVV		= govvv
GOLINT		= gometalinter
GOBUILD		= GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOCMD) build
GOFMT		= $(GOCMD)fmt
GOTEST		= ginkgo
GOTEST_OPT	= -r -p -v
GO_FILES	= $(shell $(GOCMD) list ./... | grep -v /vendor/)

GIT_COMMIT	= $(shell git rev-parse HEAD)
GIT_DIRTY	= $(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
VERSION		= $(shell cat ./VERSION)
