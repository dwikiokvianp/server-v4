package test

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"reflect"
	"server-v2/controllers"
	"testing"
	"time"
)

// MockDB is a mock implementation of the database.
type MockDB struct {
}

// Find is a mock method for finding transactions in the database.
func (db *MockDB) Find(transactions *[]Transaction) error {
	// Mock implementation to return transactions for testing
	*transactions = []Transaction{
		Transaction{
			ID:   1,
			Date: time.Now().Format("2006-01-02"),
		},
	}

	return nil
}

type Transaction struct {
	ID   int
	Date string
}

func TestGetTodayV2Transaction(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/transactions", nil)

	rec := httptest.NewRecorder()

	// Call the handler function
	router.ServeHTTP(rec, req)

	// Check the response status code
	if rec.Code != http.StatusOK {
		t.Errorf("expected status OK; got %v", rec.Code)
	}

	// Decode the response body
	var response struct {
		Data []Transaction `json:"data"`
	}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Check the response data
	expectedTransaction := Transaction{
		ID:   1,
		Date: time.Now().Format("2006-01-02"),
	}
	if len(response.Data) != 1 || !reflect.DeepEqual(response.Data[0], expectedTransaction) {
		t.Errorf("expected transaction %+v; got %+v", expectedTransaction, response.Data)
	}
}

func setupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/transactions", controllers.GetTodayV2Transaction)

	return router
}
