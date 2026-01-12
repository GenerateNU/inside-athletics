package routeTests

import (
	"encoding/json"
	"inside-athletics/internal/server"
	"log"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupTestAPI(t *testing.T) humatest.TestAPI {
	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	_, api := humatest.New(t) // setup test API

	dbUrl := os.Getenv("PROD_DB_CONNECTION_STRING")

	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})

	server.CreateRoutes(db, api)

	return api
}

/*
Decode the given response JSON into the given struct entity. Reads the value into the struct
*/
func DecodeTo[T any](entity *T, resp *httptest.ResponseRecorder) {
	dec := json.NewDecoder(resp.Body)

	err := dec.Decode(entity)
	if err != nil {
		log.Fatalf("decode error: %v", err)
	}
}
