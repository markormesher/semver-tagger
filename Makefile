.DEFAULT_GOAL: build

.PHONY: lint
lint:
	bad_files=$$(gofmt -l . | (grep -v ^vendor/ || :))
	if [[ "$${bad_files}" != "" ]]; then
		echo "Lint failed on the following files:"
		echo $${bad_files}
		exit 1
	fi

.PHONY: test
test:
	go test ./...

.PHONY: build
build:
	go build -o ./build/main ./main.go

.PHONY: run
run:
	go run ./main.go

.ONESHELL:
