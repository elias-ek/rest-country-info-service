package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"assignment-1/conf"
	"assignment-1/models"
	"assignment-1/util"
)

// handler for the status endpoint, checks known endpoints that should return status 200, also returns time since started
func StatHandler(w http.ResponseWriter, r *http.Request) {

	restURL := "http://129.241.150.113:8080/v3.1/all"
	restStatus := util.CheckAPIStatus(restURL)

	countriesNowURL := "http://129.241.150.113:3500/api/v0.1/countries/codes/q?iso2=no"
	countriesNowStatus := util.CheckAPIStatus(countriesNowURL)

	uptime := time.Since(conf.StartTime).Seconds()

	statusResp := models.StatusResponse{
		CountriesNowAPI:  countriesNowStatus,
		RestCountriesAPI: restStatus,
		Version:          "v1",
		Uptime:           uptime,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(statusResp)
}
