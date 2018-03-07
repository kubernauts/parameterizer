
test:
	cd $(GOPATH)/src/github.com/kubernauts/parameterizer/pkg/parameterizer && go test


install:
	go install github.com/kubernauts/parameterizer/cli/krm

build:
	mkdir -p build
	go build -o build/krm github.com/kubernauts/parameterizer/cli/krm
.PHONY: test
