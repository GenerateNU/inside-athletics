# Makerfile definding some targets for this project. To 
# run any of these commands run make <target name>

# Generates a new migration if there have been changes
# made to any of the models
gen-migration:
	atlas migrate dif --env gorm

# Applies migration to development (local) supabase db
migrate-dev:
	atlas migrate apply --env dev

# Applies migration to production supabase db
migrate-prod:
	atlas migrate apply --env prod

# starts server utilizing local supabase db
start-dev:
	cd cmd && APP_ENV=development go run . 

# starts server utilizing prod supabase db
start-prod:
	cd cmd && APP_ENV=production go run . 

# runs the linter
lint:
	golangci-lint run