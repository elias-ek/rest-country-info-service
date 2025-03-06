package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"assignment-1/models"
)

// Handler for the population endpoint, provides population per year and calculated mean

func PopHandler(w http.ResponseWriter, r *http.Request) {

	parts := strings.Split(r.URL.Path, "/")

	//check for country code
	if len(parts) < 5 || parts[4] == "" {
		http.Error(w, "Please provide a two-letter country code (no, gb, us etc.) and optionally ?limit=<startYear>-<endYear> to limit the years to get population values from.", http.StatusBadRequest)
		return
	}

	countryCode := parts[4]

	// check for optional limit parameter and determine year range
	limitParam := r.URL.Query().Get("limit")

	var startYear, endYear int
	limitProvided := false

	//check validity of start- and end year
	if limitParam != "" {
		limitProvided = true
		years := strings.Split(limitParam, "-")
		if len(years) != 2 {
			http.Error(w, "Invalid limit format, expected startYear-endYear", http.StatusBadRequest)
			return
		}

		var err error
		startYear, err = strconv.Atoi(years[0])
		if err != nil {
			http.Error(w, "Invalid start year", http.StatusBadRequest)
			return
		}

		endYear, err = strconv.Atoi(years[1])

		if err != nil {
			http.Error(w, "Invalid end year", http.StatusBadRequest)
			return
		}

		if startYear > endYear {
			http.Error(w, "start year cannot be greater than end year", http.StatusBadRequest)
			return
		}

	}

	//CountriesNow population endpoint requires iso3 or full country name, so get name from the codes endpoint (iso3 is not listed here so have to use name)
	nameURL := fmt.Sprintf("http://129.241.150.113:3500/api/v0.1/countries/codes/q?iso2=%s", countryCode)
	resp, err := http.Get(nameURL)
	if err != nil {
		http.Error(w, "Could not get data from CountriesNow API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "countries now api returned an error", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "could not read response", http.StatusInternalServerError)
		return
	}

	var countryName models.CountryName
	if err := json.Unmarshal(body, &countryName); err != nil {
		http.Error(w, "could not decode response", http.StatusInternalServerError)
		return
	}

	//use country name to get population data, replacing all spaces with %20 for url formatting
	popURL := fmt.Sprintf("http://129.241.150.113:3500/api/v0.1/countries/population/q?country=%s", strings.ReplaceAll(countryName.Data.Name, " ", "%20"))

	popRes, err := http.Get(popURL)
	if err != nil {
		http.Error(w, "could not get population data", http.StatusInternalServerError)
	}
	defer popRes.Body.Close()

	if popRes.StatusCode != http.StatusOK {
		http.Error(w, "countries now api returned an error", popRes.StatusCode)
		return
	}

	popBody, err := io.ReadAll(popRes.Body)
	if err != nil {
		http.Error(w, "could not read response", http.StatusInternalServerError)
		return
	}

	//decode the response
	var popInfo models.PopulationValue
	if err := json.Unmarshal(popBody, &popInfo); err != nil {
		http.Error(w, "error decoding population response", http.StatusInternalServerError)
		return
	}

	var filteredCounts []models.PopulationCount
	var total int
	var yearsCount int

	//for all counts received, if limited check year, add to filteredCounts array
	for _, count := range popInfo.Data.PopulationCounts {

		if limitProvided {
			if count.Year < startYear || count.Year > endYear {
				continue
			}
		}

		filteredCounts = append(filteredCounts, models.PopulationCount{
			Year:  count.Year,
			Value: count.Value,
		})

		//add value of each population count, increment counter for how many years
		total += count.Value
		yearsCount++
	}

	//calculate mean
	var mean float64
	if yearsCount > 0 {
		mean = float64(total) / float64(yearsCount)
	}

	//place mean and values in final response
	finalResp := models.PopulationResponse{
		Mean:   mean,
		Values: filteredCounts,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(finalResp)
}
