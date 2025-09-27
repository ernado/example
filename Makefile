test:
	@./go.test.sh
.PHONY: test

coverage:
	@./go.coverage.sh
.PHONY: coverage

test_fast:
	go test ./...

tidy:
	go mod tidy

migrate:
	atlas --config file://prod.atlas.hcl --env prod migrate apply

generate:
	go generate ./...
