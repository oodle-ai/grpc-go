all: vet test testrace

build:
	go build github.com/oodle-ai/grpc-go/...

clean:
	go clean -i github.com/oodle-ai/grpc-go/...

deps:
	GO111MODULE=on go get -d -v github.com/oodle-ai/grpc-go/...

proto:
	@ if ! which protoc > /dev/null; then \
		echo "error: protoc not installed" >&2; \
		exit 1; \
	fi
	go generate github.com/oodle-ai/grpc-go/...

test:
	go test -cpu 1,4 -timeout 7m github.com/oodle-ai/grpc-go/...

testsubmodule:
	cd security/advancedtls && go test -cpu 1,4 -timeout 7m github.com/oodle-ai/grpc-go/security/advancedtls/...
	cd security/authorization && go test -cpu 1,4 -timeout 7m github.com/oodle-ai/grpc-go/security/authorization/...

testrace:
	go test -race -cpu 1,4 -timeout 7m github.com/oodle-ai/grpc-go/...

testdeps:
	GO111MODULE=on go get -d -v -t github.com/oodle-ai/grpc-go/...

vet: vetdeps
	./scripts/vet.sh

vetdeps:
	./scripts/vet.sh -install

.PHONY: \
	all \
	build \
	clean \
	deps \
	proto \
	test \
	testsubmodule \
	testrace \
	testdeps \
	vet \
	vetdeps
