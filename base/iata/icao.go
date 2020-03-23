package iata

//        "icao": "00AL",
//        "iata": "",
//        "name": "Epps Airpark",
//        "city": "Harvest",
//        "state": "Alabama",
//        "country": "US",
//        "elevation": 820,
//        "lat": 34.8647994995,
//        "lon": -86.7703018188,
//        "tz": "America\/Chicago"

type Icao struct {
	Icao      string  `json:"icao"`
	Iata      string  `json:"iata"`
	Name      string  `json:"name"`
	City      string  `json:"city"`
	State     string  `json:"state"`
	Country   string  `json:"country"`
	Elevation int     `json:"elevation"`
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	Tz        string  `json:"tz"`
}
