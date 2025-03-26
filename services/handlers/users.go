package handlers

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "time"

    "coffeeApi/services/db"

    "github.com/golang-jwt/jwt"
    "golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("super_sekretny_klucz_kofola_5mlnZł")

type User struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Password string `json:"password"`
    Email    string `json:"email"`
    Role     string `json:"role"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    if user.Username == "" || user.Password == "" || user.Email == "" {
        http.Error(w, "Missing required fields", http.StatusBadRequest)
        return
    }

    // Hashujemy hasło
    hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Error hashing password", http.StatusInternalServerError)
        return
    }
    user.Password = string(hashed)
    if user.Role == "" {
        user.Role = "user"
    }

    err = db.DB.QueryRow(
        `INSERT INTO users (username, password, email, role) VALUES ($1, $2, $3, $4) RETURNING id`,
        user.Username, user.Password, user.Email, user.Role,
    ).Scan(&user.ID)
    if err != nil {
        http.Error(w, "Error inserting user: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Nie zwracamy hasła
    user.Password = ""
    json.NewEncoder(w).Encode(user)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    var credentials struct {
        Username  string `json:"username"`
        Passwords string `json:"passwords"`
    }
    if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    var user User
    err := db.DB.QueryRow(
        `SELECT id, username, password, email, role FROM users WHERE username=$1`,
        credentials.Username,
    ).Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Role)
    if err == sql.ErrNoRows {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    } else if err != nil {
        http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Porównujemy zahaszowane hasło
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Passwords)); err != nil {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    expirationTime := time.Now().Add(24 * time.Hour)
    claims := jwt.MapClaims{
        "userId": user.ID,
        "exp":    expirationTime.Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        http.Error(w, "Error generating token", http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}