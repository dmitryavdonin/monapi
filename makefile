commit: ## build and push container to docker hub| provide tag=(...)
	@docker build -t dmitryavdonin/promis-monapi:latest . && docker push dmitryavdonin/promis-monapi:latest
.PHONY: commit
