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

type Roastery struct {
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

// GetRoasteriesHandler retrieves all roasteries from the database.
func GetRoasteriesHandler(w http.ResponseWriter, r *http.Request) {
    rows, err := db.DB.Query(`
        SELECT id, name, country, city, address, website, description, avg_rating, lat, lon
        FROM roasteries`)
    if err != nil {
        http.Error(w, "Database query error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var roasteries []Roastery
    for rows.Next() {
        var rastery Roastery
        if err := rows.Scan(&rastery.ID, &rastery.Name, &rastery.Country, &rastery.City, &rastery.Address,
            &rastery.Website, &rastery.Description, &rastery.AvgRating, &rastery.Lat, &rastery.Lon); err != nil {
            http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
            return
        }
        roasteries = append(roasteries, rastery)
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(roasteries)
}

// GetRoasteryHandler retrieves a single roastery by its ID.
func GetRoasteryHandler(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    roasteryID, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid roastery ID", http.StatusBadRequest)
        return
    }

    var rastery Roastery
    err = db.DB.QueryRow(`
        SELECT id, name, country, city, address, website, description, avg_rating, lat, lon
        FROM roasteries WHERE id = $1`, roasteryID).
        Scan(&rastery.ID, &rastery.Name, &rastery.Country, &rastery.City, &rastery.Address, &rastery.Website, &rastery.Description, &rastery.AvgRating, &rastery.Lat, &rastery.Lon)
    if err == sql.ErrNoRows {
        http.Error(w, "Roastery not found", http.StatusNotFound)
        return
    } else if err != nil {
        http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(rastery)
}

// CreateRoasteryHandler inserts a new roastery record into the database.
// It automatically obtains coordinates from the external geocoding API.
func CreateRoasteryHandler(w http.ResponseWriter, r *http.Request) {
    var rastery Roastery
    if err := json.NewDecoder(r.Body).Decode(&rastery); err != nil {
        http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
        return
    }
    if rastery.Name == "" || rastery.Country == "" || rastery.City == "" || rastery.Address == "" {
        http.Error(w, "Missing required fields", http.StatusBadRequest)
        return
    }

    // Construct full address from address, city, and country.
    fullAddress := fmt.Sprintf("%s, %s, %s", rastery.Address, rastery.City, rastery.Country)
    lat, lon, err := geocoding.GetCoordinates(fullAddress)
    if err != nil {
        http.Error(w, "Geocoding error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    rastery.Lat = lat
    rastery.Lon = lon

    // Set starting avg rating to 0.
    err = db.DB.QueryRow(`
        INSERT INTO roasteries (name, country, city, address, website, description, avg_rating, lat, lon)
        VALUES ($1, $2, $3, $4, $5, $6, 0, $7, $8) RETURNING id`,
        rastery.Name, rastery.Country, rastery.City, rastery.Address, rastery.Website, rastery.Description, rastery.Lat, rastery.Lon).
        Scan(&rastery.ID)
    if err != nil {
        http.Error(w, "Database insert error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    rastery.AvgRating = 0
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(rastery)
}

// UpdateRoasteryHandler updates an existing roastery record.
// It re-geocodes the address if any related fields are modified.
func UpdateRoasteryHandler(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    roasteryID, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid roastery ID", http.StatusBadRequest)
        return
    }

    var rastery Roastery
    if err := json.NewDecoder(r.Body).Decode(&rastery); err != nil {
        http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
        return
    }

    // Re-geocode the full address.
    fullAddress := fmt.Sprintf("%s, %s, %s", rastery.Address, rastery.City, rastery.Country)
    lat, lon, err := geocoding.GetCoordinates(fullAddress)
    if err != nil {
        http.Error(w, "Geocoding error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    rastery.Lat = lat
    rastery.Lon = lon

    result, err := db.DB.Exec(`
        UPDATE roasteries
        SET name=$1, country=$2, city=$3, address=$4, website=$5, description=$6, lat=$7, lon=$8
        WHERE id=$9`,
        rastery.Name, rastery.Country, rastery.City, rastery.Address, rastery.Website, rastery.Description, rastery.Lat, rastery.Lon, roasteryID)
    if err != nil {
        http.Error(w, "Database update error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    rowsAffected, err := result.RowsAffected()
    if err != nil || rowsAffected == 0 {
        http.Error(w, "Roastery not found", http.StatusNotFound)
        return
    }
    rastery.ID = roasteryID
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(rastery)
}

// DeleteRoasteryHandler deletes a roastery record from the database.
func DeleteRoasteryHandler(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    roasteryID, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid roastery ID", http.StatusBadRequest)
        return
    }

    result, err := db.DB.Exec(`DELETE FROM roasteries WHERE id = $1`, roasteryID)
    if err != nil {
        http.Error(w, "Database delete error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    rowsAffected, err := result.RowsAffected()
    if err != nil || rowsAffected == 0 {
        http.Error(w, "Roastery not found", http.StatusNotFound)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}