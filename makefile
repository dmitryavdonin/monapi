commit: ## build and push container to docker hub| provide tag=(...)
	@docker build -t dmitryavdonin/otus-dz10-billing:latest . && docker push dmitryavdonin/otus-dz10-billing:latest
.PHONY: commit
