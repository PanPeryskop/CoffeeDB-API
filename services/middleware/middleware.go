package middleware

import (
    "database/sql"
    "fmt"
    "net/http"
    "strconv"
    "strings"

    "coffeeApi/services/db"
    "github.com/golang-jwt/jwt"
)

var jwtKey = []byte("super_sekretny_klucz_kofola_5mlnZł")

// AuthMiddleware validates the JWT token and sets the "X-User-ID" header.
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Authorization header required", http.StatusUnauthorized)
            return
        }
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 {
            http.Error(w, "Invalid token format", http.StatusUnauthorized)
            return
        }
        token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            return jwtKey, nil
        })
        if err != nil || !token.Valid {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            http.Error(w, "Invalid token claims", http.StatusUnauthorized)
            return
        }
        userIDFloat, ok := claims["userId"].(float64)
        if !ok {
            http.Error(w, "Invalid userId in token", http.StatusUnauthorized)
            return
        }
        r.Header.Set("X-User-ID", strconv.Itoa(int(userIDFloat)))
        next.ServeHTTP(w, r)
    })
}

// AdminMiddleware ensures that the authenticated user has admin privileges.
// It queries the database for the user's role.
func AdminMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        userIDStr := r.Header.Get("X-User-ID")
        if userIDStr == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        userID, err := strconv.Atoi(userIDStr)
        if err != nil {
            http.Error(w, "Invalid user ID", http.StatusUnauthorized)
            return
        }
        var role string
        err = db.DB.QueryRow(`SELECT role FROM users WHERE id = $1`, userID).Scan(&role)
        if err == sql.ErrNoRows {
            http.Error(w, "User not found", http.StatusUnauthorized)
            return
        } else if err != nil {
            http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
            return
        }
        if role != "admin" {
            http.Error(w, "Admin privileges required", http.StatusForbidden)
            return
        }
        next.ServeHTTP(w, r)
    })
}