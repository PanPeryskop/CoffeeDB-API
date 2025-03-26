# Coffee API

A RESTful API built in Go for managing coffee data, including coffees, roasteries, coffee shops, reviews, and user authentication using JWT. The API uses PostgreSQL as the database and integrates with external geocoding services (Nominatim & Photon) to obtain coordinates from addresses.

## Features

- **User Authentication:**  
  - Register and login endpoints with password hashing via bcrypt  
  - JWT-based authentication (see [`handlers.LoginHandler`](services/handlers/users.go))

- **Coffees:**  
  - CRUD operations for coffee records, including querying by country, process, flavour notes, etc.  
  - Flavour notes are stored as comma-separated strings (see [`handlers.coffees.go`](services/handlers/coffees.go))

- **Roasteries:**  
  - Manage roastery information with auto-geocoding to determine coordinates  
  - Re-geocode updated addresses when necessary  
  - Deletion prevented if associated coffees exist (see [`handlers.roasteries.go`](services/handlers/roasteries.go))

- **Coffee Shops:**  
  - CRUD endpoints for coffee shop records with geolocation data  
  - Average rating calculation based on reviews (see [`handlers.coffeeShops.go`](services/handlers/cofeeShops.go))

- **Reviews:**  
  - Endpoints to create, update, and delete reviews with rating validations  
  - Role-based review deletion allowed for admins (see [`handlers.reviews.go`](services/handlers/reviews.go))

- **Geolocation:**  
  - Integration with Nominatim and Photon APIs for address-to-coordinate conversion  
  - Fallback mechanism for geocoding failures (see [`geocoding.GetCoordinates`](services/geocoding/geocoding.go))


## Endpoints Overview

- **User Endpoints:**  
  - `POST /register` – Register a new user  
  - `POST /login` – Log in and obtain a JWT token

- **Coffee Endpoints:**  
  - `GET /coffees` – Retrieve all coffees  
  - `GET /coffees/{id}` – Retrieve a coffee by id  
  - `POST /coffees` – Create a new coffee (requires authentication)  
  - `PUT /coffees/{id}` – Update a coffee (requires authentication)  
  - `DELETE /coffees/{id}` – Delete a coffee (requires authentication)

- **Roastery Endpoints:**  
  - `GET /roasteries` – Retrieve all roasteries  
  - `GET /roasteries/{id}` – Retrieve a roastery by id  
  - `POST /roasteries` – Create a new roastery (requires authentication)  
  - `PUT /roasteries/{id}` – Update a roastery (requires authentication)  
  - `DELETE /roasteries/{id}` – Delete a roastery (admin only)

- **Coffee Shop Endpoints:**  
  - `GET /shops` – Retrieve all coffee shops  
  - `GET /shops/{id}` – Retrieve a coffee shop by id  
  - `POST /shops` – Create a new coffee shop (requires authentication)  
  - `PUT /shops/{id}` – Update a coffee shop (requires authentication)  
  - `DELETE /shops/{id}` – Delete a coffee shop (admin only)

- **Review Endpoints:**  
  - `GET /reviews` – Retrieve reviews, optionally filtered by coffee, roastery, shop, or user  
  - `POST /reviews` – Create a review (requires authentication)  
  - `PUT /reviews/{id}` – Update a review (requires authentication)  
  - `DELETE /reviews/{id}` – Delete a review (owner or admin only)

## Tools & Dependencies

- Go 1.24.1  
- PostgreSQL  
- [Gorilla Mux](https://github.com/gorilla/mux)  
- [golang-jwt/jwt](https://github.com/golang-jwt/jwt)  
- [joho/godotenv](https://github.com/joho/godotenv)  
- [lib/pq](https://github.com/lib/pq)  

## To Do

See the [todo](todo) file for upcoming enhancements such as improved error handling, logging, caching geocode results, and unit/integration tests.

## License

Specify your license here.
