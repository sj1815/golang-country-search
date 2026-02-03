package model

// Country represents the country information returned by the API.
type Country struct {
	Name       string `json:"name"`
	Capital    string `json:"capital"`
	Currency   string `json:"currency"`
	Population int    `json:"population"`
}

// RESTCountryResponse represents the response structure from the REST Countries API.
type RESTCountryResponse struct {
	Name       CountryName             `json:"name"`
	Capital    []string                `json:"capital"`
	Currencies map[string]CurrencyInfo `json:"currencies"`
	Population int                     `json:"population"`
}

// CountryName represents the name structure in the REST Countries API response.
type CountryName struct {
	Common   string `json:"common"`
	Official string `json:"official"`
}

// CurrencyInfo represents the currency information in the REST Countries API response.
type CurrencyInfo struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

// ErrorResponse represents a standard error response structure.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}
