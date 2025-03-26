package handlers

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "strconv"
    "strings"

    "coffeeApi/services/db"
    "github.com/gorilla/mux"
)

type Coffee struct {
    ID           int      `json:"id"`
    Name         string   `json:"name"`
    RoasteryId   int      `json:"roasteryId"`
    Country      string   `json:"country"`
    Region       string   `json:"region"`
    Farm         string   `json:"farm"`
    Variety      string   `json:"variety"`
    Process      string   `json:"process"`
    RoastProfile string   `json:"roastProfile"`
    FlavourNotes []string `json:"flavourNotes"`
    Description  string   `json:"description"`
}

// GetCoffeesHandler retrieves all coffees from the database.
func GetCoffeesHandler(w http.ResponseWriter, r *http.Request) {
    rows, err := db.DB.Query(`
        SELECT id, name, roastery_id, country, region, farm, variety, process, roast_profile, flavour_notes, description 
        FROM coffees`)
    if err != nil {
        http.Error(w, "Database query error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var coffees []Coffee
    for rows.Next() {
        var c Coffee
        var notes string
        if err := rows.Scan(&c.ID, &c.Name, &c.RoasteryId, &c.Country, &c.Region, &c.Farm, &c.Variety, &c.Process, &c.RoastProfile, &notes, &c.Description); err != nil {
            http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
            return
        }
        // Assume flavourNotes are stored as comma-separated string.
        c.FlavourNotes = strings.Split(notes, ",")
        coffees = append(coffees, c)
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(coffees)
}

// GetCoffeeHandler retrieves a single coffee by id.
func GetCoffeeHandler(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    coffeeID, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid coffee ID", http.StatusBadRequest)
        return
    }

    var c Coffee
    var notes string
    err = db.DB.QueryRow(`SELECT id, name, roastery_id, country, region, farm, variety, process, roast_profile, flavour_notes, description 
        FROM coffees WHERE id = $1`, coffeeID).Scan(&c.ID, &c.Name, &c.RoasteryId, &c.Country, &c.Region, &c.Farm, &c.Variety, &c.Process, &c.RoastProfile, &notes, &c.Description)
    if err == sql.ErrNoRows {
        http.Error(w, "Coffee not found", http.StatusNotFound)
        return
    } else if err != nil {
        http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    c.FlavourNotes = strings.Split(notes, ",")
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(c)
}

// CreateCoffeeHandler inserts a new coffee record into the database.
func CreateCoffeeHandler(w http.ResponseWriter, r *http.Request) {
    var c Coffee
    if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
        http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
        return
    }
    if c.Name == "" || c.Country == "" || c.Process == "" || c.RoastProfile == "" {
        http.Error(w, "Missing required fields", http.StatusBadRequest)
        return
    }
    // Convert flavour notes to comma-separated string.
    notes := strings.Join(c.FlavourNotes, ",")

    err := db.DB.QueryRow(`
        INSERT INTO coffees 
        (name, roastery_id, country, region, farm, variety, process, roast_profile, flavour_notes, description) 
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id`,
        c.Name, c.RoasteryId, c.Country, c.Region, c.Farm, c.Variety, c.Process, c.RoastProfile, notes, c.Description).Scan(&c.ID)
    if err != nil {
        http.Error(w, "Database insert error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(c)
}

// UpdateCoffeeHandler updates an existing coffee record.
func UpdateCoffeeHandler(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    coffeeID, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid coffee ID", http.StatusBadRequest)
        return
    }

    var c Coffee
    if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
        http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
        return
    }
    notes := strings.Join(c.FlavourNotes, ",")

    result, err := db.DB.Exec(`
        UPDATE coffees 
        SET name=$1, roastery_id=$2, country=$3, region=$4, farm=$5, variety=$6, process=$7, roast_profile=$8, flavour_notes=$9, description=$10 
        WHERE id=$11`,
        c.Name, c.RoasteryId, c.Country, c.Region, c.Farm, c.Variety, c.Process, c.RoastProfile, notes, c.Description, coffeeID)
    if err != nil {
        http.Error(w, "Database update error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil || rowsAffected == 0 {
        http.Error(w, "Coffee not found", http.StatusNotFound)
        return
    }
    c.ID = coffeeID
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(c)
}

// DeleteCoffeeHandler deletes a coffee record from the database.
func DeleteCoffeeHandler(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    coffeeID, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid coffee ID", http.StatusBadRequest)
        return
    }

    result, err := db.DB.Exec(`DELETE FROM coffees WHERE id = $1`, coffeeID)
    if err != nil {
        http.Error(w, "Database delete error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    rowsAffected, err := result.RowsAffected()
    if err != nil || rowsAffected == 0 {
        http.Error(w, "Coffee not found", http.StatusNotFound)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}