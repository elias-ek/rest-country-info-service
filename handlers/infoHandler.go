package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"assignment-1/models"
)

// Handler for the /info endpoint, provides general country information

func InfoHandler(w http.ResponseWriter, r *http.Request) {

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 || parts[4] == "" {
		http.Error(w, "Please specify a country to get information about, using its two-letter ISO code (no, us, gb, etc.). \n You may also use ?limit=<number> to limit the amount of cities returned.", http.StatusBadRequest)
		return
	}

	countryCode := parts[4]

	// Get optional limit query parameter
	limitStr := r.URL.Query().Get("limit")
	limit := 0
	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 0 {
			http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
			return
		}
	}

	// Get country information from REST Countries API
	restURL := fmt.Sprintf("http://129.241.150.113:8080/v3.1/alpha/%s", countryCode)
	resp, err := http.Get(restURL)
	if err != nil {
		http.Error(w, "Could not fetch data from REST Countries API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "REST Countries API returned an error", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Could not read response from REST Countries API", http.StatusInternalServerError)
		return
	}

	//Decode the response
	var countries []models.RestCountry
	if err := json.Unmarshal(body, &countries); err != nil || len(countries) == 0 {
		http.Error(w, "Could not decode REST Countries API response", http.StatusInternalServerError)
		return
	}

	country := countries[0]

	// Send a post requset to the CountriesNow API to get the cities

	cityURL := "http://129.241.150.113:3500/api/v0.1/countries/cities"

	cityReq, err := json.Marshal(map[string]string{"country": country.Name.Common})
	if err != nil {
		http.Error(w, "Error forming request for Cities", http.StatusInternalServerError)
		return
	}

	cityRes, err := http.Post(cityURL, "application/json", bytes.NewReader(cityReq))
	if err != nil {
		http.Error(w, "Error fetching data from CountriesNow API", http.StatusInternalServerError)
		return
	}
	defer cityRes.Body.Close()

	cityBody, err := io.ReadAll(cityRes.Body)
	if err != nil {
		http.Error(w, "Error reading response from CountriesNow API", http.StatusInternalServerError)
		return
	}

	var citiesResponse models.CitiesResponse
	if err := json.Unmarshal(cityBody, &citiesResponse); err != nil {
		http.Error(w, "Error decoding Cities response", http.StatusInternalServerError)
		return
	}
	cities := citiesResponse.Data

	// Sort cities in ascending alphabetical order
	sort.Strings(cities)
	if limit > 0 && limit < len(cities) {
		cities = cities[:limit]
	}

	// Build response
	infoRes := models.InfoResponse{

		Name:       country.Name.Common,
		Continents: country.Continents,
		Population: country.Population,
		Languages:  country.Languages,
		Borders:    country.Borders,
		Flag:       country.Flags.PNG,
		Capital:    country.Capital[0],
		Cities:     cities,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(infoRes)
}
