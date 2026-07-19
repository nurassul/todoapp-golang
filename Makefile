include .env
export

export PROJECT_ROOT=$(shell pwd)

env-up:
	@docker compose up -d todoapp-postgres

env-down:
	@docker compose down todoapp-postgres

env-cleanup:
	@read -p "Are you sure to clean all volumes? [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
	  docker compose down todoapp-postgres && \
	  rm -rf ${PROJECT_ROOT}/out/pgdata && \
	  echo "Files were cleaned"; \
	else \
	  echo "Cleaning files canceled"; \
	fi

migrate-create:
	@if [ -z "$(seq)" ]; then \
  		echo "Not found required params seq." \
  		exit 1; \
	fi; \
	docker compose run --rm todoapp-postgres-migrate \
	  create \
	  -ext sql \
	  -dir . /migrations \
	  -seq "$(seq)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
      	echo "Not found required params seq." \
      	exit 1; \
    fi; \
	docker compose run --rm todoapp-postgres-migrate \
    		-path ./migrations \
    		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todoapp-postgres:5432/${POSTGRES_DB}?sslmode=disable \
    		"$(action)"

env-port-forward:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down todoapp-postgres


logs-cleanup:
	@read -p "Clean all log files? [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
	  rm -rf ${PROJECT_ROOT}/out/logs && \
	  echo "Logs were cleaned up"; \
	else \
	  echo "Cleaning was canceled"; \
	fi

todoapp-run:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	export POSTGRES_HOST=localhost && \
	go mod tidy && \
	go run ${PROJECT_ROOT}/cmd/todoapp/main.go

















