PKG_INTERNAL = ./internal/...

.PHONY: run-tests
run-tests:
	@go generate $(PKG_INTERNAL)
	@if go test -race $(PKG_INTERNAL) -coverprofile=test.coverage.tmp $(PKG_INTERNAL) ; then \
    	cat test.coverage.tmp | grep -v 'mocks' | grep -v 'proto' | grep -v 'easyjson' > test.coverage ; \
    	go tool cover -func test.coverage | tail -n 1 ; \
    	go tool cover -func test.coverage | grep 'total:' | awk '{print $$3}' && rm test.coverage.tmp && rm test.coverage ; \
    	echo "\033[0;32mТесты прошли успешно\033[0m" ; \
    else \
    	echo "\033[0;31mТесты обнаружили проблемы\033[0m" ; \
    	exit 1 ; \
    fi

.PHONY: gen-swagger
gen-swagger:
	swag init --dir cmd/app,internal/delivery/http --parseDependency

.PHONY: run-dev
run-dev:
	go run cmd/app/main.go

.PHONY: run-docker-deploy
run-docker-deploy:
	docker-compose up --build -d
	sleep 5
	make run-migrations

.PHONY: gen-example-config
gen-example-config:
	go run cmd/app/main.go --generate-example-config=true
	go run cmd/auth/main.go --generate-example-config=true
	go run cmd/static/main.go --generate-example-config=true

.PHONY: run-session-storage-container
run-session-storage-container:
	@if [ $$(docker ps -q -f name=session-storage) ]; then \
		echo "Контейнер с сессиями уже запущен"; \
	else \
		if [ $$(docker ps -aq -f name=session-storage) ]; then \
			echo "Контейнер с сессиями существует, но не запущен. Запускаем..."; \
			docker start session-storage; \
		else \
			echo "Создаём и запускаем контейнер с сессиями..."; \
			docker run -d --publish=6379:6379 --name session-storage redis --maxmemory 1gb; \
		fi \
	fi

.PHONY: run-db-container
run-db-container:
	@if [ $$(docker ps -q -f name=db) ]; then \
		echo "Контейнер с базой данных уже запущен"; \
	else \
		if [ $$(docker ps -aq -f name=db) ]; then \
			echo "Контейнер с базой данных существует, но не запущен. Запускаем..."; \
			docker start db; \
		else \
			echo "Создаём и запускаем контейнер с базой данных..."; \
			docker create --publish=5432:5432 --name db -e POSTGRES_PASSWORD=default_password postgres; \
			docker start db; \
		fi \
	fi
	sleep 3
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
	@if [ -z $(name) ]; then \
		echo "Нужно указать название миграции параметром 'name'"; \
		exit 1; \
	fi
	goose -dir=db/migrations create $(name) sql


.PHONY: easyjson
easyjson:
	@echo "Generating easyjson..."
	@for file in $$(find ./internal/entity/dto -name '*.go' | grep -v "_easyjson.go"); do \
		easyjson -all $$file; \
	done

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
