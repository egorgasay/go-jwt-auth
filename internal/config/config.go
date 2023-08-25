package config

type Storage struct {
	DatabaseDSN string `json:"dsn"`
}

type JWT struct {
	Key        string `json:"key"`
	AccessTTL  string `json:"access_ttl"`
	RefreshTTL string `json:"refresh_ttl"`
}
