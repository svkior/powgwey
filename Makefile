# Main targets
.PHONY: build-base
build-base: ## build base golang working image
	@docker build \
		-t localhost:5000/air:latest \
		-f  deploy/base-image/Dockerfile \
		deploy/base-image/

.PHONY: generate
generate: ## regenerate all generated code
	docker run -it --rm \
		-v ${PWD}:/project \
		-v golang-cache-vol:/go/pkg/mod \
		-v go-build-vol:/root/.cache/go-build \
		--workdir="/project" \
		--entrypoint=go \
		localhost:5000/air \
		generate ./...

.PHONY: build-server
build-server: ## build server production image
	@docker build \
		-t localhost:5000/server:latest \
		-f deploy/server-image/Dockerfile \
		.

.PHONY: build-client
build-client: ## build server production image
	@docker build \
		-t localhost:5000/client:latest \
		-f deploy/client-image/Dockerfile \
		.

.PHONY: start-server
start-server: ## start server in detached container
	@docker run --rm --name server \
		--env-file "./deploy/server.env" \
		-d --network local \
		localhost:5000/server:latest

.PHONY: start-client
start-client: ## start server in detached container
	@docker run --rm --name client \
		--network local \
		localhost:5000/client:latest


.PHONY: watch-server
watch-server: ## start server in autoreload mode
	@docker run -it --rm --name server \
		--network local \
		--env-file "./deploy/server.env" \
		-v ${PWD}:/project \
		-v golang-cache-vol:/go/pkg/mod \
		-v go-build-vol:/root/.cache/go-build \
		--workdir="/project" \
		localhost:5000/air -c deploy/server.air.toml

.PHONY: watch-client
watch-client: ## start server in autoreload mode
	@docker run -it --rm --name client \
		--network local \
		-v ${PWD}:/project \
		-v golang-cache-vol:/go/pkg/mod \
		-v go-build-vol:/root/.cache/go-build \
		--workdir="/project" \
		localhost:5000/air -c deploy/client.air.toml

.PHONY: update-quotes
update-quotes: ## Update quotes from external link
	@mkdir -p data/quotes
	@curl -o data/quotes/movies.json https://raw.githubusercontent.com/msramalho/json-tv-quotes/master/quotes.json 

.PHONY: lint
lint: ## Lint all golang code
	docker run --rm \
		-v ${PWD}:/project \
		-v golang-cache-vol:/go/pkg/mod \
		-v go-build-vol:/root/.cache/go-build \
		--workdir="/project" \
		--entrypoint=golangci-lint \
		localhost:5000/air run -v ./...

.PHONY: test
test: ## test all golang code
	docker run --rm \
		-v ${PWD}:/project \
		-v golang-cache-vol:/go/pkg/mod \
		-v go-build-vol:/root/.cache/go-build \
		--workdir="/project" \
		--entrypoint=go \
		localhost:5000/air \
		test -race -covermode=atomic \
		-coverprofile=coverage.out ./...

.PHONY: watch-test
watch-test: ## test all golang code
	@docker run --rm --name watch-test \
		-v ${PWD}:/project \
		-v golang-cache-vol:/go/pkg/mod \
		-v go-build-vol:/root/.cache/go-build \
		-p 8080:8080 \
		--workdir="/project" \
		--entrypoint=goconvey \
		localhost:5000/air \
		-excludedDirs "data,deploy,tmp,cmd" -host "0.0.0.0" 

.PHONY: watch-testlocal
watch-testlocal: ## test all golang code on local environment
	~/go/bin/goconvey \
	-excludedDirs "data,deploy,tmp,cmd,server_internal/models,server_internal/app" -host "0.0.0.0" 

.PHONY: help
help: ## Print this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

