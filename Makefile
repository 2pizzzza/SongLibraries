MAIN_PACKAGE_PATH := ./cmd/songLibraries
BINARY_NAME := songLibraries

.PHONY: build
build:
	go build -o=/tmp/bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH}

.PHONY: run
run: build
	/tmp/bin/${BINARY_NAME}
