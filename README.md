# Coffee API

REST API zbudowane w Go do zarządzania danymi dotyczącymi kawy, w tym danymi kaw, palarni, kawiarni, recenzji oraz uwierzytelniania użytkowników przy użyciu JWT. API wykorzystuje bazę danych PostgreSQL i integruje się z zewnętrznymi usługami geokodowania.

## Funkcje

- **Uwierzytelnianie:**  
  - Rejestracja i logowanie z wykorzystaniem haszowania haseł (bcrypt)  
  - Autoryzacja oparta na tokenach JWT

- **Kawy:**  
  - Operacje CRUD dla kaw, z możliwością filtrowania po kraju, procesie, nutach smakowych  
  - Nuty smakowe zapisywane jako ciągi tekstowe rozdzielane przecinkami

- **Palarnie kawy:**  
  - Zarządzanie informacjami o palarniach z automatycznym geokodowaniem  
  - Aktualizacja współrzędnych przy zmianie adresu  
  - Zabezpieczenie przed usunięciem palarni, które mają powiązane kawy

- **Kawiarnie:**  
  - Endpointy CRUD dla kawiarni z danymi lokalizacyjnymi  
  - Automatyczne obliczanie średniej ocen na podstawie recenzji

- **Recenzje:**  
  - Tworzenie, aktualizacja i usuwanie recenzji z walidacją ocen  
  - Zarządzanie uprawnieniami do usuwania recenzji (właściciel lub admin)

- **Geolokalizacja:**  
  - Integracja z API Nominatim i Photon do konwersji adresów na współrzędne  
  - Mechanizm awaryjny w przypadku problemów z geokodowaniem

## Przegląd Endpointów

- **Użytkownicy:**  
  - `POST /register` – Rejestracja nowego użytkownika  
  - `POST /login` – Logowanie i otrzymanie tokena JWT

- **Kawy:**  
  - `GET /coffees` – Pobieranie wszystkich kaw  
  - `GET /coffees/{id}` – Pobieranie kawy po ID  
  - `POST /coffees` – Dodawanie nowej kawy (wymaga uwierzytelnienia)  
  - `PUT /coffees/{id}` – Aktualizacja kawy (wymaga uwierzytelnienia)  
  - `DELETE /coffees/{id}` – Usuwanie kawy (wymaga uwierzytelnienia)

- **Palarnie:**  
  - `GET /roasteries` – Pobieranie wszystkich palarni  
  - `GET /roasteries/{id}` – Pobieranie palarni po ID  
  - `POST /roasteries` – Dodawanie nowej palarni (wymaga uwierzytelnienia)  
  - `PUT /roasteries/{id}` – Aktualizacja palarni (wymaga uwierzytelnienia)  
  - `DELETE /roasteries/{id}` – Usuwanie palarni (tylko admin)

- **Kawiarnie:**  
  - `GET /shops` – Pobieranie wszystkich kawiarni  
  - `GET /shops/{id}` – Pobieranie kawiarni po ID  
  - `POST /shops` – Dodawanie nowej kawiarni (wymaga uwierzytelnienia)  
  - `PUT /shops/{id}` – Aktualizacja kawiarni (wymaga uwierzytelnienia)  
  - `DELETE /shops/{id}` – Usuwanie kawiarni (tylko admin)

- **Recenzje:**  
  - `GET /reviews` – Pobieranie recenzji z opcjonalnym filtrowaniem  
  - `POST /reviews` – Dodawanie recenzji (wymaga uwierzytelnienia)  
  - `PUT /reviews/{id}` – Aktualizacja recenzji (wymaga uwierzytelnienia)  
  - `DELETE /reviews/{id}` – Usuwanie recenzji (właściciel lub admin)

## Narzędzia i Zależności

- Go
- PostgreSQL  
- Gorilla Mux
- JWT
- godotenv
- lib/pq