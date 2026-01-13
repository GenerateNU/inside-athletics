gen-migration:
	atlas migrate dif --env gorm

migrate-dev:
	atlas migrate apply --env dev

migrate-prod:
	atlas migrate apply --env prod

start-dev:
	APP_ENV=development go run . 

start-prod:
	APP_ENV=production go run . 