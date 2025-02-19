.PHONY: generate run health auth check-routes

generate:
	go mod tidy
	go mod download
	go generate ./...

run: generate
	go run .

health:
	@echo "Check /health route..."
	@{ \
  	HEALTH_STATUS=$$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/health); \
	if [ "$$HEALTH_STATUS" -ne 200 ]; then \
		echo "Error: /health route returned status $$HEALTH_STATUS instead of 200"; \
		exit 1; \
	fi; \
	echo "/health route returned status 200"; \
	}


USER := testUserName
auth:
	@echo "Get /api/auth authorization..."
	@{ \
	TOKEN=$$(curl -s -X POST http://localhost:8080/api/auth \
		-H 'Content-Type: application/json' \
		-d '{"username": "$(USER)", "password": "pass"}' \
		| jq -r '.token'); \
	if [ -z "$$TOKEN" ]; then \
	  	echo 'Error: token not received'; \
	  	exit 1; \
	fi; \
	echo "Token received: $$TOKEN"; \
	echo 'Get private /api/private route...'; \
	API_USER=$$(curl -s -X GET http://localhost:8080/api/private \
		-H "Authorization: Bearer $$TOKEN" \
		| jq -r '.username'); \
	if [ "$$API_USER" != "$(USER)" ]; then \
		echo "Error: API username ($$API_USER) does not match expected username ($(USER))"; \
		exit 1; \
	fi; \
	echo "API username ($$API_USER) matches expected username ($(USER))"; \
	}

check-routes: health auth
