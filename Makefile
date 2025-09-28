test:
	@./go.test.sh
.PHONY: test

coverage:
	@./go.coverage.sh
.PHONY: coverage

test-fast:
	go test ./...

tidy:
	go mod tidy

migrate:
	atlas --config file://prod.atlas.hcl --env prod migrate apply

generate:
	go generate ./...

lint-openapi:
	go tool vacuum lint -d _oas/openapi.yaml

lint-otel-registry:
	weaver registry check -r _otel

generate-otel-bundle:
	weaver registry resolve -r _otel > .otel.registry.yml
	go run ./internal/cmd/otel-sort -f .otel.registry.yml

generate-otel: generate-otel-bundle
