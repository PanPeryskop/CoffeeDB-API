package handlers

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"

    "coffeeApi/services/db"
    "coffeeApi/services/geocoding"
    "github.com/gorilla/mux"
)

type CoffeeShop struct {
    ID          int     `json:"id"`
    Name        string  `json:"name"`
    Country     string  `json:"country"`
    City        string  `json:"city"`
    Address     string  `json:"address"`
    Website     string  `json:"website"`
    Description string  `json:"description"`
    AvgRating   float32 `json:"avgRating"`
    Lat         float64 `json:"lat"`
    Lon         float64 `json:"lon"`
}

func GetCoffeeShopsHandler(w http.ResponseWriter, r *http.Request) {
    rows, err := db.DB.Query(`
        SELECT id, name, country, city, address, website, description, avg_rating, lat, lon
        FROM shops`)
    if err != nil {
        http.Error(w, "Database query error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var shops []CoffeeShop
    for rows.Next() {
        var shop CoffeeShop
        if err := rows.Scan(&shop.ID, &shop.Name, &shop.Country, &shop.City, &shop.Address, &shop.Website, &shop.Description, &shop.AvgRating, &shop.Lat, &shop.Lon); err != nil {
            http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
            return
        }
        shops = append(shops, shop)
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(shops)
}

func GetCoffeeShopHandler(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    shopID, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid shop ID", http.StatusBadRequest)
        return
    }

    var shop CoffeeShop
    err = db.DB.QueryRow(`
        SELECT id, name, country, city, address, website, description, avg_rating, lat, lon
        FROM shops WHERE id = $1`, shopID).
        Scan(&shop.ID, &shop.Name, &shop.Country, &shop.City, &shop.Address, &shop.Website, &shop.Description, &shop.AvgRating, &shop.Lat, &shop.Lon)
    if err == sql.ErrNoRows {
        http.Error(w, "Coffee shop not found", http.StatusNotFound)
        return
    } else if err != nil {
        http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(shop)
}


func CreateCoffeeShopHandler(w http.ResponseWriter, r *http.Request) {
    var shop CoffeeShop
    if err := json.NewDecoder(r.Body).Decode(&shop); err != nil {
        http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
        return
    }
    if shop.Name == "" || shop.Country == "" || shop.City == "" || shop.Address == "" {
        http.Error(w, "Missing required fields", http.StatusBadRequest)
        return
    }

    // Construct full address from address, city and country.
    fullAddress := fmt.Sprintf("%s, %s, %s", shop.Address, shop.City, shop.Country)
    lat, lon, err := geocoding.GetCoordinates(fullAddress)
    if err != nil {
        http.Error(w, "Geocoding error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    shop.Lat = lat
    shop.Lon = lon


    err = db.DB.QueryRow(`
        INSERT INTO shops (name, country, city, address, website, description, avg_rating, lat, lon)
        VALUES ($1, $2, $3, $4, $5, $6, 0, $7, $8) RETURNING id`,
        shop.Name, shop.Country, shop.City, shop.Address, shop.Website, shop.Description, shop.Lat, shop.Lon).
        Scan(&shop.ID)
    if err != nil {
        http.Error(w, "Database insert error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    shop.AvgRating = 0
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(shop)
}


func UpdateCoffeeShopHandler(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    shopID, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid shop ID", http.StatusBadRequest)
        return
    }

    var shop CoffeeShop
    if err := json.NewDecoder(r.Body).Decode(&shop); err != nil {
        http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
        return
    }

    fullAddress := fmt.Sprintf("%s, %s, %s", shop.Address, shop.City, shop.Country)
    lat, lon, err := geocoding.GetCoordinates(fullAddress)
    if err != nil {
        http.Error(w, "Geocoding error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    shop.Lat = lat
    shop.Lon = lon

    result, err := db.DB.Exec(`
        UPDATE shops
        SET name=$1, country=$2, city=$3, address=$4, website=$5, description=$6, lat=$7, lon=$8
        WHERE id=$9`,
        shop.Name, shop.Country, shop.City, shop.Address, shop.Website, shop.Description, shop.Lat, shop.Lon, shopID)
    if err != nil {
        http.Error(w, "Database update error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    rowsAffected, err := result.RowsAffected()
    if err != nil || rowsAffected == 0 {
        http.Error(w, "Coffee shop not found", http.StatusNotFound)
        return
    }
    shop.ID = shopID
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(shop)
}


func DeleteCoffeeShopHandler(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    shopID, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid shop ID", http.StatusBadRequest)
        return
    }

    result, err := db.DB.Exec(`DELETE FROM shops WHERE id = $1`, shopID)
    if err != nil {
        http.Error(w, "Database delete error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    rowsAffected, err := result.RowsAffected()
    if err != nil || rowsAffected == 0 {
        http.Error(w, "Coffee shop not found", http.StatusNotFound)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}