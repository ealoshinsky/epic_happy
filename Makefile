release ?= 0.1.0
commit?=$(shell git rev-parse --short HEAD)
build_time?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
app_name?=eh
os ?= linux
arch ?= amd64

.PHONY: check
check: prepare_metalinter
	gometalinter --deadline=240s --vendor ./...

.PHONY: build
build: clean
	CGO_ENABLED=0 GOOS=${os} GOARCH=${arch} go build \
		-ldflags "-X main.release=${release} -X main.commit=${commit} -X main.buildTime=${build_time} -X main.appName=${app_name}" \
		-o  build/${os}/${app_name}-${arch}

.PHONY: clean
clean:
	@rm -rf build/

.PHONY: vendor
vendor: prepare_dep
	dep ensure

HAS_DEP := $(shell command -v dep;)
HAS_METALINTER := $(shell command -v gometalinter;)

.PHONY: prepare_dep
prepare_dep:
ifndef HAS_DEP
	go get -u -v -d github.com/golang/dep/cmd/dep && \
	go install -v github.com/golang/dep/cmd/dep
endif

.PHONY: prepare_metalinter
prepare_metalinter:
ifndef HAS_METALINTER
	go get -u -v -d github.com/alecthomas/gometalinter && \
	go install -v github.com/alecthomas/gometalinter && \
	gometalinter --install --update
endif