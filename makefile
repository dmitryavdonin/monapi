commit: ## build and push container to docker hub| provide tag=(...)
	@docker build -t dmitryavdonin/promis-monapi:2.0.9 . && docker push dmitryavdonin/promis-monapi:2.0.9
.PHONY: commit
