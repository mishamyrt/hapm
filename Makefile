VERSION = 0.4.0
NAME = $(shell uname -a)
DOCKER_IMAGE = mishamyrt/hapm

.PHONY: publish
publish: clean build
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

.PHONY: docker-build
docker-build:
	docker buildx build \
		--platform linux/amd64,linux/arm64 \
		-t $(DOCKER_IMAGE):$(VERSION) \
		-t $(DOCKER_IMAGE):latest \
		.

.PHONY: docker-publish
docker-publish:
	docker buildx build \
		--platform linux/amd64,linux/arm64 \
		-t $(DOCKER_IMAGE):$(VERSION) \
		-t $(DOCKER_IMAGE):latest \
		--push \
		.