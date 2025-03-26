package handlers

import (
    "encoding/json"
    "net/http"
)

type ApiDocumentation struct {
    Name          string            `json:"name"`
    Description   string            `json:"description"`
    Version       string            `json:"version"`
    BaseURL       string            `json:"baseUrl"`
    Authorization AuthInfo          `json:"authorization"`
    Endpoints     map[string][]API  `json:"endpoints"`
}

type AuthInfo struct {
    Description string                  `json:"description"`
    Method      string                  `json:"method"`
    Header      string                  `json:"header"`
    Examples    map[string]interface{}  `json:"examples"`
}

type API struct {
    Method          string      `json:"method"`
    Path            string      `json:"path"`
    Description     string      `json:"description"`
    Auth            bool        `json:"requiresAuth"`
    AdminOnly       bool        `json:"adminOnly,omitempty"`
    PayloadExample  interface{} `json:"payloadExample,omitempty"`
    ResponseExample interface{} `json:"responseExample,omitempty"`
}

func GetApiDocumentationHandler(w http.ResponseWriter, r *http.Request) {
    doc := ApiDocumentation{
        Name:        "Coffee API",
        Description: "REST API for managing coffee data, including coffees, roasteries, coffee shops, and reviews",
        Version:     "1.0.0",
        BaseURL:     "http://localhost:40331",
        Authorization: AuthInfo{
            Description: "The API uses JWT (JSON Web Token) for authorization. Protected endpoints require a valid JWT token in the request header.",
            Method:      "Bearer Token Authentication",
            Header:      "Authentication: Bearer your_jwt_token",
            Examples: map[string]interface{}{
                "loginRequest": map[string]interface{}{
                    "url": "POST /login",
                    "body": map[string]string{
                        "username": "coffeeuser",
                        "password": "strongpassword",
                    },
                },
                "loginResponse": map[string]interface{}{
                    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
                },
                "authorizedRequest": map[string]interface{}{
                    "url": "POST /coffees",
                    "headers": map[string]string{
                        "Content-Type": "application/json",
                        "Authentication": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
                    },
                },
                "roles": map[string]string{
                    "user": "Regular user with access to create/update content",
                    "admin": "Administrator with additional privileges like deletion",
                },
            },
        },
        Endpoints: map[string][]API{
            "Authentication": {
                {
                    Method:      "POST",
                    Path:        "/register",
                    Description: "Register a new user",
                    Auth:        false,
                    PayloadExample: map[string]interface{}{
                        "username": "coffeeuser",
                        "password": "strongpassword",
                        "email":    "user@example.com",
                    },
                },
                {
                    Method:      "POST",
                    Path:        "/login",
                    Description: "Login and get JWT token",
                    Auth:        false,
                    PayloadExample: map[string]interface{}{
                        "username": "coffeeuser",
                        "password": "strongpassword",
                    },
                    ResponseExample: map[string]string{
                        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
                    },
                },
            },
            "Coffees": {
                {
                    Method:      "GET",
                    Path:        "/coffees",
                    Description: "Get all coffees with optional filtering",
                    Auth:        false,
                },
                {
                    Method:      "GET",
                    Path:        "/coffees/{id}",
                    Description: "Get coffee by ID",
                    Auth:        false,
                },
                {
                    Method:      "POST",
                    Path:        "/coffees",
                    Description: "Add a new coffee",
                    Auth:        true,
                    PayloadExample: map[string]interface{}{
                        "name":         "Ethiopia Yirgacheffe",
                        "roasteryId":   1,
                        "country":      "Ethiopia",
                        "region":       "Yirgacheffe",
                        "process":      "Washed",
                        "roastProfile": "Light",
                        "flavourNotes": []string{"floral", "citrus", "bergamot"},
                        "description":  "A bright and floral Ethiopian coffee",
                        "imageUrl":     "https://example.com/coffee.jpg",
                    },
                },
                {
                    Method:      "PUT",
                    Path:        "/coffees/{id}",
                    Description: "Update a coffee",
                    Auth:        true,
                },
                {
                    Method:      "DELETE",
                    Path:        "/coffees/{id}",
                    Description: "Delete a coffee",
                    Auth:        true,
                },
            },
            "Roasteries": {
                {
                    Method:      "GET",
                    Path:        "/roasteries",
                    Description: "Get all roasteries",
                    Auth:        false,
                },
                {
                    Method:      "GET",
                    Path:        "/roasteries/{id}",
                    Description: "Get roastery by ID",
                    Auth:        false,
                },
                {
                    Method:      "POST",
                    Path:        "/roasteries",
                    Description: "Add a new roastery",
                    Auth:        true,
                    PayloadExample: map[string]interface{}{
                        "name":        "Coffee Roasters Inc",
                        "country":     "USA",
                        "city":        "Portland",
                        "address":     "123 Roast St",
                        "website":     "https://example.com",
                        "description": "Specialty coffee roastery",
                        "imageUrl":    "https://example.com/roastery.jpg",
                    },
                },
                {
                    Method:      "PUT",
                    Path:        "/roasteries/{id}",
                    Description: "Update a roastery",
                    Auth:        true,
                },
                {
                    Method:      "DELETE",
                    Path:        "/roasteries/{id}",
                    Description: "Delete a roastery",
                    Auth:        true,
                    AdminOnly:   true,
                },
            },
            "Coffee Shops": {
                {
                    Method:      "GET",
                    Path:        "/shops",
                    Description: "Get all coffee shops",
                    Auth:        false,
                },
                {
                    Method:      "GET",
                    Path:        "/shops/{id}",
                    Description: "Get coffee shop by ID",
                    Auth:        false,
                },
                {
                    Method:      "POST",
                    Path:        "/shops",
                    Description: "Add a new coffee shop",
                    Auth:        true,
                    PayloadExample: map[string]interface{}{
                        "name":        "Coffee Corner",
                        "country":     "Germany",
                        "city":        "Berlin",
                        "address":     "Bergmannstraße 100",
                        "website":     "https://example.com",
                        "description": "Cozy specialty coffee shop",
                        "imageUrl":    "https://example.com/shop.jpg",
                    },
                },
                {
                    Method:      "PUT",
                    Path:        "/shops/{id}",
                    Description: "Update a coffee shop",
                    Auth:        true,
                },
                {
                    Method:      "DELETE",
                    Path:        "/shops/{id}",
                    Description: "Delete a coffee shop",
                    Auth:        true,
                    AdminOnly:   true,
                },
            },
            "Reviews": {
                {
                    Method:      "GET",
                    Path:        "/reviews",
                    Description: "Get all reviews with optional filtering",
                    Auth:        false,
                },
                {
                    Method:      "POST",
                    Path:        "/reviews",
                    Description: "Add a new review",
                    Auth:        true,
                    PayloadExample: map[string]interface{}{
                        "coffeeId": 1,
                        "rating":   4.5,
                        "review":   "Great coffee with bright acidity and floral notes",
                    },
                },
                {
                    Method:      "PUT",
                    Path:        "/reviews/{id}",
                    Description: "Update a review",
                    Auth:        true,
                },
                {
                    Method:      "DELETE",
                    Path:        "/reviews/{id}",
                    Description: "Delete a review",
                    Auth:        true,
                },
            },
        },
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(doc)
}