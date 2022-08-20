# Main targets
.PHONY: base-image
base-image: ## build base golang working image
	docker build \
    -t localhost:5000/air:latest \
    -f  deploy/base-image/Dockerfile \
    deploy/base-image/

.PHONY: update-quotes
update-quotes: ## Update quotes from external link
	@mkdir -p data/quotes
	@curl -o data/quotes/movies.json https://raw.githubusercontent.com/msramalho/json-tv-quotes/master/quotes.json 

.PHONY: help
help: ## Print this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

