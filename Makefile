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
	# генерируем пример конфига
	make gen-example-config
	# генерируем документацию swagger
	make gen-swagger
	# запускаем контейнер с сессиями
	make run-session-storage-container
	# запускаем контейнер с базой данных
	make run-db-container
	# накатываем миграции
	make run-migrations
	# запускаем приложение
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

.PHONY: run-db-container
run-db-container:
	# PostgreSQL запустится на дефолтном порту 5432. Если контейнер уже запущен, то ничего не произойдёт
	if ! docker inspect -f '{{.State.Running}}' db 2>/dev/null ; then \
		docker run --publish=5432:5432 --name db -e POSTGRES_PASSWORD=default_password -d postgres; \
	fi
	sleep 2
	# создаём базу данных и пользователя, выдаём права
	-docker exec -it db psql -U postgres -c "CREATE USER kinoskop_admin PASSWORD 'admin_secret_password'"
	-docker exec -it db psql -U postgres -c "CREATE DATABASE kinoskop"
	-docker exec -it db psql -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE kinoskop TO kinoskop_admin"
	-docker exec -it db psql -U postgres -c "ALTER DATABASE kinoskop OWNER TO kinoskop_admin;"

.PHONY: run-migrations
run-migrations:
	# Сначала нужно накатить goose
	# https://pressly.github.io/goose/installation/
	goose -dir=db/migrations postgres postgres://kinoskop_admin:admin_secret_password@localhost:5432/kinoskop?sslmode=disable up

.PHONY: gen-migration
gen-migration:
	# Пример использования: make gen-migration name=create_users_table
	@if [ -z "$(name)" ]; then
		echo "Нужно указать название миграции параметром 'name'";
		exit 1;
	fi
	goose -dir=db/migrations create $(name) sql


.PHONY: run-linter
run-linter:
	golangci-lint run
