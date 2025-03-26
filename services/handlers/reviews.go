package handlers

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "strconv"
    "time"
    "coffeeApi/services/db"
    "github.com/gorilla/mux"
)


type Review struct {
    ID             int       `json:"id"`
    UserId         int       `json:"userId"`
    CoffeeId       int       `json:"coffeeId"`
    RoasteryId     int       `json:"roasteryId"`
    CoffeeShopId   int       `json:"coffeeShopId"`
    Rating         float32   `json:"rating"`
    Review         string    `json:"review"`
    DateOfCreation time.Time `json:"dateOfCreation"`
}


func allowedRating(rating float32) bool {
    intRating := int(rating)
    if float32(intRating) != rating {
        return false
    }
    return intRating >= 1 && intRating <= 5
}


func GetReviewsHandler(w http.ResponseWriter, r *http.Request) {
    rows, err := db.DB.Query(`
        SELECT id, user_id, coffee_id, roastery_id, coffee_shop_id, rating, review, date_of_creation 
        FROM reviews`)
    if err != nil {
        http.Error(w, "Database query error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var reviews []Review
    for rows.Next() {
        var rev Review
        if err := rows.Scan(&rev.ID, &rev.UserId, &rev.CoffeeId, &rev.RoasteryId, &rev.CoffeeShopId, &rev.Rating, &rev.Review, &rev.DateOfCreation); err != nil {
            http.Error(w, "Error scanning review: "+err.Error(), http.StatusInternalServerError)
            return
        }
        reviews = append(reviews, rev)
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(reviews)
}


func CreateReviewHandler(w http.ResponseWriter, r *http.Request) {
    var rev Review
    if err := json.NewDecoder(r.Body).Decode(&rev); err != nil {
        http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
        return
    }


    if !allowedRating(rev.Rating) {
        http.Error(w, "Rating must be an integer between 1 and 5", http.StatusBadRequest)
        return
    }


    targetCount := 0
    if rev.CoffeeId != 0 {
        targetCount++
    }
    if rev.RoasteryId != 0 {
        targetCount++
    }
    if rev.CoffeeShopId != 0 {
        targetCount++
    }
    if targetCount != 1 {
        http.Error(w, "Review must target exactly one of: coffee, roastery, or coffee shop", http.StatusBadRequest)
        return
    }

    rev.DateOfCreation = time.Now()
    err := db.DB.QueryRow(`
        INSERT INTO reviews (user_id, coffee_id, roastery_id, coffee_shop_id, rating, review, date_of_creation)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id`,
        rev.UserId, rev.CoffeeId, rev.RoasteryId, rev.CoffeeShopId, rev.Rating, rev.Review, rev.DateOfCreation).Scan(&rev.ID)
    if err != nil {
        http.Error(w, "Database insert error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(rev)
}


func UpdateReviewHandler(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    reviewID, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid review ID", http.StatusBadRequest)
        return
    }


    var orig Review
    err = db.DB.QueryRow(`
        SELECT id, user_id, coffee_id, roastery_id, coffee_shop_id, rating, review, date_of_creation 
        FROM reviews WHERE id = $1`, reviewID).
        Scan(&orig.ID, &orig.UserId, &orig.CoffeeId, &orig.RoasteryId, &orig.CoffeeShopId, &orig.Rating, &orig.Review, &orig.DateOfCreation)
    if err == sql.ErrNoRows {
        http.Error(w, "Review not found", http.StatusNotFound)
        return
    } else if err != nil {
        http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
        return
    }


    userIDStr := r.Header.Get("X-User-ID")
    reqUserID, err := strconv.Atoi(userIDStr)
    if err != nil || reqUserID != orig.UserId {
        http.Error(w, "You can only update your own reviews", http.StatusForbidden)
        return
    }

    var rev Review
    if err := json.NewDecoder(r.Body).Decode(&rev); err != nil {
        http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
        return
    }
    if !allowedRating(rev.Rating) {
        http.Error(w, "Rating must be an integer between 1 and 5", http.StatusBadRequest)
        return
    }

    rev.UserId = orig.UserId
    rev.DateOfCreation = orig.DateOfCreation

    result, err := db.DB.Exec(`
        UPDATE reviews SET rating = $1, review = $2 
        WHERE id = $3`, rev.Rating, rev.Review, reviewID)
    if err != nil {
        http.Error(w, "Database update error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    rowsAffected, err := result.RowsAffected()
    if err != nil || rowsAffected == 0 {
        http.Error(w, "Review not found", http.StatusNotFound)
        return
    }
    rev.ID = reviewID
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(rev)
}


func DeleteReviewHandler(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    reviewID, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid review ID", http.StatusBadRequest)
        return
    }

    var orig Review
    err = db.DB.QueryRow(`
        SELECT id, user_id FROM reviews WHERE id = $1`, reviewID).
        Scan(&orig.ID, &orig.UserId)
    if err == sql.ErrNoRows {
        http.Error(w, "Review not found", http.StatusNotFound)
        return
    } else if err != nil {
        http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
        return
    }


    userIDStr := r.Header.Get("X-User-ID")
    reqUserID, err := strconv.Atoi(userIDStr)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusUnauthorized)
        return
    }


    if reqUserID != orig.UserId {
        http.Error(w, "You can only delete your own reviews", http.StatusForbidden)
        return
    }

    result, err := db.DB.Exec(`DELETE FROM reviews WHERE id = $1`, reviewID)
    if err != nil {
        http.Error(w, "Database delete error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    rowsAffected, err := result.RowsAffected()
    if err != nil || rowsAffected == 0 {
        http.Error(w, "Review not found", http.StatusNotFound)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}