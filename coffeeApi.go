package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

type Coffee struct {
    ID          int      `json:"id"`
    Name        string   `json:"name"`
    RoasteryId  int      `json:"roasteryId"`
    Country     string   `json:"country"`
    Region      string   `json:"region"`
    Farm        string   `json:"farm"`
    Variety     string   `json:"variety"`
    Process     string   `json:"process"`
    RoastProfile string  `json:"roastProfile"`
    FlavourNotes []string `json:"flavourNotes"`
    Description string   `json:"description"`
}

type Roastery struct {
    ID         int     `json:"id"`
    Name       string  `json:"name"`
    Country    string  `json:"country"`
    City       string  `json:"city"`
    Address    string  `json:"address"`
    Website    string  `json:"website"`
    Description string `json:"description"`
    AvgRating  float32 `json:"avgRating"`
}

type User struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Password string `json:"password"`
    Email    string `json:"email"`
    Role     string `json:"role"`
}

type CoffeeShop struct {
    ID         int     `json:"id"`
    Name       string  `json:"name"`
    Country    string  `json:"country"`
    City       string  `json:"city"`
    Address    string  `json:"address"`
    Website    string  `json:"website"`
    Description string `json:"description"`
    AvgRating  float32 `json:"avgRating"`
}

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

var (
    users      []User
    coffees    []Coffee
    roasteries []Roastery
    shops      []CoffeeShop
    reviews    []Review

    userID     = 1
    coffeeID   = 1
    roasteryID = 1
    shopID     = 1
    reviewID   = 1
    jwtKey     = []byte("super_sekretny_klucz_kofola_5mlnZł")
)

// autentykacja pośrednia (middleware)
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authentication")
        if authHeader == "" {
            http.Error(w, "Authorization header required", http.StatusUnauthorized)
            return
        }

        bearerToken := strings.Split(authHeader, " ")
        if len(bearerToken) != 2 {
            http.Error(w, "Invalid token format", http.StatusUnauthorized)
            return
        }

        token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
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
            http.Error(w, "Invalid user Id in token", http.StatusUnauthorized)
            return
        }

        r.Header.Set("X-User-ID", strconv.Itoa(int(userIDFloat)))

        next.ServeHTTP(w, r)
    })
}

// Admin posrednio
func adminMiddleware(next http.Handler) http.Handler {
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

        var user User
        for _, u := range users {
            if u.ID == userID {
                user = u
                break
            }
        }

        if user.Role != "admin" {
            http.Error(w, "Admin privilages required", http.StatusForbidden)
            return
        }

        next.ServeHTTP(w, r)
    })
}

// rejestracja 
func register(w http.ResponseWriter, r *http.Request) {
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    for _, u := range users {
        if u.Username == user.Username {
            http.Error(w, "Username already exist", http.StatusBadRequest)
            return
        }
    }

    // ustawianie id usera i domyslna role (user)
    user.ID = userID
    userID++

    if user.Role == "" {
        user.Role = "user"
    }

    // hasło do hashowania ??!?!?
    users = append(users, user)

    user.Password = "" // ukrywamy 🐗
    json.NewEncoder(w).Encode(user)
}

func login(w http.ResponseWriter, r *http.Request) {
    var credentials struct {
        Username  string `json:"username"`
        Passwords string `json:"passwords"`
    }

    if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // sprawdzanko czy user istnieje
    var user User
    found := false
    for _, u := range users {
        if u.Username == credentials.Username && u.Password == credentials.Passwords {
            user = u
            found = true
            break
        }
    }

    if !found {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    // jwt token (pomocy)
    expirationTime := time.Now().Add(24 * time.Hour)
    claims := jwt.MapClaims{
        "userId": user.ID,
        "exp":    expirationTime.Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        http.Error(w, "Failed to generate token", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{
        "token": tokenString,
    })
}

func getCoffees(w http.ResponseWriter, r *http.Request) {
    country := r.URL.Query().Get("country")
    process := r.URL.Query().Get("process")
    flavor := r.URL.Query().Get("flavor")
    roastery := r.URL.Query().Get("roastery")
    roastProfile := r.URL.Query().Get("roastProfile")
    name := r.URL.Query().Get("name")

    filteredCoffees := []Coffee{}

    for _, coffee := range coffees {
        if country != "" && !strings.EqualFold(coffee.Country, country) {
            continue
        }

        if process != "" && !strings.EqualFold(coffee.Process, process) {
            continue
        }

        if flavor != "" {
            hasFlavorNote := false
            for _, note := range coffee.FlavourNotes {
                if strings.Contains(strings.ToLower(note), strings.ToLower(flavor)) {
                    hasFlavorNote = true
                    break
                }
            }
            if !hasFlavorNote {
                continue
            }
        }

        if roastProfile != "" && strings.EqualFold(coffee.RoastProfile, roastProfile) {
            continue
        }

        if name != "" && !strings.Contains(strings.ToLower(coffee.Name), strings.ToLower(name)) {
            continue
        }

        if roastery != "" {
            roasteryFound := false
            roasteryID, err := strconv.Atoi(roastery)
            if err == nil {
                roasteryFound = coffee.RoasteryId == roasteryID
            } else {
                for _, r := range roasteries {
                    if r.ID == coffee.RoasteryId && strings.Contains(strings.ToLower(r.Name), strings.ToLower(roastery)) {
                        roasteryFound = true
                        break
                    }
                }
            }
            if !roasteryFound {
                continue
            }
        }

        filteredCoffees = append(filteredCoffees, coffee)
    }

    json.NewEncoder(w).Encode(filteredCoffees)
}

func getCoffee(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "invalid coffee id", http.StatusBadRequest)
        return
    }

    for _, coffee := range coffees {
        if coffee.ID == id {
            json.NewEncoder(w).Encode(coffee)
            return
        }
    }

    http.Error(w, "Coffee not found", http.StatusNotFound)
}

func createCoffee(w http.ResponseWriter, r *http.Request) {
    var coffee Coffee
    if err := json.NewDecoder(r.Body).Decode(&coffee); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if coffee.Name == "" || coffee.Country == "" || coffee.Process == "" || coffee.RoastProfile == "" {
        http.Error(w, "Missing required fields", http.StatusBadRequest)
        return
    }

    roasteryExists := false
    for _, r := range roasteries {
        if r.ID == coffee.RoasteryId {
            roasteryExists = true
            break
        }
    }

    if !roasteryExists {
        http.Error(w, "Roastery do not exist", http.StatusBadRequest)
        return
    }

    coffee.ID = coffeeID
    coffeeID++

    coffees = append(coffees, coffee)
    json.NewEncoder(w).Encode(coffee)
}

func updateCoffee(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid coffee ID", http.StatusBadRequest)
        return
    }

    var updatedCoffee Coffee
    if err := json.NewDecoder(r.Body).Decode(&updatedCoffee); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    for i, coffee := range coffees {
        if coffee.ID == id {
            updatedCoffee.ID = coffee.ID

            if updatedCoffee.RoasteryId != coffee.RoasteryId {
                roasteryExists := false
                for _, r := range roasteries {
                    if r.ID == updatedCoffee.RoasteryId {
                        roasteryExists = true
                        break
                    }
                }
                if !roasteryExists {
                    http.Error(w, "Roastery does not exist", http.StatusBadRequest)
                    return
                }
            }

            coffees[i] = updatedCoffee
            json.NewEncoder(w).Encode(updatedCoffee)
            return
        }
    }

    http.Error(w, "Coffee not found", http.StatusNotFound)
}

func deleteCoffee(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid coffee ID", http.StatusBadRequest)
        return
    }

    for i, coffee := range coffees {
        if coffee.ID == id {
            coffees = append(coffees[:i], coffees[i+1:]...) // usunięcie tej kawuuusi

            // pozbywamy sie recenzji
            var remainingReviews []Review
            for _, review := range reviews {
                if review.CoffeeId != id {
                    remainingReviews = append(remainingReviews, review)
                }
            }
            reviews = remainingReviews

            w.WriteHeader(http.StatusNoContent)
            return
        }
    }

    http.Error(w, "Coffee not found", http.StatusNotFound)
}

func getRoasteries(w http.ResponseWriter, r *http.Request) {
    name := r.URL.Query().Get("name")
    country := r.URL.Query().Get("country")
    city := r.URL.Query().Get("city")

    filteredRoasteries := []Roastery{}
    for _, roastery := range roasteries {
        if name != "" && !strings.Contains(strings.ToLower(roastery.Name), strings.ToLower(name)) {
            continue
        }
        if country != "" && !strings.EqualFold(roastery.Country, country) {
            continue
        }
        if city != "" && !strings.EqualFold(roastery.City, city) {
            continue
        }

        filteredRoasteries = append(filteredRoasteries, roastery)
    }

    json.NewEncoder(w).Encode(filteredRoasteries)
}

func getRoastery(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid roastery ID", http.StatusBadRequest)
        return
    }

    for _, roastery := range roasteries {
        if roastery.ID == id {
            json.NewEncoder(w).Encode(roastery)
            return
        }
    }

    http.Error(w, "Roastery not found", http.StatusNotFound)
}

func createRoastery(w http.ResponseWriter, r *http.Request) {
    var roastery Roastery
    if err := json.NewDecoder(r.Body).Decode(&roastery); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if roastery.Name == "" || roastery.Country == "" || roastery.City == "" {
        http.Error(w, "Missing required fields", http.StatusBadRequest)
        return
    }

    roastery.ID = roasteryID
    roasteryID++
    roastery.AvgRating = 0

    roasteries = append(roasteries, roastery)
    json.NewEncoder(w).Encode(roastery)
}

func updateRoastery(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid roastery ID", http.StatusBadRequest)
        return
    }

    var updatedRoastery Roastery
    if err := json.NewDecoder(r.Body).Decode(&updatedRoastery); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    for i, roastery := range roasteries {
        if roastery.ID == id {
            updatedRoastery.ID = roastery.ID
            updatedRoastery.AvgRating = roastery.AvgRating

            roasteries[i] = updatedRoastery
            json.NewEncoder(w).Encode(updatedRoastery)
            return
        }
    }

    http.Error(w, "Roastery not found", http.StatusNotFound)
}

func deleteRoastery(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid roastery ID", http.StatusBadRequest)
        return
    }

    for _, coffee := range coffees {
        if coffee.RoasteryId == id {
            http.Error(w, "Cannot delete roastery that has coffees", http.StatusBadRequest)
            return
        }
    }

    for i, roastery := range roasteries {
        if roastery.ID == id {
            roasteries = append(roasteries[:i], roasteries[i+1:]...)

            var remainingReviews []Review
            for _, review := range reviews {
                if review.RoasteryId != id {
                    remainingReviews = append(remainingReviews, review)
                }
            }
            reviews = remainingReviews

            w.WriteHeader(http.StatusNoContent)
            return
        }
    }

    http.Error(w, "Roastery not found", http.StatusNotFound)
}

func getCoffeeShops(w http.ResponseWriter, r *http.Request) {
    name := r.URL.Query().Get("name")
    country := r.URL.Query().Get("country")
    city := r.URL.Query().Get("city")

    filteredShops := []CoffeeShop{}
    for _, shop := range shops {
        if name != "" && !strings.Contains(strings.ToLower(shop.Name), strings.ToLower(name)) {
            continue
        }
        if country != "" && !strings.EqualFold(shop.Country, country) {
            continue
        }
        if city != "" && !strings.EqualFold(shop.City, city) {
            continue
        }

        filteredShops = append(filteredShops, shop)
    }

    json.NewEncoder(w).Encode(filteredShops)
}

func getCoffeeShop(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid shop ID", http.StatusBadRequest)
        return
    }

    for _, shop := range shops {
        if shop.ID == id {
            json.NewEncoder(w).Encode(shop)
            return
        }
    }

    http.Error(w, "Coffee shop not found", http.StatusNotFound)
}

func createCoffeeShop(w http.ResponseWriter, r *http.Request) {
    var shop CoffeeShop
    if err := json.NewDecoder(r.Body).Decode(&shop); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if shop.Name == "" || shop.Country == "" || shop.City == "" {
        http.Error(w, "Missing required fields", http.StatusBadRequest)
        return
    }

    shop.ID = shopID
    shopID++
    shop.AvgRating = 0

    shops = append(shops, shop)
    json.NewEncoder(w).Encode(shop)
}

func updateCoffeeShop(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid shop ID", http.StatusBadRequest)
        return
    }

    var updatedShop CoffeeShop
    if err := json.NewDecoder(r.Body).Decode(&updatedShop); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    for i, shop := range shops {
        if shop.ID == id {
            updatedShop.ID = shop.ID
            updatedShop.AvgRating = shop.AvgRating

            shops[i] = updatedShop
            json.NewEncoder(w).Encode(updatedShop)
            return
        }
    }

    http.Error(w, "Coffee shop not found", http.StatusNotFound)
}

func deleteCoffeeShop(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid shop ID", http.StatusBadRequest)
        return
    }

    for i, shop := range shops {
        if shop.ID == id {
            shops = append(shops[:i], shops[i+1:]...)

            var remainingReviews []Review
            for _, review := range reviews {
                if review.CoffeeShopId != id {
                    remainingReviews = append(remainingReviews, review)
                }
            }
            reviews = remainingReviews

            w.WriteHeader(http.StatusNoContent)
            return
        }
    }

    http.Error(w, "Coffee shop not found", http.StatusNotFound)
}

func getReviews(w http.ResponseWriter, r *http.Request) {
    coffeeIDStr := r.URL.Query().Get("coffeeId")
    roasteryIDStr := r.URL.Query().Get("roasteryId")
    shopIDStr := r.URL.Query().Get("shopId")
    userIDStr := r.URL.Query().Get("userId")

    filteredReviews := []Review{}
    for _, review := range reviews {
        if coffeeIDStr != "" {
            id, err := strconv.Atoi(coffeeIDStr)
            if err != nil || review.CoffeeId != id {
                continue
            }
        }
        if roasteryIDStr != "" {
            id, err := strconv.Atoi(roasteryIDStr)
            if err != nil || review.RoasteryId != id {
                continue
            }
        }
        if shopIDStr != "" {
            id, err := strconv.Atoi(shopIDStr)
            if err != nil || review.CoffeeShopId != id {
                continue
            }
        }
        if userIDStr != "" {
            id, err := strconv.Atoi(userIDStr)
            if err != nil || review.UserId != id {
                continue
            }
        }

        filteredReviews = append(filteredReviews, review)
    }

    json.NewEncoder(w).Encode(filteredReviews)
}

func createReview(w http.ResponseWriter, r *http.Request) {
    var review Review
    if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    userIDStr := r.Header.Get("X-User-ID")
    if userIDStr == "" {
        http.Error(w, "You must be authenticated to create a review", http.StatusUnauthorized)
        return
    }

    userID, err := strconv.Atoi(userIDStr)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusUnauthorized)
        return
    }

    review.UserId = userID

    // Validate required fields
    if review.Rating < 0 || review.Rating > 10 {
        http.Error(w, "Rating must be between 0 and 10", http.StatusBadRequest)
        return
    }

    targetCount := 0
    if review.CoffeeId != 0 {
        targetCount++
        coffeeExists := false
        for _, coffee := range coffees {
            if coffee.ID == review.CoffeeId {
                coffeeExists = true
                break
            }
        }
        if !coffeeExists {
            http.Error(w, "Coffee does not exist", http.StatusBadRequest)
            return
        }
    }
    if review.RoasteryId != 0 {
        targetCount++
        roasteryExists := false
        for _, roastery := range roasteries {
            if roastery.ID == review.RoasteryId {
                roasteryExists = true
                break
            }
        }
        if !roasteryExists {
            http.Error(w, "Roastery does not exist", http.StatusBadRequest)
            return
        }
    }
    if review.CoffeeShopId != 0 {
        targetCount++
        shopExists := false
        for _, shop := range shops {
            if shop.ID == review.CoffeeShopId {
                shopExists = true
                break
            }
        }
        if !shopExists {
            http.Error(w, "Coffee shop does not exist", http.StatusBadRequest)
            return
        }
    }

    if targetCount != 1 {
        http.Error(w, "Review must target exactly one of: coffee, roastery, or coffee shop", http.StatusBadRequest)
        return
    }

    review.ID = reviewID
    reviewID++
    review.DateOfCreation = time.Now()

    reviews = append(reviews, review)

    updateAverageRatings()

    json.NewEncoder(w).Encode(review)
}

func updateReview(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid review ID", http.StatusBadRequest)
        return
    }

    var updatedReview Review
    if err := json.NewDecoder(r.Body).Decode(&updatedReview); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    userIDStr := r.Header.Get("X-User-ID")
    if userIDStr == "" {
        http.Error(w, "You must be authenticated to update a review", http.StatusUnauthorized)
        return
    }

    userID, err := strconv.Atoi(userIDStr)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusUnauthorized)
        return
    }

    for i, review := range reviews {
        if review.ID == id {
            if review.UserId != userID {
                http.Error(w, "You can only update your own reviews", http.StatusForbidden)
                return
            }

            if updatedReview.Rating < 0 || updatedReview.Rating > 10 {
                http.Error(w, "Rating must be between 0 and 10", http.StatusBadRequest)
                return
            }

            updatedReview.ID = review.ID
            updatedReview.UserId = review.UserId
            updatedReview.CoffeeId = review.CoffeeId
            updatedReview.RoasteryId = review.RoasteryId
            updatedReview.CoffeeShopId = review.CoffeeShopId
            updatedReview.DateOfCreation = review.DateOfCreation

            reviews[i] = updatedReview

            updateAverageRatings()

            json.NewEncoder(w).Encode(updatedReview)
            return
        }
    }

    http.Error(w, "Review not found", http.StatusNotFound)
}

func deleteReview(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid review ID", http.StatusBadRequest)
        return
    }

    userIDStr := r.Header.Get("X-User-ID")
    if userIDStr == "" {
        http.Error(w, "You must be authenticated to delete a review", http.StatusUnauthorized)
        return
    }

    userID, err := strconv.Atoi(userIDStr)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusUnauthorized)
        return
    }

    isAdmin := false
    for _, u := range users {
        if u.ID == userID && u.Role == "admin" {
            isAdmin = true
            break
        }
    }

    for i, review := range reviews {
        if review.ID == id {
            if review.UserId != userID && !isAdmin {
                http.Error(w, "You can only delete your own reviews", http.StatusForbidden)
                return
            }

            reviews = append(reviews[:i], reviews[i+1:]...)

            updateAverageRatings()

            w.WriteHeader(http.StatusNoContent)
            return
        }
    }

    http.Error(w, "Review not found", http.StatusNotFound)
}

func updateAverageRatings() {
    for i := range shops {
        var sum float64
        count := 0
        for _, review := range reviews {
            if review.CoffeeShopId == shops[i].ID {
                sum += float64(review.Rating)
                count++
            }
        }
        if count > 0 {
            shops[i].AvgRating = float32(sum / float64(count))
        } else {
            shops[i].AvgRating = 0
        }
    }

    for i := range roasteries {
        var sum float64
        count := 0
        for _, review := range reviews {
            if review.RoasteryId == roasteries[i].ID {
                sum += float64(review.Rating)
                count++
            }
        }
        if count > 0 {
            roasteries[i].AvgRating = float32(sum / float64(count))
        } else {
            roasteries[i].AvgRating = 0
        }
    }
}

func seedData() {
    file, err := os.Open("data.json")
    if err != nil {
        fmt.Println("Błąd podczas otwierania pliku data.json:", err)
        return
    }
    defer file.Close()

    var data struct {
        Users      []User       `json:"users"`
        Coffees    []Coffee     `json:"coffees"`
        Roasteries []Roastery   `json:"roasteries"`
        Shops      []CoffeeShop `json:"shops"`
        Reviews    []Review     `json:"reviews"`
    }

    decoder := json.NewDecoder(file)
    if err := decoder.Decode(&data); err != nil {
        fmt.Println("Błąd podczas dekodowania JSON:", err)
        return
    }

    users = data.Users
    coffees = data.Coffees
    roasteries = data.Roasteries
    shops = data.Shops
    reviews = data.Reviews

    maxUserID := 0
    maxCoffeeID := 0
    maxRoasteryID := 0
    maxShopID := 0
    maxReviewID := 0

    for _, user := range users {
        if user.ID > maxUserID {
            maxUserID = user.ID
        }
    }

    for _, coffee := range coffees {
        if coffee.ID > maxCoffeeID {
            maxCoffeeID = coffee.ID
        }
    }

    for _, roastery := range roasteries {
        if roastery.ID > maxRoasteryID {
            maxRoasteryID = roastery.ID
        }
    }

    for _, shop := range shops {
        if shop.ID > maxShopID {
            maxShopID = shop.ID
        }
    }

    for _, review := range reviews {
        if review.ID > maxReviewID {
            maxReviewID = review.ID
        }
    }

    userID = maxUserID + 1
    coffeeID = maxCoffeeID + 1
    roasteryID = maxRoasteryID + 1
    shopID = maxShopID + 1
    reviewID = maxReviewID + 1

    updateAverageRatings()

    fmt.Println("Dane zostały wczytane pomyślnie!")
    fmt.Printf("Załadowano: %d użytkowników, %d kaw, %d palarni, %d kawiarni, %d recenzji\n",
        len(users), len(coffees), len(roasteries), len(shops), len(reviews))
}

func main() {
    seedData()

    router := mux.NewRouter()

    router.HandleFunc("/register", register).Methods("POST")
    router.HandleFunc("/login", login).Methods("POST")

    router.HandleFunc("/coffees", getCoffees).Methods("GET")
    router.HandleFunc("/coffees/{id}", getCoffee).Methods("GET")
    router.HandleFunc("/coffees", createCoffee).Methods("POST").Handler(authMiddleware(http.HandlerFunc(createCoffee)))
    router.HandleFunc("/coffees/{id}", updateCoffee).Methods("PUT").Handler(authMiddleware(http.HandlerFunc(updateCoffee)))
    router.HandleFunc("/coffees/{id}", deleteCoffee).Methods("DELETE").Handler(authMiddleware(adminMiddleware(http.HandlerFunc(deleteCoffee))))

    router.HandleFunc("/roasteries", getRoasteries).Methods("GET")
    router.HandleFunc("/roasteries/{id}", getRoastery).Methods("GET")
    router.HandleFunc("/roasteries", createRoastery).Methods("POST").Handler(authMiddleware(http.HandlerFunc(createRoastery)))
    router.HandleFunc("/roasteries/{id}", updateRoastery).Methods("PUT").Handler(authMiddleware(http.HandlerFunc(updateRoastery)))
    router.HandleFunc("/roasteries/{id}", deleteRoastery).Methods("DELETE").Handler(authMiddleware(adminMiddleware(http.HandlerFunc(deleteRoastery))))

    router.HandleFunc("/shops", getCoffeeShops).Methods("GET")
    router.HandleFunc("/shops/{id}", getCoffeeShop).Methods("GET")
    router.HandleFunc("/shops", createCoffeeShop).Methods("POST").Handler(authMiddleware(http.HandlerFunc(createCoffeeShop)))
    router.HandleFunc("/shops/{id}", updateCoffeeShop).Methods("PUT").Handler(authMiddleware(http.HandlerFunc(updateCoffeeShop)))
    router.HandleFunc("/shops/{id}", deleteCoffeeShop).Methods("DELETE").Handler(authMiddleware(adminMiddleware(http.HandlerFunc(deleteCoffeeShop))))

    router.HandleFunc("/reviews", getReviews).Methods("GET")
    router.HandleFunc("/reviews", createReview).Methods("POST").Handler(authMiddleware(http.HandlerFunc(createReview)))
    router.HandleFunc("/reviews/{id}", updateReview).Methods("PUT").Handler(authMiddleware(http.HandlerFunc(updateReview)))
    router.HandleFunc("/reviews/{id}", deleteReview).Methods("DELETE").Handler(authMiddleware(http.HandlerFunc(deleteReview)))

    port := ":40331"
    fmt.Printf("Serwer uruchomiony na porcie %s\n", port)
    err := http.ListenAndServe(port, router)
    if err != nil {
        fmt.Println("Błąd podczas uruchamiania serwera:", err)
    }
}


