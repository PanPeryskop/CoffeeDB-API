package geocoding

import (
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
    "net/url"
    "strconv"
    "strings"
    "time"
)

type NominatimResponse []struct {
    Lat string `json:"lat"`
    Lon string `json:"lon"`
}

type PhotonResponse struct {
    Features []struct {
        Geometry struct {
            Coordinates []float64 `json:"coordinates"` // [lon, lat]
        } `json:"geometry"`
    } `json:"features"`
}

func translatePolishChars(s string) string {
    polishToAscii := map[rune]string{
        'ą': "a", 'ć': "c", 'ę': "e", 'ł': "l", 'ń': "n", 'ó': "o", 'ś': "s", 'ź': "z", 'ż': "z",
        'Ą': "A", 'Ć': "C", 'Ę': "E", 'Ł': "L", 'Ń': "N", 'Ó': "O", 'Ś': "S", 'Ź': "Z", 'Ż': "Z",
    }
    var result strings.Builder
    for _, r := range s {
        if replacement, ok := polishToAscii[r]; ok {
            result.WriteString(replacement)
        } else {
            result.WriteRune(r)
        }
    }
    return result.String()
}

func removeStreetPrefixes(address string) string {
    prefixes := []string{"ul. ", "ul.", "straße ", "str. ", "ул. ", "улица "}
    for _, prefix := range prefixes {
        address = strings.ReplaceAll(address, prefix, "")
    }
    return address
}

func getNominatimCoordinates(address string) (float64, float64, error) {
    encodedAddress := url.QueryEscape(address)
    requestURL := fmt.Sprintf("https://nominatim.openstreetmap.org/search?format=json&q=%s", encodedAddress)
    client := &http.Client{Timeout: 10 * time.Second}
    req, err := http.NewRequest("GET", requestURL, nil)
    if err != nil {
        return 0, 0, err
    }
    req.Header.Set("User-Agent", "CoffeeApiGeocoder/1.0")
    resp, err := client.Do(req)
    if err != nil {
        return 0, 0, err
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        return 0, 0, errors.New("geocoding API request failed")
    }
    var results NominatimResponse
    if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
        return 0, 0, err
    }
    if len(results) == 0 {
        return 0, 0, errors.New("no results found for the address")
    }
    lat, err := strconv.ParseFloat(results[0].Lat, 64)
    if err != nil {
        return 0, 0, err
    }
    lon, err := strconv.ParseFloat(results[0].Lon, 64)
    if err != nil {
        return 0, 0, err
    }
    return lat, lon, nil
}

func getPhotonCoordinates(address string) (float64, float64, error) {
    encodedAddress := url.QueryEscape(address)
    requestURL := fmt.Sprintf("https://photon.komoot.io/api/?q=%s&limit=1", encodedAddress)
    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Get(requestURL)
    if err != nil {
        return 0, 0, err
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        return 0, 0, fmt.Errorf("photon API request failed with status: %d", resp.StatusCode)
    }
    var pr PhotonResponse
    if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
        return 0, 0, err
    }
    if len(pr.Features) == 0 || len(pr.Features[0].Geometry.Coordinates) < 2 {
        return 0, 0, errors.New("photon: no results found for the address")
    }
    coords := pr.Features[0].Geometry.Coordinates
    return coords[1], coords[0], nil
}

func GetCoordinates(address string) (float64, float64, error) {
    address = translatePolishChars(address)
    fmt.Println("Address:", address)
    lat, lon, err := getNominatimCoordinates(address)
    if err == nil {
        return lat, lon, nil
    }
    if strings.Contains(err.Error(), "no results") {
        modifiedAddress := removeStreetPrefixes(address)
        fmt.Println("Retrying with modified address:", modifiedAddress)
        lat, lon, err = getNominatimCoordinates(modifiedAddress)
        if err == nil {
            return lat, lon, nil
        }
    }
    fmt.Println("Using Photon fallback for address:", address)
    lat, lon, err = getPhotonCoordinates(address)
    if err == nil {
        return lat, lon, nil
    }
    return 0, 0, fmt.Errorf("all geocoding attempts failed: %v", err)
}