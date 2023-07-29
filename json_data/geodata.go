package json_data

type GeoData struct {
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Address      string  `json:"address"`
	XInfo        string  `json:"xinfo"`
	UserDataSign string  `json:"userDataSign"`
	Destinations []int64 `json:"destinations"`
	Locale       string  `json:"locale"`
	Shard        int     `json:"shard"`
	Currency     string  `json:"currency"`
	IP           string  `json:"ip"`
	Dt           int     `json:"dt"`
}
