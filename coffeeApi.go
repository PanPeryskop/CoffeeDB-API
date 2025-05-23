package main

import (
    "fmt"
    "log"
    "net/http"
    
    "coffeeApi/services/db"
    "coffeeApi/services/handlers"
    "coffeeApi/services/middleware"
    
    "github.com/gorilla/mux"
)

func main() {
    if err := db.Init(); err != nil {
        log.Fatal("Błąd połączenia z bazą:", err)
    }
    
    router := mux.NewRouter()
    
    // Documentation
    router.HandleFunc("/", handlers.GetApiDocumentationHandler).Methods("GET")
    router.HandleFunc("/help", handlers.GetHtmlDocumentationHandler).Methods("GET")

    // User e
    router.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")
    router.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
    router.HandleFunc("/users/{id}", handlers.GetUserByIdHandler).Methods("GET")

    // Coffee 
    router.HandleFunc("/coffees", handlers.GetCoffeesHandler).Methods("GET")
    router.HandleFunc("/coffees/{id}", handlers.GetCoffeeHandler).Methods("GET")
    router.Handle("/coffees", middleware.AuthMiddleware(http.HandlerFunc(handlers.CreateCoffeeHandler))).Methods("POST")
    router.Handle("/coffees/{id}", middleware.AuthMiddleware(http.HandlerFunc(handlers.UpdateCoffeeHandler))).Methods("PUT")
    router.Handle("/coffees/{id}", middleware.AuthMiddleware(http.HandlerFunc(handlers.DeleteCoffeeHandler))).Methods("DELETE")

    // Coffee Shop 
    router.HandleFunc("/shops", handlers.GetCoffeeShopsHandler).Methods("GET")
    router.HandleFunc("/shops/{id}", handlers.GetCoffeeShopHandler).Methods("GET")
    router.Handle("/shops", middleware.AuthMiddleware(http.HandlerFunc(handlers.CreateCoffeeShopHandler))).Methods("POST")
    router.Handle("/shops/{id}", middleware.AuthMiddleware(http.HandlerFunc(handlers.UpdateCoffeeShopHandler))).Methods("PUT")
    router.Handle("/shops/{id}", middleware.AuthMiddleware(middleware.AdminMiddleware(http.HandlerFunc(handlers.DeleteCoffeeShopHandler)))).Methods("DELETE")

    // Roasteries 
    router.HandleFunc("/roasteries", handlers.GetRoasteriesHandler).Methods("GET")
    router.HandleFunc("/roasteries/{id}", handlers.GetRoasteryHandler).Methods("GET")
    router.Handle("/roasteries", middleware.AuthMiddleware(http.HandlerFunc(handlers.CreateRoasteryHandler))).Methods("POST")
    router.Handle("/roasteries/{id}", middleware.AuthMiddleware(http.HandlerFunc(handlers.UpdateRoasteryHandler))).Methods("PUT")
    router.Handle("/roasteries/{id}", middleware.AuthMiddleware(middleware.AdminMiddleware(http.HandlerFunc(handlers.DeleteRoasteryHandler)))).Methods("DELETE")

    // Reviews
    router.HandleFunc("/reviews", handlers.GetReviewsHandler).Methods("GET")
    router.Handle("/reviews", middleware.AuthMiddleware(http.HandlerFunc(handlers.CreateReviewHandler))).Methods("POST")
    router.Handle("/reviews/{id}", middleware.AuthMiddleware(http.HandlerFunc(handlers.UpdateReviewHandler))).Methods("PUT")
    router.Handle("/reviews/{id}", middleware.AuthMiddleware(http.HandlerFunc(handlers.DeleteReviewHandler))).Methods("DELETE")

    // Stats
    router.HandleFunc("/stats", handlers.GetStatsHandler).Methods("GET")

    router.Use(middleware.CORSMiddleware)

    port := ":40331"
    server := &http.Server{
        Addr:    port,
        Handler: router,
    }

    fmt.Printf("Running on port %s\n", port)
    fmt.Println("Dokumentacja: http://localhost" + port + "/help")
    if err := server.ListenAndServe(); err != nil {
        log.Fatal("Błąd podczas uruchamiania:", err)
    }
}