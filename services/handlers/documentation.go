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
        // BaseURL:     "http://localhost:40331",
        Authorization: AuthInfo{
            Description: "The API uses JWT (JSON Web Token) for authorization. Protected endpoints require a valid JWT token in the request header.",
            Method:      "Bearer Token Authorization",
            Header:      "Authorization: Bearer your_jwt_token",
            Examples: map[string]interface{}{
                "loginRequest": map[string]interface{}{
                    "url": "POST /login",
                    "body": map[string]string{
                        "username": "coffeegeek",
                        "passwords": "kawa123",
                    },
                },
                "loginResponse": map[string]interface{}{
                    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
                },
                "authorizedRequest": map[string]interface{}{
                    "url": "POST /coffees",
                    "headers": map[string]string{
                        "Content-Type": "application/json",
                        "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
                    },
                },
                "roles": map[string]string{
                    "user": "Regular user with access to create/update content",
                    "admin": "Administrator with additional privileges like deletion",
                },
            },
        },
        Endpoints: map[string][]API{
            "Authorization": {
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
                        "username": "coffeegeek",
                        "passwords": "kawa123",
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
        // BaseURL:     "http://localhost:40331",
        Authorization: AuthInfo{
            Description: "The API uses JWT (JSON Web Token) for Authorization. Protected endpoints require a valid JWT token in the request header.",
            Method:      "Bearer Token Authorization",
            Header:      "Authorization: Bearer your_jwt_token",
            Examples: map[string]interface{}{
                "loginRequest": map[string]interface{}{
                    "url": "POST /login",
                    "body": map[string]string{
                        "username": "coffeeuser",
                        "passwords": "strongpassword",
                    },
                },
                "loginResponse": map[string]interface{}{
                    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
                },
                "authorizedRequest": map[string]interface{}{
                    "url": "POST /coffees",
                    "headers": map[string]string{
                        "Content-Type": "application/json",
                        "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
                    },
                },
                "roles": map[string]string{
                    "user": "Regular user with access to create/update content",
                    "admin": "Administrator with additional privileges like deletion",
                },
            },
        },
        Endpoints: map[string][]API{
            "Authorization": {
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
                        "username": "coffeegeek",
                        "passwords": "kawa123",
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
                        "name":         "Geisha Kolumbia Finca",
                        "roasteryId":   1,
                        "country":      "Colombia",
                        "region":       "Huila",
                        "farm":         "Monteblanco",
                        "variety":      "Geisha",
                        "process":      "Washed",
                        "roastProfile": "Light",
                        "flavourNotes": []string{"citrus", "almond", "white chocolate"},
                        "description":  "Kawa Geisha Kolumbia Finca “Monteblanco” Rodrigo Sanchez Valencia to pozycja obowiązkowa dla każdego kawowego entuzjasty, poszukującego w kawie wyjątkowych smaków i aromatu. Ziarna kawy pochodzą z bajecznego regionu Kolumbii  Finca “Monteblanco”, gdzie gleba i klimat stanowią wymarzone warunki do uprawy kawy najwyższej jakości. Plantacja “Monteblanco”  położona 1700 m.n.p.m należy do Rodrigo Sancheza Valenci, który słynie z tego że poświecą sporo uwagi każdemu etapowi produkcji kawy. Stosuje zrównoważone metody uprawy, które obejmują ręczne zbieranie dojrzałych ziaren, oraz staranny proces suszenia i palenia ziaren. Dzięki temu kawa Geisha Kolumbia Finca “Monteblanco” osiąga najwyższy poziom. Kawa Geisha Kolumbia Finca “Monteblanco” Rodrigo Sanchez Valencia charakteryzuje się bogatą paletą smaków.  Wyczuwalne są delikatne nuty cytrusowe, jaśminu oraz subtelna słodycz. W smaku można również wyczuć nuty mlecznej czekolady oraz  migdałów i subtelnej goryczki, które nadają jej wyjątkowy charakter.",
                        "imageUrl":     "https://palarniagrunt.pl/wp-content/uploads/2023/06/Geisha.jpg",
                    },
                },
                {
                    Method:      "PUT",
                    Path:        "/coffees/{id}",
                    Description: "Update a coffee",
                    PayloadExample: map[string]interface{}{
                        "name":         "Geisha Kolumbia Finca",
                        "roasteryId":   1,
                        "country":      "Colombia",
                        "region":       "Huila",
                        "farm":         "Monteblanco",
                        "variety":      "Geisha",
                        "process":      "Washed",
                        "roastProfile": "Light",
                        "flavourNotes": []string{"citrus", "almond", "white chocolate"},
                        "description":  "Kawa Geisha Kolumbia Finca “Monteblanco” Rodrigo Sanchez Valencia to pozycja obowiązkowa dla każdego kawowego entuzjasty, poszukującego w kawie wyjątkowych smaków i aromatu. Ziarna kawy pochodzą z bajecznego regionu Kolumbii  Finca “Monteblanco”, gdzie gleba i klimat stanowią wymarzone warunki do uprawy kawy najwyższej jakości. Plantacja “Monteblanco”  położona 1700 m.n.p.m należy do Rodrigo Sancheza Valenci, który słynie z tego że poświecą sporo uwagi każdemu etapowi produkcji kawy. Stosuje zrównoważone metody uprawy, które obejmują ręczne zbieranie dojrzałych ziaren, oraz staranny proces suszenia i palenia ziaren. Dzięki temu kawa Geisha Kolumbia Finca “Monteblanco” osiąga najwyższy poziom. Kawa Geisha Kolumbia Finca “Monteblanco” Rodrigo Sanchez Valencia charakteryzuje się bogatą paletą smaków.  Wyczuwalne są delikatne nuty cytrusowe, jaśminu oraz subtelna słodycz. W smaku można również wyczuć nuty mlecznej czekolady oraz  migdałów i subtelnej goryczki, które nadają jej wyjątkowy charakter.",
                        "imageUrl":     "https://palarniagrunt.pl/wp-content/uploads/2023/06/Geisha.jpg",
                    },
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
                        "name":        "Fuszera",
                        "country":     "Poland",
                        "city":        "Katowice",
                        "address":     "Korfantego 72",
                        "website":     "https://fuszera.pl/",
                        "description": "Najlepsza, a przynajmniej najzabawniejsza palarnia w polsce",
                        "imageUrl":    "https://fuszera.pl/wp-content/uploads/2025/01/P1449314-1-scaled.webp",
                    },
                },
                {
                    Method:      "PUT",
                    Path:        "/roasteries/{id}",
                    Description: "Update a roastery",
                    PayloadExample: map[string]interface{}{
                        "name":        "Fuszera",
                        "country":     "Poland",
                        "city":        "Katowice",
                        "address":     "Korfantego 72",
                        "website":     "https://fuszera.pl/",
                        "description": "Najlepsza, a przynajmniej najzabawniejsza palarnia w polsce",
                        "imageUrl":    "https://fuszera.pl/wp-content/uploads/2025/01/P1449314-1-scaled.webp",
                    },
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
                        "name":        "Bez cukru",
                        "country":     "Poland",
                        "city":        "Katowice",
                        "address":     "Wawelska 1",
                        "website":     "http://www.kawiarniabezcukru.pl",
                        "description": "Dobra kawa w centrum Katowic",
                        "imageUrl":    "https://bi.im-g.pl/im/e4/47/14/z21265124AMP,Kawiarnia--Bez-cukru-w-Katowicach.jpg",
                    },
                },
                {
                    Method:      "PUT",
                    Path:        "/shops/{id}",
                    Description: "Update a coffee shop",
                    PayloadExample: map[string]interface{}{
                        "name":        "Bez cukru",
                        "country":     "Poland",
                        "city":        "Katowice",
                        "address":     "Wawelska 1",
                        "website":     "http://www.kawiarniabezcukru.pl",
                        "description": "Dobra kawa w centrum Katowic",
                        "imageUrl":    "https://bi.im-g.pl/im/e4/47/14/z21265124AMP,Kawiarnia--Bez-cukru-w-Katowicach.jpg",
                    },
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
                        "rating":   4,
                        "review":   "Great coffee with bright acidity and floral notes",
                    },
                },
                {
                    Method:      "PUT",
                    Path:        "/reviews/{id}",
                    Description: "Update a review",
                    PayloadExample: map[string]interface{}{
                        "rating": 5,
                        "review": "Exceptional coffee with chocolate and hazelnut notes, perfectly balanced acidity.",
                    },
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
            background: linear-gradient(135deg, #2c1e18, #0f0906);
            color: white;
            padding: 60px 0;
            text-align: center;
            position: relative;
            overflow: hidden;
            box-shadow: 0 4px 12px rgba(0,0,0,0.3);
        }
        
        header::before {
            content: '';
            position: absolute;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            background-image: url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%239C6F48' fill-opacity='0.1'%3E%3Cpath d='M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E");
            opacity: 0.2;
        }
        
        .coffee-beans {
            position: absolute;
            width: 100%;
            height: 100%;
            top: 0;
            left: 0;
            overflow: hidden;
            z-index: 1;
        }
        
        .bean {
            position: absolute;
            background: radial-gradient(ellipse at center, rgba(156, 111, 72, 0.3) 0%, rgba(156, 111, 72, 0) 70%);
            border-radius: 50%;
            transform: scale(0);
            animation: float-up 15s infinite ease-out;
        }
        
        @keyframes float-up {
            0% {
                transform: translateY(100%) scale(0);
                opacity: 0;
            }
            10% {
                opacity: 0.8;
                transform: translateY(80%) scale(1);
            }
            100% {
                transform: translateY(-20%) scale(0.5);
                opacity: 0;
            }
        }
        
        .header-content {
            position: relative;
            z-index: 2;
        }
        
        .logo-container {
            margin-bottom: 25px;
            display: inline-block;
            position: relative;
        }
        
        .coffee-cup {
            font-size: 4.5rem;
            display: inline-block;
            color: #d4a574;
            text-shadow: 0 2px 10px rgba(0,0,0,0.3);
            animation: glow 3s infinite ease-in-out;
        }
        
        @keyframes glow {
            0%, 100% { 
                color: #d4a574;
                text-shadow: 0 2px 10px rgba(0,0,0,0.3);
            }
            50% { 
                color: #e8c496;
                text-shadow: 0 2px 20px rgba(212, 165, 116, 0.6);
            }
        }
        
        .steam {
            position: absolute;
            top: -15px;
            left: 50%;
            transform: translateX(-50%);
            width: 60px;
            height: 30px;
            background: transparent;
        }
        
        .steam-particle {
            position: absolute;
            bottom: 0;
            width: 8px;
            height: 8px;
            border-radius: 50%;
            background-color: rgba(255, 255, 255, 0.6);
            opacity: 0;
            animation: rise 3s infinite ease-in-out;
        }
        
        .steam-particle:nth-child(1) {
            left: 20%;
            animation-delay: 0.2s;
        }
        
        .steam-particle:nth-child(2) {
            left: 40%;
            animation-delay: 0.8s;
        }
        
        .steam-particle:nth-child(3) {
            left: 60%;
            animation-delay: 0.4s;
        }
        
        .steam-particle:nth-child(4) {
            left: 80%;
            animation-delay: 1s;
        }
        
        @keyframes rise {
            0% {
                bottom: 0;
                opacity: 0;
                width: 8px;
                height: 8px;
            }
            30% {
                opacity: 0.8;
            }
            100% {
                bottom: 100%;
                opacity: 0;
                width: 20px;
                height: 20px;
            }
        }
        
        h1 {
            font-size: 3.5rem;
            font-weight: 700;
            letter-spacing: 2px;
            margin-bottom: 15px;
            text-shadow: 0 2px 15px rgba(0,0,0,0.3);
            background: linear-gradient(45deg, #e8c496, #ffffff);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            position: relative;
            display: inline-block;
        }
        
        h1::after {
            content: '';
            position: absolute;
            bottom: -10px;
            left: 50%;
            transform: translateX(-50%);
            width: 120px;
            height: 3px;
            background: linear-gradient(90deg, rgba(255,255,255,0), rgba(255,255,255,0.8), rgba(255,255,255,0));
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
            max-height: none;
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
            overflow: auto; /* This enables both horizontal and vertical scrolling as needed */
            margin: 10px 0;
            position: relative;
            max-height: none; /* Remove the max-height restriction */
        }
            
        .response-body {
            white-space: pre;
            word-wrap: normal;
            max-height: none; /* Remove the height restriction */
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
            background-color: rgba(212, 165, 116, 0.05);
            padding: 12px;
            border-radius: 5px;
            border-left: 3px solid #d4a574;
            animation: fadeIn 0.3s ease-out;
            box-shadow: inset 0 1px 3px rgba(0,0,0,0.05);
            transition: all 0.3s ease;
        }
        
        .auth-input.visible {
            display: block;
        }
        
        .token-input {
            width: 100%;
            padding: 10px;
            border: 1px solid rgba(212, 165, 116, 0.3);
            border-radius: 5px;
            font-family: monospace;
            background-color: rgba(255, 255, 255, 0.9);
            transition: all 0.3s ease;
        }
        
        .token-input:focus {
            border-color: #d4a574;
            outline: none;
            box-shadow: 0 0 0 2px rgba(212, 165, 116, 0.2);
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
            background: linear-gradient(135deg, #1c140f, #0c0805);
            color: white;
            text-align: center;
            padding: 60px 0 40px;
            margin-top: 80px;
            position: relative;
            overflow: hidden;
            box-shadow: 0 -4px 15px rgba(0,0,0,0.4);
            perspective: 1000px;
        }
        
        footer::before {
            content: '';
            position: absolute;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            background-image: url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%239C6F48' fill-opacity='0.1'%3E%3Cpath d='M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E");
            opacity: 0.15;
            animation: pulse-bg 8s infinite alternate;
        }
        
        @keyframes pulse-bg {
            0% { opacity: 0.1; }
            100% { opacity: 0.2; }
        }
        
        .footer-wave {
            position: absolute;
            top: 0;
            left: 0;
            width: 100%;
            height: 60px;
            transform: translateY(-95%);
            filter: drop-shadow(0 -5px 5px rgba(0,0,0,0.1));
        }
        
        .wave {
            position: absolute;
            height: 60px;
            width: 100%;
            background: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 1200 120' preserveAspectRatio='none'%3E%3Cpath d='M0,0V46.29c47.79,22.2,103.59,32.17,158,28,70.36-5.37,136.33-33.31,206.8-37.5C438.64,32.43,512.34,53.67,583,72.05c69.27,18,138.3,24.88,209.4,13.08,36.15-6,69.85-17.84,104.45-29.34C989.49,25,1113-14.29,1200,52.47V0Z' opacity='.25' fill='%231c140f'%3E%3C/path%3E%3Cpath d='M0,0V15.81C13,36.92,27.64,56.86,47.69,72.05,99.41,111.27,165,111,224.58,91.58c31.15-10.15,60.09-26.07,89.67-39.8,40.92-19,84.73-46,130.83-49.67,36.26-2.85,70.9,9.42,98.6,31.56,31.77,25.39,62.32,62,103.63,73,40.44,10.79,81.35-6.69,119.13-24.28s75.16-39,116.92-43.05c59.73-5.85,113.28,22.88,168.9,38.84,30.2,8.66,59,6.17,87.09-7.5,22.43-10.89,48-26.93,60.65-49.24V0Z' opacity='.5' fill='%231c140f'%3E%3C/path%3E%3Cpath d='M0,0V5.63C149.93,59,314.09,71.32,475.83,42.57c43-7.64,84.23-20.12,127.61-26.46,59-8.63,112.48,12.24,165.56,35.4C827.93,77.22,886,95.24,951.2,90c86.53-7,172.46-45.71,248.8-84.81V0Z' fill='%231c140f'%3E%3C/path%3E%3C/svg%3E") no-repeat;
            background-size: cover;
            animation: wave 15s linear infinite;
        }
        
        .wave:nth-child(1) {
            z-index: 3;
            opacity: 0.7;
            animation: wave-move1 12s linear infinite;
        }
        
        .wave:nth-child(2) {
            z-index: 2;
            opacity: 0.5;
            animation: wave-move2 18s linear infinite;
            bottom: 10px;
        }
        
        .wave:nth-child(3) {
            z-index: 1;
            opacity: 0.3;
            animation: wave-move3 24s linear infinite;
            bottom: 15px;
        }
        
        @keyframes wave-move1 {
            0% { background-position-x: 0; }
            100% { background-position-x: 1200px; }
        }
        
        @keyframes wave-move2 {
            0% { background-position-x: 0; }
            100% { background-position-x: -1200px; }
        }
        
        @keyframes wave-move3 {
            0% { background-position-x: 0; }
            100% { background-position-x: 1500px; }
        }
        
        .footer-content {
            position: relative;
            z-index: 5;
            transform-style: preserve-3d;
        }
        
        .footer-beans {
            position: absolute;
            width: 100%;
            height: 100%;
            top: 0;
            left: 0;
            overflow: hidden;
            z-index: 1;
            transform-style: preserve-3d;
        }
        
        .footer-bean {
            position: absolute;
            background: radial-gradient(ellipse at center, rgba(156, 111, 72, 0.2) 0%, rgba(156, 111, 72, 0) 70%);
            border-radius: 50%;
            transform: translateZ(0) rotateX(45deg) scale(0);
            animation: float-beans 20s infinite ease-out;
            bottom: 0;
            filter: blur(1px);
        }
        
        @keyframes float-beans {
            0% {
                transform: translateY(20px) translateZ(0) rotateX(45deg) scale(0) rotate(0deg);
                opacity: 0;
            }
            5% {
                opacity: 0.8;
                transform: translateY(15px) translateZ(20px) rotateX(45deg) scale(1) rotate(45deg);
            }
            90% {
                opacity: 0.3;
                transform: translateY(-140px) translateZ(100px) rotateX(45deg) scale(0.7) rotate(300deg);
            }
            100% {
                transform: translateY(-150px) translateZ(0) rotateX(45deg) scale(0) rotate(360deg);
                opacity: 0;
            }
        }
        
        .coffee-trails {
            position: absolute;
            width: 100%;
            height: 70%;
            bottom: 0;
            left: 0;
            overflow: hidden;
            perspective: 500px;
        }
        
        .coffee-trail {
            position: absolute;
            bottom: -50px;
            width: 2px;
            background: linear-gradient(to bottom, rgba(212, 165, 116, 0), rgba(212, 165, 116, 0.5) 50%, rgba(212, 165, 116, 0));
            animation: rise-trail 6s ease-in-out infinite;
            opacity: 0;
            transform-style: preserve-3d;
            transform-origin: bottom center;
        }
        
        @keyframes rise-trail {
            0% {
                height: 0;
                opacity: 0;
                transform: translateZ(0) rotateY(0deg);
                filter: blur(0);
            }
            20% {
                opacity: 0.9;
                transform: translateZ(20px) rotateY(15deg);
                filter: blur(1px);
            }
            80% {
                opacity: 0.4;
                transform: translateZ(50px) rotateY(-15deg);
                filter: blur(2px);
            }
            100% {
                height: 100%;
                opacity: 0;
                transform: translateZ(0) rotateY(0deg);
                filter: blur(0);
            }
        }
        
        .footer-logo {
            margin-bottom: 20px;
            display: inline-block;
            position: relative;
            animation: float-logo 6s ease-in-out infinite;
        }
        
        @keyframes float-logo {
            0%, 100% { transform: translateY(0); }
            50% { transform: translateY(-10px); }
        }
        
        .footer-logo i {
            font-size: 2.5rem;
            color: #d4a574;
            display: inline-block;
            animation: glow-logo 4s infinite alternate;
        }
        
        @keyframes glow-logo {
            0% { 
                color: #d4a574;
                text-shadow: 0 0 5px rgba(212, 165, 116, 0.3);
                transform: scale(1);
            }
            100% { 
                color: #e8c496; 
                text-shadow: 0 0 20px rgba(212, 165, 116, 0.8), 0 0 30px rgba(212, 165, 116, 0.4);
                transform: scale(1.1);
            }
        }
        
        .footer-content h3 {
            color: white;
            margin-bottom: 15px;
            text-align: center;
            font-size: 2rem;
            font-weight: 700;
            background: linear-gradient(45deg, #d4a574, #ffffff);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            position: relative;
            display: inline-block;
            animation: shimmer 10s infinite;
        }
        
        @keyframes shimmer {
            0% { background-position: -100% 0; }
            100% { background-position: 200% 0; }
        }
        
        .footer-content h3::after {
            content: '';
            position: absolute;
            bottom: -8px;
            left: 50%;
            transform: translateX(-50%);
            width: 50%;
            height: 2px;
            background: linear-gradient(90deg, rgba(255,255,255,0), rgba(212, 165, 116, 0.8), rgba(255,255,255,0));
            animation: width-pulse 4s infinite alternate;
        }
        
        @keyframes width-pulse {
            0% { width: 30%; opacity: 0.5; }
            100% { width: 70%; opacity: 1; }
        }
        
        .footer-content p {
            color: rgba(255,255,255,0.7);
            max-width: 600px;
            margin: 0 auto 20px;
            line-height: 1.6;
            transform: translateZ(10px);
        }
        
        .footer-social {
            margin-top: 30px;
            position: relative;
            z-index: 5;
            perspective: 1500px;
            transform-style: preserve-3d;
        }
        
        .social-icon {
            display: inline-flex;
            justify-content: center;
            align-items: center;
            width: 40px;
            height: 40px;
            border-radius: 50%;
            margin: 0 12px;
            background: rgba(255,255,255,0.1);
            color: white;
            text-decoration: none;
            transition: all 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275);
            position: relative;
            overflow: visible;
            transform-style: preserve-3d;
            transform: translateZ(0);
            box-shadow: 0 5px 15px rgba(0,0,0,0.2);
            animation: social-float 6s infinite ease-in-out;
            animation-delay: calc(var(--i) * 0.2s);
            z-index: 10;
            will-change: transform;
            transform-origin: center center;
        }
        
        .social-icon:hover {
            background: linear-gradient(135deg, #d4a574, #b87333);
            transform: translateZ(80px) scale(1.4) rotate(8deg);
            box-shadow: 0 20px 40px rgba(212, 165, 116, 0.6), 0 0 30px rgba(212, 165, 116, 0.4);
            z-index: 20;
        }
        
        .social-icon i {
            position: relative;
            z-index: 22;
            transform: translateZ(10px);
            transition: all 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275);
        }
        
        .social-icon:hover i {
            color: white;
            transform: translateZ(30px) scale(1.3);
            text-shadow: 0 0 15px rgba(255,255,255,0.7);
        }
        
        .social-icon::before {
            content: '';
            position: absolute;
            width: 160%;
            height: 160%;
            top: -30%;
            left: -30%;
            background: radial-gradient(circle, rgba(255,255,255,0.4) 0%, rgba(255,255,255,0) 70%);
            transform: scale(0);
            opacity: 0;
            transition: all 0.5s;
            z-index: 21;
            pointer-events: none;
        }
        
        .social-icon:hover::before {
            transform: scale(1);
            opacity: 1;
            animation: social-pulse 1.5s infinite;
        }
        
        @keyframes social-pulse {
            0% { transform: scale(0.8); opacity: 0.8; }
            100% { transform: scale(1.5); opacity: 0; }
        }
        
        .copyright {
            margin-top: 30px;
            color: rgb(255, 255, 255);
            font-size: 0.9rem;
            position: relative;
            display: inline-block;
            padding: 10px 20px;
            border-radius: 30px;
            background: transparent;
            backdrop-filter: blur(5px);
            transform: translateZ(5px);
        }
        
        .copyright::before,
        .copyright::after {
            content: '';
            position: absolute;
            height: 1px;
            width: 70px;
            background: linear-gradient(90deg, transparent, rgba(212, 165, 116, 0.5), transparent);
            top: 50%;
            animation: width-pulse 4s infinite alternate-reverse;
        }
        
        .copyright::before {
            left: -90px;
        }
        
        .copyright::after {
            right: -90px;
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

        .wrapper {
        display: inline-flex;
        list-style: none;
        height: 120px;
        width: 100%;
        padding-top: 40px;
        font-family: "Poppins", sans-serif;
        justify-content: center;
        }

        .wrapper .icon {
        position: relative;
        background: transparent;  /* Zmienione z #fff na transparent */
        color: #d4a574;  /* Dodany kolor ikon pasujący do motywu kawy */
        border-radius: 50%;
        margin: 10px;
        width: 50px;
        height: 50px;
        font-size: 18px;
        display: flex;
        justify-content: center;
        align-items: center;
        flex-direction: column;
        box-shadow: 0 5px 10px rgba(0, 0, 0, 0.2);  /* Delikatniejszy cień */
        cursor: pointer;
        transition: all 0.2s cubic-bezier(0.68, -0.55, 0.265, 1.55);
        border: 1px solid rgba(212, 165, 116, 0.3);  /* Delikatna obwódka */
        }

        .wrapper .icon a {
            display: flex;
            justify-content: center;
            align-items: center;
            color: #fff;
        }

        .wrapper .tooltip {
        position: absolute;
        top: 0;
        font-size: 14px;
        background: #fff;
        color: #fff;
        padding: 5px 8px;
        border-radius: 5px;
        box-shadow: 0 10px 10px rgba(0, 0, 0, 0.1);
        opacity: 0;
        pointer-events: none;
        transition: all 0.3s cubic-bezier(0.68, -0.55, 0.265, 1.55);
        }

        .wrapper .tooltip::before {
        position: absolute;
        content: "";
        height: 8px;
        width: 8px;
        background: #fff;
        bottom: -3px;
        left: 50%;
        transform: translate(-50%) rotate(45deg);
        transition: all 0.3s cubic-bezier(0.68, -0.55, 0.265, 1.55);
        }

        .wrapper .icon:hover .tooltip {
        top: -45px;
        opacity: 1;
        visibility: visible;
        pointer-events: auto;
        }

        .wrapper .icon:hover span,
        .wrapper .icon:hover .tooltip {
        text-shadow: 0px -1px 0px rgba(0, 0, 0, 0.1);
        }

        .wrapper .facebook:hover,
        .wrapper .facebook:hover .tooltip,
        .wrapper .facebook:hover .tooltip::before {
        background: #333;
        color: #fff;
        }

        .wrapper .twitter:hover,
        .wrapper .twitter:hover .tooltip,
        .wrapper .twitter:hover .tooltip::before {
        background: #0077b5;
        color: #fff;
        }

        .wrapper .instagram:hover,
        .wrapper .instagram:hover .tooltip,
        .wrapper .instagram:hover .tooltip::before {
        background: #e4405f;
        color: #fff;
        }
    </style>
</head>
<body>
    <div class="theme-toggle" onclick="toggleTheme()">
        <i class="fas fa-moon"></i>
    </div>
    
    <header>
        <div class="coffee-beans">
            <div class="bean" style="left: 10%; width: 120px; height: 120px; animation-duration: 12s;"></div>
            <div class="bean" style="left: 25%; width: 80px; height: 80px; animation-duration: 18s; animation-delay: 2s;"></div>
            <div class="bean" style="left: 40%; width: 100px; height: 100px; animation-duration: 15s; animation-delay: 4s;"></div>
            <div class="bean" style="left: 60%; width: 110px; height: 110px; animation-duration: 14s; animation-delay: 1s;"></div>
            <div class="bean" style="left: 75%; width: 90px; height: 90px; animation-duration: 17s; animation-delay: 3s;"></div>
            <div class="bean" style="left: 85%; width: 130px; height: 130px; animation-duration: 13s; animation-delay: 0.5s;"></div>
        </div>
        <div class="container header-content">
            <div class="logo-container">
 
                <div class="coffee-cup">
                    <i class="fas fa-mug-hot"></i>
                </div>
            </div>
            <h1>CoffeeDB API</h1>
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
        
        <section class="section" id="Authorization">
            <h2><i class="fas fa-lock"></i> Authorization</h2>
            <div class="auth-info">
                <p>` + doc.Authorization.Description + `</p>
                <div><strong>Method:</strong> ` + doc.Authorization.Method + `</div>
                <div><strong>Header:</strong> <code>` + doc.Authorization.Header + `</code></div>
                
                <h4>Authorization Examples</h4>
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
        "Authorization": "fas fa-lock",
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
                        <input type="checkbox" onclick="toggleAuthInput(this)"> Include Authorization token
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
        <div class="footer-wave">
            <div class="wave"></div>
            <div class="wave"></div>
            <div class="wave"></div>
        </div>
        
        <div class="footer-beans">
            <div class="footer-bean" style="left: 5%; width: 60px; height: 60px; animation-duration: 15s;"></div>
            <div class="footer-bean" style="left: 15%; width: 40px; height: 40px; animation-duration: 18s; animation-delay: 1s;"></div>
            <div class="footer-bean" style="left: 30%; width: 70px; height: 70px; animation-duration: 12s; animation-delay: 2s;"></div>
            <div class="footer-bean" style="left: 50%; width: 50px; height: 50px; animation-duration: 16s; animation-delay: 0.5s;"></div>
            <div class="footer-bean" style="left: 65%; width: 45px; height: 45px; animation-duration: 14s; animation-delay: 1.5s;"></div>
            <div class="footer-bean" style="left: 80%; width: 55px; height: 55px; animation-duration: 17s; animation-delay: 2.5s;"></div>
            <div class="footer-bean" style="left: 92%; width: 60px; height: 60px; animation-duration: 13s; animation-delay: 1s;"></div>
        </div>
        
        <div class="coffee-trails">
            <div class="coffee-trail" style="left: 10%; animation-delay: 0.2s;"></div>
            <div class="coffee-trail" style="left: 25%; animation-delay: 2.1s;"></div>
            <div class="coffee-trail" style="left: 40%; animation-delay: 0.5s;"></div>
            <div class="coffee-trail" style="left: 55%; animation-delay: 1.8s;"></div>
            <div class="coffee-trail" style="left: 70%; animation-delay: 1.2s;"></div>
            <div class="coffee-trail" style="left: 85%; animation-delay: 0.8s;"></div>
            <div class="coffee-trail" style="left: 18%; animation-delay: 1.4s;"></div>
            <div class="coffee-trail" style="left: 33%; animation-delay: 2.7s;"></div>
            <div class="coffee-trail" style="left: 48%; animation-delay: 0.9s;"></div>
            <div class="coffee-trail" style="left: 63%; animation-delay: 2.3s;"></div>
            <div class="coffee-trail" style="left: 78%; animation-delay: 1.6s;"></div>
            <div class="coffee-trail" style="left: 93%; animation-delay: 3.0s;"></div>
        </div>
        
        <div class="container footer-content">
            <div class="footer-logo">
                <i class="fas fa-mug-hot"></i>
            </div>
            
            <h3>CoffeeDB API</h3>
            
            <p>The ultimate RESTful API for coffee enthusiasts and professionals. Access data on coffees, roasteries, shops and more.</p>
            
            <div class="footer-social">
            <!-- From Uiverse.io by david-mohseni --> 
            <ul class="wrapper">
                <li class="icon facebook">
                <span class="tooltip">GitHub</span>
                <a href="https://github.com/PanPeryskop">
                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 496 512" height="1.2em" fill="currentColor">
                    <path d="M165.9 397.4c0 2-2.3 3.6-5.2 3.6-3.3.3-5.6-1.3-5.6-3.6 0-2 2.3-3.6 5.2-3.6 3-.3 5.6 1.3 5.6 3.6zm-31.1-4.5c-.7 2 1.3 4.3 4.3 4.9 2.6 1 5.6 0 6.2-2s-1.3-4.3-4.3-5.2c-2.6-.7-5.5.3-6.2 2.3zm44.2-1.7c-2.9.7-4.9 2.6-4.6 4.9.3 2 2.9 3.3 5.9 2.6 2.9-.7 4.9-2.6 4.6-4.6-.3-1.9-3-3.2-5.9-2.9zM244.8 8C106.1 8 0 113.3 0 252c0 110.9 69.8 205.8 169.5 239.2 12.8 2.3 17.3-5.6 17.3-12.1 0-6.2-.3-40.4-.3-61.4 0 0-70 15-84.7-29.8 0 0-11.4-29.1-27.8-36.6 0 0-22.9-15.7 1.6-15.4 0 0 24.9 2 38.6 25.8 21.9 38.6 58.6 27.5 72.9 20.9 2.3-16 8.8-27.1 16-33.7-55.9-6.2-112.3-14.3-112.3-110.5 0-27.5 7.6-41.3 23.6-58.9-2.6-6.5-11.1-33.3 2.6-67.9 20.9-6.5 69 27 69 27 20-5.6 41.5-8.5 62.8-8.5s42.8 2.9 62.8 8.5c0 0 48.1-33.6 69-27 13.7 34.7 5.2 61.4 2.6 67.9 16 17.7 25.8 31.5 25.8 58.9 0 96.5-58.9 104.2-114.8 110.5 9.2 7.9 17 22.9 17 46.4 0 33.7-.3 75.4-.3 83.6 0 6.5 4.6 14.4 17.3 12.1C428.2 457.8 496 362.9 496 252 496 113.3 383.5 8 244.8 8zM97.2 352.9c-1.3 1-1 3.3.7 5.2 1.6 1.6 3.9 2.3 5.2 1 1.3-1 1-3.3-.7-5.2-1.6-1.6-3.9-2.3-5.2-1zm-10.8-8.1c-.7 1.3.3 2.9 2.3 3.9 1.6 1 3.6.7 4.3-.7.7-1.3-.3-2.9-2.3-3.9-2-.6-3.6-.3-4.3.7zm32.4 35.6c-1.6 1.3-1 4.3 1.3 6.2 2.3 2.3 5.2 2.6 6.5 1 1.3-1.3.7-4.3-1.3-6.2-2.2-2.3-5.2-2.6-6.5-1zm-11.4-14.7c-1.6 1-1.6 3.6 0 5.9 1.6 2.3 4.3 3.3 5.6 2.3 1.6-1.3 1.6-3.9 0-6.2-1.4-2.3-4-3.3-5.6-2z"/>
                    </svg>
                </a>
                </li>
                <li class="icon twitter">
                <span class="tooltip">LinkedIn</span>
                <a href="https://www.linkedin.com/in/stanislaw-gadek/">
                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 448 512" height="1.2em" fill="currentColor">
                    <path d="M416 32H31.9C14.3 32 0 46.5 0 64.3v383.4C0 465.5 14.3 480 31.9 480H416c17.6 0 32-14.5 32-32.3V64.3c0-17.8-14.4-32.3-32-32.3zM135.4 416H69V202.2h66.5V416zm-33.2-243c-21.3 0-38.5-17.3-38.5-38.5S80.9 96 102.2 96c21.2 0 38.5 17.3 38.5 38.5 0 21.3-17.2 38.5-38.5 38.5zm282.1 243h-66.4V312c0-24.8-.5-56.7-34.5-56.7-34.6 0-39.9 27-39.9 54.9V416h-66.4V202.2h63.7v29.2h.9c8.9-16.8 30.6-34.5 62.9-34.5 67.2 0 79.7 44.3 79.7 101.9V416z"/>
                    </svg>
                </a>
                </li>
                <li class="icon instagram">
                <span class="tooltip">Instagram</span>
                <a href="https://www.instagram.com/_peryskop/">
                    <svg xmlns="http://www.w3.org/2000/svg" height="1.2em" fill="currentColor" class="bi bi-instagram" viewBox="0 0 16 16">
                    <path d="M8 0C5.829 0 5.556.01 4.703.048 3.85.088 3.269.222 2.76.42a3.917 3.917 0 0 0-1.417.923A3.927 3.927 0 0 0 .42 2.76C.222 3.268.087 3.85.048 4.7.01 5.555 0 5.827 0 8.001c0 2.172.01 2.444.048 3.297.04.852.174 1.433.372 1.942.205.526.478.972.923 1.417.444.445.89.719 1.416.923.51.198 1.09.333 1.942.372C5.555 15.99 5.827 16 8 16s2.444-.01 3.298-.048c.851-.04 1.434-.174 1.943-.372a3.916 3.916 0 0 0 1.416-.923c.445-.445.718-.891.923-1.417.197-.509.332-1.09.372-1.942C15.99 10.445 16 10.173 16 8s-.01-2.445-.048-3.299c-.04-.851-.175-1.433-.372-1.941a3.926 3.926 0 0 0-.923-1.417A3.911 3.911 0 0 0 13.24.42c-.51-.198-1.092-.333-1.943-.372C10.443.01 10.172 0 7.998 0h.003zm-.717 1.442h.718c2.136 0 2.389.007 3.232.046.78.035 1.204.166 1.486.275.373.145.64.319.92.599.28.28.453.546.598.92.11.281.24.705.275 1.485.039.843.047 1.096.047 3.231s-.008 2.389-.047 3.232c-.035.78-.166 1.203-.275 1.485a2.47 2.47 0 0 1-.599.919c-.28.28-.546.453-.92.598-.28.11-.704.24-1.485.276-.843.038-1.096.047-3.232.047s-2.39-.009-3.233-.047c-.78-.036-1.203-.166-1.485-.276a2.478 2.478 0 0 1-.92-.598 2.48 2.48 0 0 1-.6-.92c-.109-.281-.24-.705-.275-1.485-.038-.843-.046-1.096-.046-3.233 0-2.136.008-2.388.046-3.231.036-.78.166-1.204.276-1.486.145-.373.319-.64.599-.92.28-.28.546-.453.92-.598.282-.11.705-.24 1.485-.276.738-.034 1.024-.044 2.515-.045v.002zm4.988 1.328a.96.96 0 1 0 0 1.92.96.96 0 0 0 0-1.92zm-4.27 1.122a4.109 4.109 0 1 0 0 8.217 4.109 4.109 0 0 0 0-8.217zm0 1.441a2.667 2.667 0 1 1 0 5.334 2.667 2.667 0 0 1 0-5.334z"></path>
                    </svg>
                </a>
                </li>
            </ul>
            </div>
            
            <div class="copyright">
                &copy; 2025 CoffeeDB API | Version ` + doc.Version + `
            </div>
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
			const authToggle = checkbox.closest('.auth-toggle');
			if (!authToggle) {
				alert("Error: Could not find .auth-toggle element");
				return;
			}
			
			const authInput = authToggle.querySelector('.auth-input');
			if (!authInput) {
				alert("Error: Could not find .auth-input element");
				return;
			}
			
			// Toggle visibility
			if (checkbox.checked) {
				authInput.style.display = 'block';
			} else {
				authInput.style.display = 'none';
			}
			
			// Debug info
			console.log("Toggle state:", checkbox.checked);
			console.log("Auth input element:", authInput);
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
            const authToggle = endpointBody.querySelector('.auth-toggle');
            if (authToggle) {
                const tokenCheckbox = authToggle.querySelector('input[type="checkbox"]');
                if (tokenCheckbox && tokenCheckbox.checked) {
                    const tokenInput = authToggle.querySelector('.token-input');
                    if (tokenInput && tokenInput.value.trim()) {
                        options.headers['Authorization'] = 'Bearer ' + tokenInput.value.trim();
                    }
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