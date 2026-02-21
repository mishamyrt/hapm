VERSION = 0.4.0
NAME = $(shell uname -a)

.PHONY: publish
publish: clean build
	$(VENV) python3 -m twine upload --repository pypi dist/*
	git add Makefile
	git commit -m "chore: release v$(VERSION)"
	git tag "v$(VERSION)"
	git push
	git push --tags

.PHONY: clean
clean:
	rm -rf build

.PHONY: build
build:
	CGO_ENABLED=0 \
		go build \
			-ldflags "-s -w" \
			-o build/hapm \
			hapm.go

.PHONY: install
install:
	go install

.PHONY: lint
lint:
	golangci-lint run ./...
	revive -config ./revive.toml  ./...

.PHONY: test
test:
	@go test ./...

.PHONY: test-e2e
test-e2e:
	@go test \
		-race \
		-count=1 \
		-timeout=30s \
		-tags=e2e \
		e2e_test.go