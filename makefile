commit: ## build and push container to docker hub| provide tag=(...)
	@docker build -t dmitryavdonin/promis-monapi:2.0.1 . && docker push dmitryavdonin/promis-monapi:2.0.1
.PHONY: commit
