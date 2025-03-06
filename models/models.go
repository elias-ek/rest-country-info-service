package models

//response sent from the info endpoint
type InfoResponse struct {
	Name       string            `json:"name"`
	Continents []string          `json:"continents"`
	Population int               `json:"population"`
	Languages  map[string]string `json:"languages"`
	Borders    []string          `json:"borders"`
	Flag       string            `json:"flag"`
	Capital    string            `json:"capital"`
	Cities     []string          `json:"cities"`
}

//individual popluation counts gathered from countries now api
type PopulationCount struct {
	Year  int `json:"year"`
	Value int `json:"value"`
}

//whole list of population counts
type PopulationValue struct {
	Data struct {
		PopulationCounts []PopulationCount `json:"populationCounts"`
	} `json:"data"`
}

//response sent from population endpoint
type PopulationResponse struct {
	Mean   float64           `json:"mean"`
	Values []PopulationCount `json:"values"`
}

//response sent from status endpoint
type StatusResponse struct {
	CountriesNowAPI  string  `json:"countriesnowapi"`
	RestCountriesAPI string  `json:"restcountriesapi"`
	Version          string  `json:"version"`
	Uptime           float64 `json:"uptime"`
}

//name of country, as used by population handler to find correct population
type CountryName struct {
	Data struct {
		Name string `json:"name"`
	} `json:"data"`
}

//country gathered from rest countries api
type RestCountry struct {
	Name struct {
		Common string `json:"common"`
	} `json:"name"`
	Continents []string          `json:"continents"`
	Population int               `json:"population"`
	Languages  map[string]string `json:"languages"`
	Borders    []string          `json:"borders"`
	Flags      struct {
		PNG string `json:"png"`
	} `json:"flags"`
	Capital []string `json:"capital"`
}

//cities gathered from countries now
type CitiesResponse struct {
	Data []string `json:"data"`
}
