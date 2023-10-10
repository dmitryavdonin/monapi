commit: ## build and push container to docker hub| provide tag=(...)
	@docker build -t dmitryavdonin/promis-monapi:2.0.7 . && docker push dmitryavdonin/promis-monapi:2.0.7
.PHONY: commit
