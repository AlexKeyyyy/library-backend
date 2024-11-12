package handlers

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"library-backend/db"
	"net/http"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Хэшируем пароль
	hash := sha256.New()
	hash.Write([]byte(user.Password))
	hashedPassword := hex.EncodeToString(hash.Sum(nil))

	// Проверяем логин и пароль в базе данных
	var storedPassword string
	err = db.DB.QueryRow("SELECT password FROM users WHERE username = $1", user.Username).Scan(&storedPassword)
	if err == sql.ErrNoRows || storedPassword != hashedPassword {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	w.Write([]byte("Login successful"))
}
