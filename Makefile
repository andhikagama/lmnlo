BINARY=lmnlo
OS=linux

build:
		GOOS=${OS} GOHOSTOS=${OS} go build -o ${BINARY}

test:
	./test-cover.sh

install: test
		go build -o ${BINARY}

unittest:
	go test -short $$(go list ./... | grep -v /vendor/)

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

.PHONY: clean install unittest test
