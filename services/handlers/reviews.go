package handlers

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"
    "strings"
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

type ReviewResponse struct {
    ID             int       `json:"id"`
    UserId         int       `json:"userId"`
    UserName       string    `json:"userName"`
    CoffeeId       int       `json:"coffeeId"`
    CoffeeName     string    `json:"coffeeName,omitempty"`
    RoasteryId     int       `json:"roasteryId"`
    RoasteryName   string    `json:"roasteryName,omitempty"`
    CoffeeShopId   int       `json:"coffeeShopId"`
    CoffeeShopName string    `json:"coffeeShopName,omitempty"`
    Rating         float32   `json:"rating"`
    Review         string    `json:"review"`
    DateOfCreation time.Time `json:"dateOfCreation"`
    TargetType     string    `json:"targetType"`
    TargetName     string    `json:"targetName"`
}

func allowedRating(rating float32) bool {
    intRating := int(rating)
    if float32(intRating) != rating {
        return false
    }
    return intRating >= 1 && intRating <= 5
}

func GetReviewsHandler(w http.ResponseWriter, r *http.Request) {
    q := r.URL.Query()
    
    userId := q.Get("userId")
    coffeeId := q.Get("coffeeId")
    roasteryId := q.Get("roasteryId")
    coffeeShopId := q.Get("coffeeShopId")
    minRating := q.Get("minRating")
    maxRating := q.Get("maxRating")
    fromDate := q.Get("fromDate")
    toDate := q.Get("toDate")
    
    coffeeCountry := q.Get("coffeeCountry")
    coffeeProcess := q.Get("coffeeProcess")
    coffeeRoastProfile := q.Get("coffeeRoastProfile")
    coffeeFlavour := q.Get("coffeeFlavour")
    
    roasteryCountry := q.Get("roasteryCountry")
    roasteryCity := q.Get("roasteryCity")
    
    shopCountry := q.Get("shopCountry")
    shopCity := q.Get("shopCity")
    
    baseQuery := `
        SELECT r.id, r.user_id, r.coffee_id, r.roastery_id, r.coffee_shop_id, r.rating, r.review, r.date_of_creation,
               u.username AS user_name,
               c.name AS coffee_name,
               ro.name AS roastery_name,
               s.name AS shop_name
        FROM reviews r
        LEFT JOIN users u ON r.user_id = u.id
        LEFT JOIN coffees c ON r.coffee_id = c.id
        LEFT JOIN roasteries ro ON r.roastery_id = ro.id
        LEFT JOIN shops s ON r.coffee_shop_id = s.id`
    
    conditions := []string{}
    args := []interface{}{}
    argIdx := 1
    
    if userId != "" {
        if id, err := strconv.Atoi(userId); err == nil {
            conditions = append(conditions, fmt.Sprintf("r.user_id = $%d", argIdx))
            args = append(args, id)
            argIdx++
        }
    }
    if coffeeId != "" {
        if id, err := strconv.Atoi(coffeeId); err == nil {
            conditions = append(conditions, fmt.Sprintf("r.coffee_id = $%d", argIdx))
            args = append(args, id)
            argIdx++
        }
    }
    if roasteryId != "" {
        if id, err := strconv.Atoi(roasteryId); err == nil {
            conditions = append(conditions, fmt.Sprintf("r.roastery_id = $%d", argIdx))
            args = append(args, id)
            argIdx++
        }
    }
    if coffeeShopId != "" {
        if id, err := strconv.Atoi(coffeeShopId); err == nil {
            conditions = append(conditions, fmt.Sprintf("r.coffee_shop_id = $%d", argIdx))
            args = append(args, id)
            argIdx++
        }
    }
    if minRating != "" {
        if rating, err := strconv.ParseFloat(minRating, 32); err == nil {
            conditions = append(conditions, fmt.Sprintf("r.rating >= $%d", argIdx))
            args = append(args, float32(rating))
            argIdx++
        }
    }
    if maxRating != "" {
        if rating, err := strconv.ParseFloat(maxRating, 32); err == nil {
            conditions = append(conditions, fmt.Sprintf("r.rating <= $%d", argIdx))
            args = append(args, float32(rating))
            argIdx++
        }
    }
    if fromDate != "" {
        date, err := time.Parse("2006-01-02", fromDate)
        if err == nil {
            conditions = append(conditions, fmt.Sprintf("r.date_of_creation >= $%d", argIdx))
            args = append(args, date)
            argIdx++
        }
    }
    if toDate != "" {
        date, err := time.Parse("2006-01-02", toDate)
        if err == nil {
            conditions = append(conditions, fmt.Sprintf("r.date_of_creation <= $%d", argIdx))
            args = append(args, date.Add(24*time.Hour))
            argIdx++
        }
    }
    
    if coffeeCountry != "" {
        conditions = append(conditions, fmt.Sprintf("c.country ILIKE $%d", argIdx))
        args = append(args, "%"+coffeeCountry+"%")
        argIdx++
    }
    if coffeeProcess != "" {
        conditions = append(conditions, fmt.Sprintf("c.process ILIKE $%d", argIdx))
        args = append(args, "%"+coffeeProcess+"%")
        argIdx++
    }
    if coffeeRoastProfile != "" {
        conditions = append(conditions, fmt.Sprintf("c.roast_profile ILIKE $%d", argIdx))
        args = append(args, "%"+coffeeRoastProfile+"%")
        argIdx++
    }
    if coffeeFlavour != "" {
        conditions = append(conditions, fmt.Sprintf("c.flavour_notes ILIKE $%d", argIdx))
        args = append(args, "%"+coffeeFlavour+"%")
        argIdx++
    }
    
    if roasteryCountry != "" {
        conditions = append(conditions, fmt.Sprintf("ro.country ILIKE $%d", argIdx))
        args = append(args, "%"+roasteryCountry+"%")
        argIdx++
    }
    if roasteryCity != "" {
        conditions = append(conditions, fmt.Sprintf("ro.city ILIKE $%d", argIdx))
        args = append(args, "%"+roasteryCity+"%")
        argIdx++
    }
    
    if shopCountry != "" {
        conditions = append(conditions, fmt.Sprintf("s.country ILIKE $%d", argIdx))
        args = append(args, "%"+shopCountry+"%")
        argIdx++
    }
    if shopCity != "" {
        conditions = append(conditions, fmt.Sprintf("s.city ILIKE $%d", argIdx))
        args = append(args, "%"+shopCity+"%")
        argIdx++
    }

    if len(conditions) > 0 {
        baseQuery += " WHERE " + strings.Join(conditions, " AND ")
    }

    baseQuery += " ORDER BY r.date_of_creation DESC"

    rows, err := db.DB.Query(baseQuery, args...)
    if err != nil {
        http.Error(w, "Database query error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var reviewResponses []ReviewResponse
    for rows.Next() {
        var rev Review
        var userName, coffeeName, roasteryName, shopName sql.NullString
        
        if err := rows.Scan(
            &rev.ID, &rev.UserId, &rev.CoffeeId, &rev.RoasteryId, &rev.CoffeeShopId, 
            &rev.Rating, &rev.Review, &rev.DateOfCreation,
            &userName, &coffeeName, &roasteryName, &shopName); err != nil {
            http.Error(w, "Error scanning review: "+err.Error(), http.StatusInternalServerError)
            return
        }
        
        response := ReviewResponse{
            ID:             rev.ID,
            UserId:         rev.UserId,
            UserName:       nullStringValue(userName, "Anonymous User"),
            CoffeeId:       rev.CoffeeId,
            CoffeeName:     nullStringValue(coffeeName, ""),
            RoasteryId:     rev.RoasteryId,
            RoasteryName:   nullStringValue(roasteryName, ""),
            CoffeeShopId:   rev.CoffeeShopId,
            CoffeeShopName: nullStringValue(shopName, ""),
            Rating:         rev.Rating,
            Review:         rev.Review,
            DateOfCreation: rev.DateOfCreation,
            TargetType:     getTargetType(rev.CoffeeId, rev.RoasteryId, rev.CoffeeShopId),
            TargetName:     getTargetName(coffeeName, roasteryName, shopName),
        }
        
        reviewResponses = append(reviewResponses, response)
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(reviewResponses)
}

func GetReviewHandler(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    reviewID, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid review ID", http.StatusBadRequest)
        return
    }
    
    var rev Review
    var userName, coffeeName, roasteryName, shopName sql.NullString
    
    err = db.DB.QueryRow(`
        SELECT r.id, r.user_id, r.coffee_id, r.roastery_id, r.coffee_shop_id, r.rating, r.review, r.date_of_creation,
               u.username AS user_name,
               c.name AS coffee_name,
               ro.name AS roastery_name,
               s.name AS shop_name
        FROM reviews r
        LEFT JOIN users u ON r.user_id = u.id
        LEFT JOIN coffees c ON r.coffee_id = c.id
        LEFT JOIN roasteries ro ON r.roastery_id = ro.id
        LEFT JOIN shops s ON r.coffee_shop_id = s.id
        WHERE r.id = $1`, reviewID).
        Scan(&rev.ID, &rev.UserId, &rev.CoffeeId, &rev.RoasteryId, &rev.CoffeeShopId, 
             &rev.Rating, &rev.Review, &rev.DateOfCreation,
             &userName, &coffeeName, &roasteryName, &shopName)
             
    if err == sql.ErrNoRows {
        http.Error(w, "Review not found", http.StatusNotFound)
        return
    } else if err != nil {
        http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    
    response := ReviewResponse{
        ID:             rev.ID,
        UserId:         rev.UserId,
        UserName:       nullStringValue(userName, "Anonymous User"),
        CoffeeId:       rev.CoffeeId,
        CoffeeName:     nullStringValue(coffeeName, ""),
        RoasteryId:     rev.RoasteryId,
        RoasteryName:   nullStringValue(roasteryName, ""),
        CoffeeShopId:   rev.CoffeeShopId,
        CoffeeShopName: nullStringValue(shopName, ""),
        Rating:         rev.Rating,
        Review:         rev.Review,
        DateOfCreation: rev.DateOfCreation,
        TargetType:     getTargetType(rev.CoffeeId, rev.RoasteryId, rev.CoffeeShopId),
        TargetName:     getTargetName(coffeeName, roasteryName, shopName),
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func CreateReviewHandler(w http.ResponseWriter, r *http.Request) {
    var rev Review
    if err := json.NewDecoder(r.Body).Decode(&rev); err != nil {
        http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
        return
    }

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
    rev.UserId = userID

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
    err = db.DB.QueryRow(`
        INSERT INTO reviews (user_id, coffee_id, roastery_id, coffee_shop_id, rating, review, date_of_creation)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id`,
        rev.UserId, rev.CoffeeId, rev.RoasteryId, rev.CoffeeShopId, rev.Rating, rev.Review, rev.DateOfCreation).Scan(&rev.ID)
    if err != nil {
        http.Error(w, "Database insert error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    updateAverageRating(rev.CoffeeId, rev.RoasteryId, rev.CoffeeShopId)

    var userName, coffeeName, roasteryName, shopName sql.NullString
    
    err = db.DB.QueryRow(`
        SELECT r.id, r.user_id, r.coffee_id, r.roastery_id, r.coffee_shop_id, r.rating, r.review, r.date_of_creation,
               u.username AS user_name,
               c.name AS coffee_name,
               ro.name AS roastery_name,
               s.name AS shop_name
        FROM reviews r
        LEFT JOIN users u ON r.user_id = u.id
        LEFT JOIN coffees c ON r.coffee_id = c.id
        LEFT JOIN roasteries ro ON r.roastery_id = ro.id
        LEFT JOIN shops s ON r.coffee_shop_id = s.id
        WHERE r.id = $1`, rev.ID).
        Scan(&rev.ID, &rev.UserId, &rev.CoffeeId, &rev.RoasteryId, &rev.CoffeeShopId, 
             &rev.Rating, &rev.Review, &rev.DateOfCreation,
             &userName, &coffeeName, &roasteryName, &shopName)
             
    if err != nil {
        http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    
    response := ReviewResponse{
        ID:             rev.ID,
        UserId:         rev.UserId,
        UserName:       nullStringValue(userName, "Anonymous User"),
        CoffeeId:       rev.CoffeeId,
        CoffeeName:     nullStringValue(coffeeName, ""),
        RoasteryId:     rev.RoasteryId,
        RoasteryName:   nullStringValue(roasteryName, ""),
        CoffeeShopId:   rev.CoffeeShopId,
        CoffeeShopName: nullStringValue(shopName, ""),
        Rating:         rev.Rating,
        Review:         rev.Review,
        DateOfCreation: rev.DateOfCreation,
        TargetType:     getTargetType(rev.CoffeeId, rev.RoasteryId, rev.CoffeeShopId),
        TargetName:     getTargetName(coffeeName, roasteryName, shopName),
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
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
    userRoleStr := r.Header.Get("X-User-Role")
    reqUserID, err := strconv.Atoi(userIDStr)
    if err != nil || (reqUserID != orig.UserId && userRoleStr != "admin") {
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
    rev.CoffeeId = orig.CoffeeId
    rev.RoasteryId = orig.RoasteryId
    rev.CoffeeShopId = orig.CoffeeShopId
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
    
    updateAverageRating(rev.CoffeeId, rev.RoasteryId, rev.CoffeeShopId)
    
    var userName, coffeeName, roasteryName, shopName sql.NullString
    
    err = db.DB.QueryRow(`
        SELECT r.id, r.user_id, r.coffee_id, r.roastery_id, r.coffee_shop_id, r.rating, r.review, r.date_of_creation,
               u.username AS user_name,
               c.name AS coffee_name,
               ro.name AS roastery_name,
               s.name AS shop_name
        FROM reviews r
        LEFT JOIN users u ON r.user_id = u.id
        LEFT JOIN coffees c ON r.coffee_id = c.id
        LEFT JOIN roasteries ro ON r.roastery_id = ro.id
        LEFT JOIN shops s ON r.coffee_shop_id = s.id
        WHERE r.id = $1`, reviewID).
        Scan(&rev.ID, &rev.UserId, &rev.CoffeeId, &rev.RoasteryId, &rev.CoffeeShopId, 
             &rev.Rating, &rev.Review, &rev.DateOfCreation,
             &userName, &coffeeName, &roasteryName, &shopName)
             
    if err != nil {
        http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    
    response := ReviewResponse{
        ID:             rev.ID,
        UserId:         rev.UserId,
        UserName:       nullStringValue(userName, "Anonymous User"),
        CoffeeId:       rev.CoffeeId,
        CoffeeName:     nullStringValue(coffeeName, ""),
        RoasteryId:     rev.RoasteryId,
        RoasteryName:   nullStringValue(roasteryName, ""),
        CoffeeShopId:   rev.CoffeeShopId,
        CoffeeShopName: nullStringValue(shopName, ""),
        Rating:         rev.Rating,
        Review:         rev.Review,
        DateOfCreation: rev.DateOfCreation,
        TargetType:     getTargetType(rev.CoffeeId, rev.RoasteryId, rev.CoffeeShopId),
        TargetName:     getTargetName(coffeeName, roasteryName, shopName),
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
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
        SELECT id, user_id, coffee_id, roastery_id, coffee_shop_id FROM reviews WHERE id = $1`, reviewID).
        Scan(&orig.ID, &orig.UserId, &orig.CoffeeId, &orig.RoasteryId, &orig.CoffeeShopId)
    if err == sql.ErrNoRows {
        http.Error(w, "Review not found", http.StatusNotFound)
        return
    } else if err != nil {
        http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    userIDStr := r.Header.Get("X-User-ID")
    userRoleStr := r.Header.Get("X-User-Role")
    reqUserID, err := strconv.Atoi(userIDStr)
    
    if err != nil || (reqUserID != orig.UserId && userRoleStr != "admin") {
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
    
    updateAverageRating(orig.CoffeeId, orig.RoasteryId, orig.CoffeeShopId)
    
    w.WriteHeader(http.StatusNoContent)
}

func updateAverageRating(coffeeId, roasteryId, coffeeShopId int) {
    if coffeeId != 0 {
        _, err := db.DB.Exec(`
            UPDATE coffees SET avg_rating = 
            (SELECT COALESCE(AVG(rating), 0) FROM reviews WHERE coffee_id = $1)
            WHERE id = $1`, coffeeId)
        if err != nil {
            fmt.Printf("Error updating coffee rating: %v\n", err)
        }
    }
    
    if roasteryId != 0 {
        _, err := db.DB.Exec(`
            UPDATE roasteries SET avg_rating = 
            (SELECT COALESCE(AVG(rating), 0) FROM reviews WHERE roastery_id = $1)
            WHERE id = $1`, roasteryId)
        if err != nil {
            fmt.Printf("Error updating roastery rating: %v\n", err)
        }
    }
    
    if coffeeShopId != 0 {
        _, err := db.DB.Exec(`
            UPDATE shops SET avg_rating = 
            (SELECT COALESCE(AVG(rating), 0) FROM reviews WHERE coffee_shop_id = $1)
            WHERE id = $1`, coffeeShopId)
        if err != nil {
            fmt.Printf("Error updating coffee shop rating: %v\n", err)
        }
    }
}

func nullStringValue(ns sql.NullString, defaultValue string) string {
    if ns.Valid {
        return ns.String
    }
    return defaultValue
}

func getTargetType(coffeeId, roasteryId, coffeeShopId int) string {
    if coffeeId != 0 {
        return "coffee"
    }
    if roasteryId != 0 {
        return "roastery"
    }
    if coffeeShopId != 0 {
        return "coffee_shop"
    }
    return "unknown"
}

func getTargetName(coffeeName, roasteryName, shopName sql.NullString) string {
    if coffeeName.Valid {
        return coffeeName.String
    }
    if roasteryName.Valid {
        return roasteryName.String
    }
    if shopName.Valid {
        return shopName.String
    }
    return "Unknown"
}