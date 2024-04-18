package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/beny996/rss-agregator/internal/auth"
	"github.com/beny996/rss-agregator/internal/database"
	"github.com/google/uuid"
)

func (apiConfig *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	user, err := apiConfig.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})

	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Couldn't create user: %v", err))
		return
	}

	respondWithJSON(w, 201, databaseUserToUser(user))
}

func (apiConfig *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request) {

	apiKey, err := auth.GetApiKey(r.Header)

	if err != nil {
		responseWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
		return
	}

	user, err := apiConfig.DB.GetUserByApiKey(r.Context(), apiKey)

	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Couldn't get user: %v", err))
		return

	}

	respondWithJSON(w, 200, databaseUserToUser(user))
}
