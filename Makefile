migrate-create:
	docker compose run --rm migrate $(name)

sqlc:
	docker compose run --rm sqlc

down:
	docker compose down

gen-mock-service:
	mockgen \
	-destination=internal/mocks/service/mock_$(service)_service.go \
	-package=service_mock \
	github.com/loanem-backend/inventory-service/internal/service \
	$(interface)
