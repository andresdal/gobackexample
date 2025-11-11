package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andresdal/gobackexample/types"
	"github.com/gorilla/mux"
)

func TestUserServiceHandlers(t *testing.T) {
	userStore := &MockUserStore{}
	handler := NewHandler(userStore)

	t.Run("should fail if user payload is invalid", func(t *testing.T) {
		payload := types.RegisterUserPayload{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "invalid",
			Password:  "123",  
		}
		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		
		router.HandleFunc("/register", handler.handleRegister).Methods("POST")
		router.ServeHTTP(rr, req)

		if(rr.Code != http.StatusBadRequest) {
			t.Errorf("expected status %d but got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should succeed if user payload is valid", func(t *testing.T) {
		payload := types.RegisterUserPayload{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "valid@gmail.com",
			Password:  "123",  
		}
		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		
		router.HandleFunc("/register", handler.handleRegister).Methods("POST")
		router.ServeHTTP(rr, req)

		if(rr.Code != http.StatusCreated) {
			t.Errorf("expected status %d but got %d", http.StatusCreated, rr.Code)
		}
	})
}

type MockUserStore struct{}

func (m *MockUserStore) GetUserByEmail(email string) (*types.User, error) {
	return nil, fmt.Errorf("user not found")
}

func (m *MockUserStore) GetUserByID(id int) (*types.User, error) {
	return nil, nil
}

func (m *MockUserStore) CreateUser(user types.User) error {
	return nil
}