package handlers

import (
    "encoding/json"
    "net/http"
    "strings"
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
        Description: "REST API for managing coffee data, including coffees, roasteries, coffee shops, and reviews. For interactive documentation, visit the /help endpoint.",
        Version:     "1.0.0",
        BaseURL:     "http://srv17.mikr.us:40331",
        Authorization: AuthInfo{
            Description: "The API uses JWT (JSON Web Token) for authorization. Protected endpoints require a valid JWT token in the request header.",
            Method:      "Bearer Token Authentication",
            Header:      "Authentication: Bearer your_jwt_token",
            Examples: map[string]interface{}{
                "loginRequest": map[string]interface{}{
                    "url": "POST /login",
                    "body": map[string]string{
                        "username": "coffeegeek",
                        "password": "kawa123",
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

func GetHtmlDocumentationHandler(w http.ResponseWriter, r *http.Request) {
    doc := ApiDocumentation{
        Name:        "Coffee API",
        Description: "REST API for managing coffee data, including coffees, roasteries, coffee shops, and reviews",
        Version:     "1.0.0",
        BaseURL:     "http://srv17.mikr.us:40331",
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

    html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>` + doc.Name + ` Documentation</title>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css">
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.7.0/styles/atom-one-dark.min.css">
    <script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.7.0/highlight.min.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.7.0/languages/json.min.js"></script>
    <style>
        :root {
            --primary-color: #6F4E37;
            --secondary-color: #B87333;
            --background-color: #FFFAF0;
            --text-color: #333;
            --header-color: #5D4037;
            --get-color: #61affe;
            --post-color: #49cc90;
            --put-color: #fca130;
            --delete-color: #f93e3e;
            --accent-color: #795548;
        }
        
        * {
            box-sizing: border-box;
            margin: 0;
            padding: 0;
        }
        
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            line-height: 1.6;
            color: var(--text-color);
            background-color: var(--background-color);
            overflow-x: hidden;
        }
        
        .container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 0 20px;
        }
        
        header {
            background: linear-gradient(135deg, var(--primary-color), var(--secondary-color));
            color: white;
            padding: 40px 0;
            text-align: center;
            position: relative;
            overflow: hidden;
            box-shadow: 0 4px 12px rgba(0,0,0,0.1);
        }
        
        header::before {
            content: '';
            position: absolute;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            background: url("data:image/svg+xml,%3Csvg width='100' height='100' viewBox='0 0 100 100' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath d='M11 18c3.866 0 7-3.134 7-7s-3.134-7-7-7-7 3.134-7 7 3.134 7 7 7zm48 25c3.866 0 7-3.134 7-7s-3.134-7-7-7-7 3.134-7 7 3.134 7 7 7zm-43-7c1.657 0 3-1.343 3-3s-1.343-3-3-3-3 1.343-3 3 1.343 3 3 3zm63 31c1.657 0 3-1.343 3-3s-1.343-3-3-3-3 1.343-3 3 1.343 3 3 3zM34 90c1.657 0 3-1.343 3-3s-1.343-3-3-3-3 1.343-3 3 1.343 3 3 3zm56-76c1.657 0 3-1.343 3-3s-1.343-3-3-3-3 1.343-3 3 1.343 3 3 3zM12 86c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm28-65c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm23-11c2.76 0 5-2.24 5-5s-2.24-5-5-5-5 2.24-5 5 2.24 5 5 5zm-6 60c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm29 22c2.76 0 5-2.24 5-5s-2.24-5-5-5-5 2.24-5 5 2.24 5 5 5zM32 63c2.76 0 5-2.24 5-5s-2.24-5-5-5-5 2.24-5 5 2.24 5 5 5zm57-13c2.76 0 5-2.24 5-5s-2.24-5-5-5-5 2.24-5 5 2.24 5 5 5zm-9-21c1.105 0 2-.895 2-2s-.895-2-2-2-2 .895-2 2 .895 2 2 2zM60 91c1.105 0 2-.895 2-2s-.895-2-2-2-2 .895-2 2 .895 2 2 2zM35 41c1.105 0 2-.895 2-2s-.895-2-2-2-2 .895-2 2 .895 2 2 2zM12 60c1.105 0 2-.895 2-2s-.895-2-2-2-2 .895-2 2 .895 2 2 2z' fill='rgba(255,255,255,.075)' fill-rule='evenodd'/%3E%3C/svg%3E");
            opacity: 0.3;
        }
        
        .coffee-cup {
            margin-bottom: 15px;
            font-size: 4rem;
            animation: steam 3s infinite ease-in-out;
        }
        
        @keyframes steam {
            0%, 100% { 
                transform: translateY(0) rotate(0deg); 
                opacity: 0.8;
            }
            50% { 
                transform: translateY(-10px) rotate(5deg); 
                opacity: 1;
            }
        }
        
        h1 {
            font-size: 2.8rem;
            margin-bottom: 15px;
            letter-spacing: 1px;
            text-shadow: 2px 2px 4px rgba(0,0,0,0.2);
        }
        
        .version-badge {
            display: inline-block;
            background-color: rgba(255,255,255,0.2);
            padding: 5px 12px;
            border-radius: 20px;
            font-size: 0.9rem;
            margin-bottom: 15px;
            backdrop-filter: blur(5px);
        }
        
        header p {
            max-width: 700px;
            margin: 0 auto;
            font-size: 1.2rem;
            opacity: 0.9;
        }
        
        .main-content {
            padding: 40px 0;
        }
        
        .section {
            margin-bottom: 50px;
            animation: fadeIn 0.5s ease-in-out;
        }
        
        @keyframes fadeIn {
            from { opacity: 0; transform: translateY(20px); }
            to { opacity: 1; transform: translateY(0); }
        }
        
        h2 {
            color: var(--header-color);
            font-size: 2rem;
            margin-bottom: 20px;
            padding-bottom: 10px;
            border-bottom: 3px solid var(--primary-color);
            display: inline-block;
        }
        
        h3 {
            color: var(--header-color);
            font-size: 1.5rem;
            margin: 30px 0 15px;
            display: flex;
            align-items: center;
        }
        
        h3 i {
            margin-right: 10px;
            color: var(--primary-color);
        }
        
        .base-url {
            background-color: #2d2d2d;
            color: white;
            padding: 15px;
            border-radius: 5px;
            margin-bottom: 20px;
            font-family: monospace;
            position: relative;
            box-shadow: 0 3px 6px rgba(0,0,0,0.1);
        }
        
        .copy-btn {
            position: absolute;
            top: 8px;
            right: 8px;
            background-color: rgba(255,255,255,0.15);
            border: none;
            color: white;
            padding: 5px 10px;
            border-radius: 3px;
            cursor: pointer;
            font-size: 0.8rem;
            transition: all 0.2s;
        }
        
        .copy-btn:hover {
            background-color: rgba(255,255,255,0.25);
        }
        
        .auth-info {
            background-color: white;
            border-radius: 8px;
            padding: 20px;
            box-shadow: 0 4px 10px rgba(0,0,0,0.05);
            margin-bottom: 20px;
        }
        
        .auth-info p, .auth-info div {
            margin-bottom: 15px;
        }
			
        .auth-info code {
            padding: 3px 6px;
            border-radius: 3px;
            font-family: monospace;
            color: #e0e0e0;
        }
        
        .categories-nav {
            display: flex;
            flex-wrap: wrap;
            gap: 10px;
            margin-bottom: 30px;
        }
        
        .category-btn {
            background-color: var(--primary-color);
            color: white;
            border: none;
            padding: 8px 15px;
            border-radius: 20px;
            cursor: pointer;
            transition: all 0.2s;
            font-weight: 500;
        }
        
        .category-btn:hover, .category-btn.active {
            background-color: var(--secondary-color);
            transform: translateY(-2px);
            box-shadow: 0 3px 8px rgba(0,0,0,0.15);
        }
        
        .category {
            margin-bottom: 40px;
            padding: 25px;
            background-color: white;
            border-radius: 10px;
            box-shadow: 0 4px 15px rgba(0,0,0,0.05);
            transition: all 0.3s;
        }
        
        .category:hover {
            transform: translateY(-5px);
            box-shadow: 0 8px 25px rgba(0,0,0,0.1);
        }
        
        .category h3 {
            margin-top: 0;
            color: var(--header-color);
            border-bottom: 2px solid var(--accent-color);
            padding-bottom: 10px;
            margin-bottom: 20px;
        }
        
        .endpoints {
            display: grid;
            gap: 20px;
        }
        
        .endpoint {
            background-color: #f9f9f9;
            border-radius: 8px;
            overflow: hidden;
            box-shadow: 0 2px 5px rgba(0,0,0,0.05);
            transition: all 0.3s;
        }
        
        .endpoint:hover {
            box-shadow: 0 5px 15px rgba(0,0,0,0.1);
        }
        
        .endpoint-header {
            padding: 15px;
            display: flex;
            align-items: center;
            gap: 15px;
            background-color: #f1f1f1;
            cursor: pointer;
        }
        
        .method {
            padding: 5px 10px;
            border-radius: 5px;
            font-weight: bold;
            color: white;
            min-width: 70px;
            text-align: center;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        
        .get { background-color: var(--get-color); }
        .post { background-color: var(--post-color); }
        .put { background-color: var(--put-color); }
        .delete { background-color: var(--delete-color); }
        
        .path {
            font-family: monospace;
            font-size: 1.1em;
            font-weight: 500;
            flex-grow: 1;
        }
        
        .auth-badge {
            background-color: #e74c3c;
            color: white;
            padding: 3px 8px;
            border-radius: 3px;
            font-size: 0.8em;
            display: inline-flex;
            align-items: center;
            gap: 5px;
        }
        
        .admin-badge {
            background-color: #9b59b6;
            color: white;
            padding: 3px 8px;
            border-radius: 3px;
            font-size: 0.8em;
            display: inline-flex;
            align-items: center;
            gap: 5px;
        }
        
        .endpoint-body {
            padding: 0;
            max-height: 0;
            overflow: hidden;
            transition: all 0.3s ease;
        }
        
        .endpoint-body.active {
            padding: 20px;
            max-height: 2000px;
        }
        
        .description {
            margin-bottom: 20px;
            padding-bottom: 15px;
            border-bottom: 1px solid #eee;
        }
        
        .code-example h4 {
            margin: 20px 0 10px;
            color: var(--header-color);
        }
        
        pre {
            background-color: #282c34;
            border-radius: 5px;
            padding: 15px;
            overflow-x: auto;
            margin: 10px 0;
            position: relative;
        }
        
        .try-it {
            margin-top: 20px;
            padding-top: 20px;
            border-top: 1px solid #eee;
        }
        
        .try-it h4 {
            margin-bottom: 15px;
            display: flex;
            align-items: center;
            gap: 10px;
        }
        
        .try-it h4 i {
            color: var(--secondary-color);
        }
        
        .input-group {
            margin-bottom: 15px;
        }
        
        .input-group label {
            display: block;
            margin-bottom: 5px;
            font-weight: 500;
        }
        
        .input-group input, .input-group textarea {
            width: 100%;
            padding: 10px;
            border: 1px solid #ddd;
            border-radius: 5px;
            font-family: monospace;
        }
        
        .auth-toggle {
            margin-bottom: 15px;
        }
        
        .auth-input {
            display: none;
            margin-top: 10px;
        }
        
        .execute-btn {
            background-color: var(--secondary-color);
            color: white;
            border: none;
            padding: 10px 20px;
            border-radius: 5px;
            cursor: pointer;
            font-weight: 500;
            display: flex;
            align-items: center;
            gap: 8px;
            transition: all 0.2s;
        }
        
        .execute-btn:hover {
            background-color: #a25f28;
            transform: translateY(-2px);
        }
        
        .response-container {
            margin-top: 20px;
            display: none;
        }
        
        .response-container.show {
            display: block;
            animation: slideDown 0.3s ease-out;
        }
        
        @keyframes slideDown {
            from { opacity: 0; transform: translateY(-20px); }
            to { opacity: 1; transform: translateY(0); }
        }
        
        .response-info {
            display: flex;
            align-items: center;
            gap: 10px;
            margin-bottom: 10px;
        }
        
        .status {
            padding: 3px 8px;
            border-radius: 3px;
            font-weight: 500;
        }
        
        .status.success { background-color: #2ecc71; color: white; }
        .status.error { background-color: #e74c3c; color: white; }
        
        footer {
            background-color: var(--primary-color);
            color: white;
            text-align: center;
            padding: 30px 0;
            margin-top: 50px;
        }
        
        .scroll-top {
            position: fixed;
            bottom: 30px;
            right: 30px;
            background-color: var(--primary-color);
            color: white;
            width: 50px;
            height: 50px;
            border-radius: 50%;
            display: flex;
            justify-content: center;
            align-items: center;
            cursor: pointer;
            box-shadow: 0 4px 10px rgba(0,0,0,0.2);
            opacity: 0;
            visibility: hidden;
            transition: all 0.3s;
        }
        
        .scroll-top.show {
            opacity: 1;
            visibility: visible;
        }
        
        .scroll-top:hover {
            background-color: var(--secondary-color);
            transform: translateY(-5px);
        }
        
        @media (max-width: 768px) {
            h1 {
                font-size: 2.2rem;
            }
            
            .endpoint-header {
                flex-wrap: wrap;
            }
            
            .auth-badge, .admin-badge {
                margin-top: 10px;
            }
        }
        
        .theme-toggle {
            position: fixed;
            top: 20px;
            right: 20px;
            background-color: rgba(0,0,0,0.2);
            color: white;
            width: 40px;
            height: 40px;
            border-radius: 50%;
            display: flex;
            justify-content: center;
            align-items: center;
            cursor: pointer;
            z-index: 100;
            transition: all 0.3s;
        }
        
        .theme-toggle:hover {
            background-color: rgba(0,0,0,0.3);
            transform: rotate(30deg);
        }
        
        body.dark-theme {
            --background-color: #1a1a1a;
            --text-color: #f1f1f1;
            --header-color: #e0e0e0;
        }
        
        body.dark-theme .category,
        body.dark-theme .auth-info {
            background-color: #2d2d2d;
        }
        
        body.dark-theme .endpoint {
            background-color: #2d2d2d;
        }
        
        body.dark-theme .endpoint-header {
            background-color: #232323;
        }
        
        body.dark-theme .path,
        body.dark-theme .description,
        body.dark-theme h4,
        body.dark-theme label {
            color: #e0e0e0;
        }
        
        body.dark-theme .auth-info code {
            background-color: #3a3a3a;
            color: #e0e0e0;
        }
        
        body.dark-theme input, 
        body.dark-theme textarea {
            background-color: #3a3a3a;
            color: #e0e0e0;
            border-color: #555;
        }
        
        .loading {
            display: inline-block;
            width: 20px;
            height: 20px;
            border: 3px solid rgba(255,255,255,.3);
            border-radius: 50%;
            border-top-color: white;
            animation: spin 1s ease-in-out infinite;
        }
        
        @keyframes spin {
            to { transform: rotate(360deg); }
        }
    </style>
</head>
<body>
    <div class="theme-toggle" onclick="toggleTheme()">
        <i class="fas fa-moon"></i>
    </div>
    
    <header>
        <div class="container">
            <div class="coffee-cup">
                <i class="fas fa-mug-hot"></i>
            </div>
            <h1>` + doc.Name + `</h1>
            <div class="version-badge">v` + doc.Version + `</div>
            <p>` + doc.Description + `</p>
        </div>
    </header>
    
    <main class="container main-content">
        <section class="section" id="overview">
            <h2><i class="fas fa-info-circle"></i> Overview</h2>
            <p>This API allows you to interact with our coffee database, providing endpoints to manage coffee data, roasteries, coffee shops, and reviews.</p>
            
            <div class="base-url" id="base-url">
                ` + doc.BaseURL + `
                <button class="copy-btn" onclick="copyToClipboard('base-url')">Copy</button>
            </div>
        </section>
        
        <section class="section" id="authentication">
            <h2><i class="fas fa-lock"></i> Authentication</h2>
            <div class="auth-info">
                <p>` + doc.Authorization.Description + `</p>
                <div><strong>Method:</strong> ` + doc.Authorization.Method + `</div>
                <div><strong>Header:</strong> <code>` + doc.Authorization.Header + `</code></div>
                
                <h4>Authentication Examples</h4>
                <pre><code class="language-json">` + formatAuthExamples(doc.Authorization.Examples) + `</code></pre>
            </div>
        </section>
        
        <section class="section" id="endpoints">
            <h2><i class="fas fa-plug"></i> API Endpoints</h2>
            
            <div class="categories-nav">
                <button class="category-btn active" onclick="filterCategory('all')">All</button>`

    for category := range doc.Endpoints {
        html += `<button class="category-btn" onclick="filterCategory('` + category + `')">` + category + `</button>`
    }
    
    html += `</div>`

    // Icons for categories
    categoryIcons := map[string]string{
        "Authentication": "fas fa-lock",
        "Coffees": "fas fa-coffee",
        "Roasteries": "fas fa-industry",
        "Coffee Shops": "fas fa-store",
        "Reviews": "fas fa-star",
    }

    // For each category
    for category, endpoints := range doc.Endpoints {
        icon := "fas fa-folder"
        if val, ok := categoryIcons[category]; ok {
            icon = val
        }

        html += `<div class="category" data-category="` + category + `">
            <h3><i class="` + icon + `"></i> ` + category + `</h3>
            <div class="endpoints">`

        for _, endpoint := range endpoints {
            methodClass := strings.ToLower(endpoint.Method)
            
            html += `<div class="endpoint">
                <div class="endpoint-header" onclick="toggleEndpoint(this)">
                    <div class="method ` + methodClass + `">` + endpoint.Method + `</div>
                    <div class="path">` + endpoint.Path + `</div>`
            
            if endpoint.Auth {
                html += `<div class="auth-badge"><i class="fas fa-lock"></i> Auth Required</div>`
            }
            
            if endpoint.AdminOnly {
                html += `<div class="admin-badge"><i class="fas fa-crown"></i> Admin Only</div>`
            }
            
            html += `</div>
                <div class="endpoint-body">
                    <div class="description">` + endpoint.Description + `</div>`
            
            if endpoint.PayloadExample != nil {
                payloadJSON, _ := json.MarshalIndent(endpoint.PayloadExample, "", "    ")
                html += `<div class="code-example">
                    <h4>Request Example:</h4>
                    <pre><code class="language-json">` + string(payloadJSON) + `</code></pre>
                </div>`
            }
            
            if endpoint.ResponseExample != nil {
                responseJSON, _ := json.MarshalIndent(endpoint.ResponseExample, "", "    ")
                html += `<div class="code-example">
                    <h4>Response Example:</h4>
                    <pre><code class="language-json">` + string(responseJSON) + `</code></pre>
                </div>`
            }
            
            // Try it section
            html += `<div class="try-it">
                <h4><i class="fas fa-play-circle"></i> Try it out</h4>
                <div class="input-group">
                    <label>Path:</label>
                    <input type="text" class="path-input" value="` + endpoint.Path + `" placeholder="Enter path with parameters replaced">
                </div>`
            
            if endpoint.Auth {
                html += `<div class="auth-toggle">
                    <label>
                        <input type="checkbox" onclick="toggleAuthInput(this)"> Include authentication token
                    </label>
                    <div class="auth-input">
                        <input type="text" class="token-input" placeholder="Enter your JWT token">
                    </div>
                </div>`
            }
            
            if endpoint.Method == "POST" || endpoint.Method == "PUT" {
                var payloadStr string
                if endpoint.PayloadExample != nil {
                    payloadJSON, _ := json.MarshalIndent(endpoint.PayloadExample, "", "    ")
                    payloadStr = string(payloadJSON)
                }
                
                html += `<div class="input-group">
                    <label>Request Body:</label>
                    <textarea class="body-input" rows="5" placeholder="Enter request JSON">` + payloadStr + `</textarea>
                </div>`
            }
            
            html += `<button class="execute-btn" onclick="executeRequest(this, '` + endpoint.Method + `')">
                    <i class="fas fa-paper-plane"></i> Execute Request
                </button>
                
                <div class="response-container">
                    <h4>Response</h4>
                    <div class="response-info">
                        <span>Status:</span> <span class="status"></span>
                        <span class="time"></span>
                    </div>
                    <pre><code class="response-body language-json"></code></pre>
                </div>
            </div>`
            
            html += `</div>
            </div>`
        }
        
        html += `</div>
        </div>`
    }
    
    html += `</section>
    </main>
    
    <div class="scroll-top" onclick="scrollToTop()">
        <i class="fas fa-arrow-up"></i>
    </div>
    
    <footer>
        <div class="container">
            <p>` + doc.Name + ` &copy; 2023 | Version ` + doc.Version + `</p>
        </div>
    </footer>
    
    <script>
        // Syntax highlighting
        document.addEventListener('DOMContentLoaded', () => {
            document.querySelectorAll('pre code').forEach((block) => {
                hljs.highlightBlock(block);
            });
            
            // Show scroll to top button when scrolled down
            window.addEventListener('scroll', () => {
                const scrollBtn = document.querySelector('.scroll-top');
                if (window.pageYOffset > 300) {
                    scrollBtn.classList.add('show');
                } else {
                    scrollBtn.classList.remove('show');
                }
            });
        });
        
        // Toggle endpoint details
        function toggleEndpoint(element) {
            const body = element.nextElementSibling;
            body.classList.toggle('active');
        }
        
        // Filter categories
        function filterCategory(category) {
            const buttons = document.querySelectorAll('.category-btn');
            const categories = document.querySelectorAll('.category');
            
            // Update active button
            buttons.forEach(btn => {
                btn.classList.remove('active');
                if (btn.textContent.toLowerCase() === category.toLowerCase() || 
                   (category === 'all' && btn.textContent.toLowerCase() === 'all')) {
                    btn.classList.add('active');
                }
            });
            
            // Show/hide categories
            if (category === 'all') {
                categories.forEach(cat => {
                    cat.style.display = 'block';
                });
            } else {
                categories.forEach(cat => {
                    if (cat.dataset.category === category) {
                        cat.style.display = 'block';
                    } else {
                        cat.style.display = 'none';
                    }
                });
            }
        }
        
        // Copy to clipboard
        function copyToClipboard(elementId) {
            const element = document.getElementById(elementId);
            const text = element.textContent.trim();
            
            navigator.clipboard.writeText(text).then(() => {
                const button = element.querySelector('.copy-btn');
                const originalText = button.textContent;
                button.textContent = 'Copied!';
                
                setTimeout(() => {
                    button.textContent = originalText;
                }, 2000);
            });
        }
        
        // Toggle auth input visibility
        function toggleAuthInput(checkbox) {
            const authInput = checkbox.parentElement.nextElementSibling;
            authInput.style.display = checkbox.checked ? 'block' : 'none';
        }
        
        // Execute API request
        function executeRequest(button, method) {
            const endpointBody = button.closest('.endpoint-body');
            const pathInput = endpointBody.querySelector('.path-input');
            const responseContainer = endpointBody.querySelector('.response-container');
            const responseStatus = responseContainer.querySelector('.status');
            const responseTime = responseContainer.querySelector('.time');
            const responseBody = responseContainer.querySelector('.response-body');
            
            // Get path (replace parameters if needed)
            let path = pathInput.value;
            if (!path.startsWith('/')) {
                path = '/' + path;
            }
            
            // Create request options
            const options = {
                method: method,
                headers: {
                    'Content-Type': 'application/json'
                }
            };
            
            // Add auth token if provided
            const tokenCheckbox = endpointBody.querySelector('.auth-toggle input[type="checkbox"]');
            if (tokenCheckbox && tokenCheckbox.checked) {
                const tokenInput = endpointBody.querySelector('.token-input');
                if (tokenInput.value.trim()) {
                    options.headers['Authentication'] = 'Bearer ' + tokenInput.value.trim();
                }
            }
            
            // Add request body for POST/PUT
            if (method === 'POST' || method === 'PUT') {
                const bodyInput = endpointBody.querySelector('.body-input');
                if (bodyInput && bodyInput.value.trim()) {
                    try {
                        options.body = JSON.parse(bodyInput.value);
                        options.body = JSON.stringify(options.body);
                    } catch (e) {
                        alert('Invalid JSON in request body');
                        return;
                    }
                }
            }
            
            // Show loading state
            button.innerHTML = '<div class="loading"></div> Loading...';
            button.disabled = true;
            
            // Execute the request
            const startTime = new Date();
            
            fetch('` + doc.BaseURL + `' + path, options)
                .then(response => {
                    const endTime = new Date();
                    const duration = endTime - startTime;
                    
                    // Update status
                    responseStatus.textContent = response.status + ' ' + response.statusText;
                    responseStatus.className = 'status';
                    responseStatus.classList.add(response.ok ? 'success' : 'error');
                    
                    // Update time
                    responseTime.textContent = '(' + duration + 'ms)';
                    
                    return response.text();
                })
                .then(text => {
                    // Try to parse as JSON
                    try {
                        const json = JSON.parse(text);
                        responseBody.textContent = JSON.stringify(json, null, 2);
                    } catch (e) {
                        // Not JSON, show as is
                        responseBody.textContent = text;
                    }
                    
                    // Highlight response
                    hljs.highlightBlock(responseBody);
                    
                    // Show response
                    responseContainer.classList.add('show');
                })
                .catch(error => {
                    responseStatus.textContent = 'Error';
                    responseStatus.className = 'status error';
                    responseBody.textContent = error.message;
                    responseContainer.classList.add('show');
                })
                .finally(() => {
                    // Reset button
                    button.innerHTML = '<i class="fas fa-paper-plane"></i> Execute Request';
                    button.disabled = false;
                });
        }
        
        // Scroll to top
        function scrollToTop() {
            window.scrollTo({
                top: 0,
                behavior: 'smooth'
            });
        }
        
        // Toggle dark/light theme
        function toggleTheme() {
            const body = document.body;
            const themeIcon = document.querySelector('.theme-toggle i');
            
            body.classList.toggle('dark-theme');
            
            if (body.classList.contains('dark-theme')) {
                themeIcon.className = 'fas fa-sun';
            } else {
                themeIcon.className = 'fas fa-moon';
            }
            
            // Store preference
            localStorage.setItem('theme', body.classList.contains('dark-theme') ? 'dark' : 'light');
        }
        
        // Apply saved theme preference
        (() => {
            const savedTheme = localStorage.getItem('theme');
            if (savedTheme === 'dark') {
                document.body.classList.add('dark-theme');
                document.querySelector('.theme-toggle i').className = 'fas fa-sun';
            }
        })();
    </script>
</body>
</html>`

    w.Header().Set("Content-Type", "text/html")
    w.Write([]byte(html))
}

func formatAuthExamples(examples map[string]interface{}) string {
    exampleJSON, _ := json.MarshalIndent(examples, "", "    ")
    return string(exampleJSON)
}