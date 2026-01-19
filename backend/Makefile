# Makerfile definding some targets for this project. To 
# run any of these commands run make <target name>

# Generates a new migration if there have been changes
# made to any of the models

.PHONY: gen-migration
gen-migration:
	if [ -z "$(MIGRATION_NAME)" ]; then \
		echo "MIGRATION_NAME must be set"; \
	else \
		doppler run --command="atlas migrate diff --env dev $(MIGRATION_NAME)";	\
	fi

# Applies migration to development (local) supabase db
.PHONY: migrate-dev
migrate-dev:
	doppler run --command="atlas migrate apply --env dev"

# Applies migration to production supabase db
.PHONY: migrate-prod
migrate-prod:
	doppler run --command="atlas migrate apply --env prod"

# starts server utilizing local supabase db
.PHONY: start-dev
start-dev:
	doppler run --command="cd cmd && APP_ENV=development go run ."

# starts server utilizing prod supabase db
.PHONY: start-prod
start-prod:
	doppler run --command="cd cmd && APP_ENV=production go run ."

# runs the linter
.PHONY: lint
lint:
	doppler run --comand="golangci-lint run --timeout=5m"

# reverts single migration for prod
.PHONY: revert-prod
revert-single-prod:
	doppler run --command="atlas migrate down --env prod"

.PHONY: revert-dev
# reverts single migration for dev
revert-single-dev:
	doppler run --comand="atlas migrate down --env dev"

VERBOSE ?=
V_FLAG := $(if $(filter true,$(VERBOSE)),-v,)

.PHONY: run-all-tests
# Runs all tests, if passed verbose, log test details
run-all-tests:
	go test $(V_FLAG) ./internal/tests/route_tests ./internal/tests/unit_tests

.PHONY: run-route-tests
# Runs all the route tests, if passed verbose, log test details
run-route-tests:
	go test $(V_FLAG) ./internal/tests/route_tests

.PHONY: run-unit-tests
# Runs all the unit tests, if passed verbose, log test details
run-unit-tests:
	go test $(V_FLAG) ./internal/tests/unit_tests





