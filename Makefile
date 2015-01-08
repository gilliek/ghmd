EXEC = ghmd
PKG  = github.com/gilliek/ghmd

all: build

install:
	go install ${PKG}

uninstall:
	rm ${GOPATH}/bin/${EXEC}

build:
	go build ${PKG}

