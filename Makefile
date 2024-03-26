PKG_INTERNAL = ./internal/...

.PHONY: run-tests
run-tests:
	go generate $(PKG_INTERNAL)
	go test -race $(PKG_INTERNAL) -coverprofile=test.coverage.tmp $(PKG_INTERNAL)
	cat test.coverage.tmp | grep -v 'mocks' > test.coverage
	go tool cover -func test.coverage | tail -n 1 && rm test.coverage.tmp && rm test.coverage

.PHONY: gen-swagger
gen-swagger:
	swag init --dir cmd/app,internal/delivery --parseDependency

.PHONY: run-dev
run-dev:
	make gen-swagger
	make run-session-storage-container
	go run cmd/app/main.go

.PHONY: gen-example-config
gen-example-config:
	go run cmd/app/main.go --generate-example-config=true

.PHONY: run-session-storage-container
run-session-storage-container:
	# Хранилище сессий запустится в redis на дефолтном порту 6379. Если контейнер уже запущен, то ничего не произойдёт
	if ! docker inspect -f '{{.State.Running}}' session-storage 2>/dev/null ; then \
		docker run --publish=6379:6379 --name session-storage -d redis redis-server --maxmemory 1gb; \
    fi

.PHONY: run-linter
run-linter:
	golangci-lint run
