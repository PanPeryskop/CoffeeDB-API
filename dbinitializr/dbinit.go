package main

import (
    "encoding/json"
    "fmt"
    "io"
    "log"
    "os"
    "strings"

    "coffeeApi/services/db"
    "coffeeApi/services/geocoding"
    _ "github.com/lib/pq"
    "golang.org/x/crypto/bcrypt"
)

func createTables() error {
    queries := []string{
        `CREATE TABLE IF NOT EXISTS users(
            id SERIAL PRIMARY KEY,
            username TEXT NOT NULL UNIQUE,
            password TEXT NOT NULL,
            email TEXT NOT NULL,
            role TEXT NOT NULL,
            avatar_url TEXT
        )`,
        `CREATE TABLE IF NOT EXISTS coffees(
            id SERIAL PRIMARY KEY,
            name TEXT NOT NULL,
            roastery_id INTEGER,
            country TEXT,
            region TEXT,
            farm TEXT,
            variety TEXT,
            process TEXT,
            roast_profile TEXT,
            flavour_notes TEXT,
            description TEXT,
            image_url TEXT
        )`,
        `CREATE TABLE IF NOT EXISTS roasteries(
            id SERIAL PRIMARY KEY,
            name TEXT NOT NULL,
            country TEXT,
            city TEXT,
            address TEXT,
            website TEXT,
            description TEXT,
            avg_rating REAL,
            lat REAL,
            lon REAL,
            image_url TEXT
        )`,
        `CREATE TABLE IF NOT EXISTS shops(
            id SERIAL PRIMARY KEY,
            name TEXT NOT NULL,
            country TEXT,
            city TEXT,
            address TEXT,
            website TEXT,
            description TEXT,
            avg_rating REAL,
            lat REAL,
            lon REAL,
            image_url TEXT
        )`,
        `CREATE TABLE IF NOT EXISTS reviews(
            id SERIAL PRIMARY KEY,
            user_id INTEGER,
            coffee_id INTEGER,
            roastery_id INTEGER,
            coffee_shop_id INTEGER,
            rating REAL,
            review TEXT,
            date_of_creation TIMESTAMPTZ
        )`,
    }

    for _, q := range queries {
        if _, err := db.DB.Exec(q); err != nil {
            return fmt.Errorf("error executing query: %v, error: %v", q, err)
        }
    }
    return nil
}

type Data struct {
    Users      []User       `json:"users"`
    Coffees    []Coffee     `json:"coffees"`
    Roasteries []Roastery   `json:"roasteries"`
    Shops      []CoffeeShop `json:"shops"`
    Reviews    []Review     `json:"reviews"`
}

type User struct {
    ID        int    `json:"id"`
    Username  string `json:"username"`
    Password  string `json:"password"`
    Email     string `json:"email"`
    Role      string `json:"role"`
    AvatarURL string `json:"avatarUrl"`
}

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
    ImageURL     string   `json:"imageUrl"`
}

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
    ImageURL    string  `json:"imageUrl"`
}

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
    ImageURL    string  `json:"imageUrl"`
}

type Review struct {
    ID             int     `json:"id"`
    UserId         int     `json:"userId"`
    CoffeeId       int     `json:"coffeeId"`
    RoasteryId     int     `json:"roasteryId"`
    CoffeeShopId   int     `json:"coffeeShopId"`
    Rating         float32 `json:"rating"`
    Review         string  `json:"review"`
    DateOfCreation string  `json:"dateOfCreation"`
}

func tableIsEmpty(query string) (bool, error) {
    var count int
    err := db.DB.QueryRow(query).Scan(&count)
    if err != nil {
        return false, err
    }
    return count == 0, nil
}

func seedData(filePath string) error {
    dataFile, err := os.Open(filePath)
    if err != nil {
        return fmt.Errorf("error opening data file: %v", err)
    }
    defer dataFile.Close()
    bytes, err := io.ReadAll(dataFile)
    if err != nil {
        return fmt.Errorf("error reading data file: %v", err)
    }

    var data Data
    if err := json.Unmarshal(bytes, &data); err != nil {
        return fmt.Errorf("error unmarshalling JSON: %v", err)
    }

    tx, err := db.DB.Begin()
    if err != nil {
        return fmt.Errorf("error beginning transaction: %v", err)
    }

    empty, err := tableIsEmpty("SELECT COUNT(*) FROM users")
    if err != nil {
        tx.Rollback()
        return err
    }
    if empty {
        for _, u := range data.Users {
            // Hashowanie hasła
            hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
            if err != nil {
                tx.Rollback()
                return fmt.Errorf("error hashing password for user %v: %v", u.Username, err)
            }
            _, err = tx.Exec(`INSERT INTO users (username, password, email, role, avatar_url) 
                              VALUES ($1, $2, $3, $4, $5)`,
                u.Username, string(hashedPassword), u.Email, u.Role, u.AvatarURL)
            if err != nil {
                tx.Rollback()
                return fmt.Errorf("error inserting user %v: %v", u.Username, err)
            }
        }
    }

    empty, err = tableIsEmpty("SELECT COUNT(*) FROM coffees")
    if err != nil {
        tx.Rollback()
        return err
    }
    if empty {
        for _, c := range data.Coffees {
            notes := ""
            if len(c.FlavourNotes) > 0 {
                notes = strings.Join(c.FlavourNotes, ",")
            }
            _, err := tx.Exec(`INSERT INTO coffees (name, roastery_id, country, region, farm, variety, process, 
                              roast_profile, flavour_notes, description, image_url)
                              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
                c.Name, c.RoasteryId, c.Country, c.Region, c.Farm, c.Variety, c.Process, 
                c.RoastProfile, notes, c.Description, c.ImageURL)
            if err != nil {
                tx.Rollback()
                return fmt.Errorf("error inserting coffee %v: %v", c.Name, err)
            }
        }
    }

    empty, err = tableIsEmpty("SELECT COUNT(*) FROM roasteries")
    if err != nil {
        tx.Rollback()
        return err
    }
    if empty {
        for _, r := range data.Roasteries {
            if r.Lat == 0 && r.Lon == 0 {
                fullAddress := fmt.Sprintf("%s, %s, %s", r.Address, r.City, r.Country)
                lat, lon, err := geocoding.GetCoordinates(fullAddress)
                if err != nil {
                    fmt.Printf("Warning: could not geocode roastery %s: %v\n", r.Name, err)
                } else {
                    r.Lat = lat
                    r.Lon = lon
                }
            }
            _, err := tx.Exec(`INSERT INTO roasteries (name, country, city, address, website, description, 
                              avg_rating, lat, lon, image_url)
                              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
                r.Name, r.Country, r.City, r.Address, r.Website, r.Description, 
                r.AvgRating, r.Lat, r.Lon, r.ImageURL)
            if err != nil {
                tx.Rollback()
                return fmt.Errorf("error inserting roastery %v: %v", r.Name, err)
            }
        }
    }

    empty, err = tableIsEmpty("SELECT COUNT(*) FROM shops")
    if err != nil {
        tx.Rollback()
        return err
    }
    if empty {
        for _, s := range data.Shops {
            if s.Lat == 0 && s.Lon == 0 {
                fullAddress := fmt.Sprintf("%s, %s, %s", s.Address, s.City, s.Country)
                lat, lon, err := geocoding.GetCoordinates(fullAddress)
                if err != nil {
                    fmt.Printf("Warning: could not geocode coffee shop %s: %v\n", s.Name, err)
                } else {
                    s.Lat = lat
                    s.Lon = lon
                }
            }
            _, err := tx.Exec(`INSERT INTO shops (name, country, city, address, website, description, 
                              avg_rating, lat, lon, image_url)
                              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
                s.Name, s.Country, s.City, s.Address, s.Website, s.Description, 
                s.AvgRating, s.Lat, s.Lon, s.ImageURL)
            if err != nil {
                tx.Rollback()
                return fmt.Errorf("error inserting shop %v: %v", s.Name, err)
            }
        }
    }

    empty, err = tableIsEmpty("SELECT COUNT(*) FROM reviews")
    if err != nil {
        tx.Rollback()
        return err
    }
    if empty {
        for _, rev := range data.Reviews {
            _, err := tx.Exec(`INSERT INTO reviews (user_id, coffee_id, roastery_id, coffee_shop_id, 
                              rating, review, date_of_creation)
                              VALUES ($1, $2, $3, $4, $5, $6, $7)`,
                rev.UserId, rev.CoffeeId, rev.RoasteryId, rev.CoffeeShopId, 
                rev.Rating, rev.Review, rev.DateOfCreation)
            if err != nil {
                tx.Rollback()
                return fmt.Errorf("error inserting review: %v", err)
            }
        }
    }

    if err := tx.Commit(); err != nil {
        return fmt.Errorf("error committing transaction: %v", err)
    }
    fmt.Println("Data seeded successfully!")
    return nil
}

func main() {
    if err := db.Init(); err != nil {
        log.Fatal("Database initialization error:", err)
    }

    if err := createTables(); err != nil {
        log.Fatal("Error creating tables:", err)
    }

    if err := seedData("dbinitializr/data.json"); err != nil {
        log.Fatal(err)
    }
}