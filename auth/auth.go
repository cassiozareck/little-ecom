package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"net/mail"
	"os"
	"time"
)

var db *sql.DB

func Register(w http.ResponseWriter, r *http.Request) {
	log.Println("Registering")

	// Parse the request body to get account credentials
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the email
	if _, err := mail.ParseAddress(creds.Email); err != nil {
		http.Error(w, "Invalid email address: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Verify if the account already exists
	var storedEmail string
	err = db.QueryRow(`SELECT username FROM account WHERE username = $1`, creds.Email).Scan(&storedEmail)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Error while querying the database: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Account already exists", http.StatusBadRequest)
		return
	}

	if len(creds.Password) < 6 {
		http.Error(w, "Password must be at least 6 characters long", http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error while hashing the password: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Insert the account into the database
	_, err = db.Exec(`INSERT INTO account(username, password) VALUES($1, $2)`, creds.Email, hashedPassword)
	if err != nil {
		http.Error(w, "Error while storing account: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

var jwtKey = []byte(os.Getenv("JWT_KEY"))

func SignIn(w http.ResponseWriter, r *http.Request) {
	log.Println("Signing in")

	// Parse the request body
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the account's hashed password from the database
	var storedHashedPassword string
	err = db.QueryRow(`SELECT password FROM account WHERE username = $1`, creds.Email).Scan(&storedHashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "account not found", http.StatusUnauthorized)
		} else {
			http.Error(w, "Error while querying the database: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Compare the stored hashed password with the provided password
	err = bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(creds.Password))
	if err != nil {
		http.Error(w, "Invalid credentials: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Create a new token object, specifying signing method and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": creds.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign the token with our secret key
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Error while signing the token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(tokenString))
}

func ValidateToken(w http.ResponseWriter, r *http.Request) {
	log.Println("Validating token")

	type TokenRequest struct {
		Token string `json:"token"`
	}

	var req TokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tokenString := req.Token

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"], err.Error())
		}
		return jwtKey, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Token is valid
		email := claims["email"].(string)
		exp := claims["exp"].(float64)

		response := struct {
			Email string  `json:"email"`
			Exp   float64 `json:"exp"`
		}{
			Email: email,
			Exp:   exp,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
	}
}

func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	log.Println("Deleting account")

	// Parse the request body
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the account's hashed password from the database
	var storedHashedPassword string
	err = db.QueryRow(`SELECT password FROM account WHERE username = $1`, creds.Email).Scan(&storedHashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "account not found", http.StatusUnauthorized)
		} else {
			http.Error(w, "Error while querying the database: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Compare the stored hashed password with the provided password
	err = bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(creds.Password))
	if err != nil {
		http.Error(w, "Invalid credentials: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Delete the account from the database
	_, err = db.Exec(`DELETE FROM account WHERE username = $1`, creds.Email)
	if err != nil {
		http.Error(w, "Error while deleting account: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	// Build a PostgreSQL connection string using the environment variables DB_account and DB_PASSWORD
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := "5432"
	dbName := "auth"

	connectionString := "postgres://" + dbUser + ":" + dbPassword + "@" + dbHost + ":" + dbPort + "/" + dbName + "?sslmode=disable"

	log.Println("Connecting with: " + connectionString)

	// Open a connection to the PostgreSQL database
	var err error
	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	// Check if the connection to the database is successful
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	// Create a new Gorilla Mux router
	r := mux.NewRouter()

	r.HandleFunc("/auth/register", Register).Methods("POST")
	r.HandleFunc("/auth/signin", SignIn).Methods("POST")
	r.HandleFunc("/auth/validate", ValidateToken).Methods("POST")
	r.HandleFunc("/auth/delete", DeleteAccount).Methods("DELETE")
	// Start the HTTP server
	http.ListenAndServe(":8080", r)
}
