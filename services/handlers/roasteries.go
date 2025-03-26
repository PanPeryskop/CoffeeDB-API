package handlers

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"
    "strings"

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

func GetRoasteriesHandler(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query()
    name := query.Get("name")
    country := query.Get("country")
    city := query.Get("city")
    address := query.Get("address")
    website := query.Get("website")
    description := query.Get("description")
    minRating := query.Get("minRating")
    maxRating := query.Get("maxRating")

    baseQuery := `SELECT id, name, country, city, address, website, description, avg_rating, lat, lon FROM roasteries`
    conditions := []string{}
    args := []interface{}{}
    argIdx := 1

    if name != "" {
        conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", argIdx))
        args = append(args, "%"+name+"%")
        argIdx++
    }
    if country != "" {
        conditions = append(conditions, fmt.Sprintf("country ILIKE $%d", argIdx))
        args = append(args, "%"+country+"%")
        argIdx++
    }
    if city != "" {
        conditions = append(conditions, fmt.Sprintf("city ILIKE $%d", argIdx))
        args = append(args, "%"+city+"%")
        argIdx++
    }
    if address != "" {
        conditions = append(conditions, fmt.Sprintf("address ILIKE $%d", argIdx))
        args = append(args, "%"+address+"%")
        argIdx++
    }
    if website != "" {
        conditions = append(conditions, fmt.Sprintf("website ILIKE $%d", argIdx))
        args = append(args, "%"+website+"%")
        argIdx++
    }
    if description != "" {
        conditions = append(conditions, fmt.Sprintf("description ILIKE $%d", argIdx))
        args = append(args, "%"+description+"%")
        argIdx++
    }
    if minRating != "" {
        if rating, err := strconv.ParseFloat(minRating, 32); err == nil {
            conditions = append(conditions, fmt.Sprintf("avg_rating >= $%d", argIdx))
            args = append(args, rating)
            argIdx++
        }
    }
    if maxRating != "" {
        if rating, err := strconv.ParseFloat(maxRating, 32); err == nil {
            conditions = append(conditions, fmt.Sprintf("avg_rating <= $%d", argIdx))
            args = append(args, rating)
            argIdx++
        }
    }

    if len(conditions) > 0 {
        baseQuery += " WHERE " + strings.Join(conditions, " AND ")
    }

    rows, err := db.DB.Query(baseQuery, args...)
    if err != nil {
        http.Error(w, "Database query error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var roasteries []Roastery
    for rows.Next() {
        var rastery Roastery
        err := rows.Scan(&rastery.ID, &rastery.Name, &rastery.Country, &rastery.City, &rastery.Address, &rastery.Website, &rastery.Description, &rastery.AvgRating, &rastery.Lat, &rastery.Lon)
        if err != nil {
            http.Error(w, "Error scanning roastery: "+err.Error(), http.StatusInternalServerError)
            return
        }
        roasteries = append(roasteries, rastery)
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(roasteries)
}

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

    fullAddress := fmt.Sprintf("%s, %s, %s", rastery.Address, rastery.City, rastery.Country)
    lat, lon, err := geocoding.GetCoordinates(fullAddress)
    if err != nil {
        http.Error(w, "Geocoding error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    rastery.Lat = lat
    rastery.Lon = lon

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

    fullAddress := fmt.Sprintf("%s, %s, %s", rastery.Address, rastery.City, rastery.Country)
    lat, lon, err := geocoding.GetCoordinates(fullAddress)
    if err != nil {
        http.Error(w, "Geocoding error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    rastery.Lat = lat
    rastery.Lon = lon

    result, err := db.DB.Exec(`
        UPDATE roasteries SET name=$1, country=$2, city=$3, address=$4, website=$5, description=$6, lat=$7, lon=$8
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

func DeleteRoasteryHandler(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    roasteryID, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid roastery ID", http.StatusBadRequest)
        return
    }

    var coffeeCount int
    err = db.DB.QueryRow("SELECT COUNT(*) FROM coffees WHERE roastery_id = $1", roasteryID).Scan(&coffeeCount)
    if err != nil {
        http.Error(w, "Database query error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    if coffeeCount > 0 {
        http.Error(w, "Cannot delete roastery that has associated coffees", http.StatusBadRequest)
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