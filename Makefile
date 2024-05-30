PKG_INTERNAL = ./internal/...

.PHONY: run-tests
run-tests:
	@go generate $(PKG_INTERNAL)
	@if go test -race $(PKG_INTERNAL) -coverprofile=test.coverage.tmp $(PKG_INTERNAL) ; then \
    	cat test.coverage.tmp | grep -v 'mocks' > test.coverage ; \
    	go tool cover -func test.coverage | tail -n 1 && rm test.coverage.tmp && rm test.coverage ; \
    	echo "\033[0;32mТесты прошли успешно\033[0m" ; \
    else \
    	echo "\033[0;31mТесты обнаружили проблемы\033[0m" ; \
    	exit 1 ; \
    fi

.PHONY: gen-swagger
gen-swagger:
	swag init --dir cmd/app,internal/delivery/http --parseDependency

.PHONY: run-full-dev
run-full-dev:
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

.PHONY: run-dev
run-dev:
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
		docker run --publish=5432:5432 --name db -e POSTGRES_PASSWORD=default_password -v $(PWD)/config/postgresql.conf:/etc/postgresql/postgresql.conf -d postgres -c 'config_file=/etc/postgresql/postgresql.conf'; \
	fi
	sleep 2
	# создаём базу данных и пользователя, выдаём права
	-docker exec -it db psql -U postgres -c "CREATE USER kinoskop_admin WITH SUPERUSER PASSWORD 'admin_secret_password'"
	-docker exec -it db psql -U postgres -c "CREATE DATABASE kinoskop"
	-docker exec -it db psql -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE kinoskop TO kinoskop_admin"
	-docker exec -it db psql -U postgres -c "ALTER DATABASE kinoskop OWNER TO kinoskop_admin;"

.PHONY: create-db-from-state
create-db-from-state:
	# Выполняем SQL-скрипт из файла state1.sql
	docker cp $(PWD)/db/states/state1.sql db:/state.sql
	docker exec -i db bash -c 'psql -U kinoskop_admin kinoskop < /state.sql'

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
		exit 1; \
	fi
	goose -dir=db/migrations create $(name) sql

.PHONY: run-linter
run-linter:
	@if golangci-lint run ; then \
    	echo "\033[0;32mЛинтер прошел успешно\033[0m" ; \
    else \
    	echo "\033[0;31mЛинтер обнаружил проблемы\033[0m" ; \
    	exit 1 ; \
    fi

.PHONY: before-push
before-push: run-linter run-tests
	@echo "\033[0;32mВсе проверки прошли успешно. Можно делать git push.\033[0m"

.PHONY: gen-wrk-report
gen-wrk-report:
	# Тест на создание отзыва
	#
	wrk -t5 -c10 -d60m -s perf_test/post.lua http://127.0.0.1:8080/api/review > perf_test/outpost.txt  
	#
	cat perf_test/outpost.txt
	#
	# Тест на получение отзыва
	#
	wrk -t5 -c10 -d60m -s perf_test/get_after_post.lua http://127.0.0.1:8080/api/review/50000 > perf_test/outget.txt
	#
	cat perf_test/outget.txt