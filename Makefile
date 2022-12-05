.PHONY: help

help: ## show help
	@awk 'BEGIN {FS = ":.*##"; printf "Usage:\n  make \033[36m\033[0m\n"} /^[$$()% 0-9a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

setup_db: ## setup database for development
	./scripts/init_db.sh

migrate: export APP_ENV=dev
migrate: ## run database migrations
	go run ./cmd/migration

test_setup_db: ## setup database for testing
	./scripts/init_db_test.sh
test: ## run repository tests WITHOUT coverage
	go test ./pkg/adapter/repository
test_cov: ## run repository tests and gerate coverage
	rm -r ./coverage
	mkdir ./coverage
	go test ./pkg/adapter/repository -cover -coverprofile=coverage/repo.out
	go tool cover -html=coverage/repo.out -o coverage/repo.html

e2e_setup_db: ## setup database for e2e testing
	./scripts/init_db_e2e.sh
e2e: ## run End-To-End tests
	go test ./test/e2e/...

start: ## start development server
	air

docs: ## generate graphql schema docs
	@echo "\033[0;33mMake sure you have run gqlgen and restart the server\033[0m"
	graphdoc -e http://localhost:8080/query -o ./docs/schema --force
	@echo "\033[0;32mOpen file://${PWD}/docs/schema/index.html in the browser to view the docs \033[0m"