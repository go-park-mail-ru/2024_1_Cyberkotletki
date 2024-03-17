PKG_INTERNAL = ./internal/...

.PHONY: run-tests
run-tests:
	go test -race $(PKG_INTERNAL) -coverprofile=test.coverage.tmp $(PKG_INTERNAL)
	cat test.coverage.tmp | grep -v 'mocks' > test.coverage
	go tool cover -func test.coverage | tail -n 1 && rm test.coverage.tmp && rm test.coverage

.PHONY: gen-swagger
gen-swagger:
	swag init --dir cmd/app,internal/delivery --parseDependency

.PHONY: run-dev
run-dev:
	make gen-swagger
	go run cmd/app/main.go

.PHONY: gen-example-config
gen-example-config:
	go run cmd/app/main.go --generate-example-config=true

