package models

// GeoServe contains the list of affected cities included in the GeoJSON alert provided by the USGS
type GeoServe struct {
	Cities []City `json:"cities"`
	Timezone struct {
		UTCOffset *int     `json:"utcOffset"`
		Latitude  *float64 `json:"latitude"`
		UtcTime   *string  `json:"utcTime"`
		Time      *string  `json:"time"`
		ShortName *string  `json:"shortName"`
		Longitude *float64 `json:"longitude"`
		LongName  *string  `json:"longName"`
	} `json:"timezone"`
	Region struct {
		Country *string `json:"country"`
		State   *string `json:"state"`
	} `json:"region"`
	Fe struct {
		Number      *int    `json:"number"`
		MediumName  *string `json:"mediumName"`
		HDS         *string `json:"hds"`
		SpanishName *string `json:"spanishName"`
		ShortName   *string `json:"shortName"`
		LongName    *string `json:"longName"`
	} `json:"fe"`
}

// City is used in GeoServe to represent a singular city element
type City struct {
	Distance         *int     `json:"distance"`
	Longitude        *float64 `json:"longitude"`
	Latitude         *float64 `json:"latitude"`
	Name             *string  `json:"name"`
	Direction        *string  `json:"direction"`
	Population       *int     `json:"population"`
	OpenStreetMapURL string   `json:"open_street_map"`
	AppleMapsURL     string   `json:"apple_maps"`
}
