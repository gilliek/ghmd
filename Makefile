EXEC = ghmd
PKG  = github.com/gilliek/ghmd

all: build

install:
	go install ${PKG}

uninstall:
	rm -f ${GOPATH}/bin/${EXEC}

build:
	go build ${PKG}

clean:
	rm -f ${ghmd}

