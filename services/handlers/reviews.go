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

// Review represents a review record.
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

// allowedRating validates that the rating is exactly one of 1,2,3,4,5.
func allowedRating(rating float32) bool {
    intRating := int(rating)
    if float32(intRating) != rating {
        return false
    }
    return intRating >= 1 && intRating <= 5
}

// GetReviewsHandler retrieves reviews from the database based on optional query parameters.
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

// CreateReviewHandler inserts a new review record into the database.
func CreateReviewHandler(w http.ResponseWriter, r *http.Request) {
    var rev Review
    if err := json.NewDecoder(r.Body).Decode(&rev); err != nil {
        http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
        return
    }

    // Validate required rating: must be one of 1,2,3,4,5.
    if !allowedRating(rev.Rating) {
        http.Error(w, "Rating must be an integer between 1 and 5", http.StatusBadRequest)
        return
    }

    // Determine target: exactly one of CoffeeId, RoasteryId, or CoffeeShopId must be non-zero.
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

// UpdateReviewHandler updates an existing review record.
// It checks that the requesting user owns the review.
func UpdateReviewHandler(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    reviewID, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid review ID", http.StatusBadRequest)
        return
    }

    // Retrieve the original review to check ownership.
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

    // Assume user id is provided in header "X-User-ID".
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
    // Maintain original review's user_id and date_of_creation.
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

// DeleteReviewHandler deletes a review record from the database.
// It checks that the requesting user is the owner or is an admin.
func DeleteReviewHandler(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    reviewID, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid review ID", http.StatusBadRequest)
        return
    }

    // Retrieve the original review to check ownership.
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

    // Assume user id is provided in header "X-User-ID".
    userIDStr := r.Header.Get("X-User-ID")
    reqUserID, err := strconv.Atoi(userIDStr)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusUnauthorized)
        return
    }

    // Here you would typically check if the user has admin privileges.
    // For simplicity, we assume only the owner can delete.
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