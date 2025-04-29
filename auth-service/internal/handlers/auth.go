package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/cushydigit/microstore/auth-service/internal/auth"
	"github.com/cushydigit/microstore/auth-service/internal/helpers"
	"github.com/cushydigit/microstore/auth-service/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func Register(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			helpers.ErrorJSON(w,err)
			return
		}
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		_, err := db.Exec("INSERT INTO users (email, passoword) VALUES ($1, $2)", user.Email, hashedPassword)
		if err != nil {
			helpers.ErrorJSON(w, err, http.StatusInternalServerError)
			return 
		}

		paylaod := helpers.ResponseJSON{
			Error: false,
			Message: "User registered successfully",
		}
		helpers.WriteJSON(w, http.StatusCreated, paylaod)

	}
}

func Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds models.User
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			helpers.ErrorJSON(w, err)
			return
		}

		var storedPasseord string
		var userID int
		if err := db.QueryRow("SELECT id, password FROM users WHERE email=$1", creds.Email).Scan(&userID, &storedPasseord); err != nil {
			helpers.ErrorJSON(w, errors.New("Invalid credentials"), http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(storedPasseord), []byte(creds.Password)); err != nil {
			helpers.ErrorJSON(w, errors.New("Invalid credentials"), http.StatusUnauthorized)
			return
		}

		token, err := auth.GenerateJWT(userID)
		if err != nil {
			helpers.ErrorJSON(w, errors.New("Internal server error"), http.StatusInternalServerError)
		}

		payload := helpers.ResponseJSON{
			Error: false,
			Message: "Login successfully",
			Data: &helpers.LoginPayload {
				Token: token,
				User: helpers.UserPayload {
					ID: userID,
					Email: creds.Email,
				},
			},
		}

		helpers.WriteJSON(w, http.StatusOK, payload )

	}
}
		

