package handlers

import (
	"encoding/json"
	"net/http"

	"coffeeApi/services/db"
)

type Stats struct {
    Users      int `json:"users"`
    Coffees    int `json:"coffees"`
    Roasteries int `json:"roasteries"`
    Shops      int `json:"shops"`
    Reviews    int `json:"reviews"`
}

func GetStatsHandler(w http.ResponseWriter, r *http.Request) {
    stats := Stats{}

    queries := map[string]*int{
        "SELECT COUNT(*) FROM users":      &stats.Users,
        "SELECT COUNT(*) FROM coffees":    &stats.Coffees,
        "SELECT COUNT(*) FROM roasteries": &stats.Roasteries,
        "SELECT COUNT(*) FROM shops":      &stats.Shops,
        "SELECT COUNT(*) FROM reviews":    &stats.Reviews,
    }

    for query, dest := range queries {
        if err := db.DB.QueryRow(query).Scan(dest); err != nil {
            http.Error(w, "Database error", http.StatusInternalServerError)
            return
        }
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(stats)
}