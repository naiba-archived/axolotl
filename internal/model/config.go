package model

type Config struct {
	Code struct {
		Endpoint      string `json:"endpoint,omitempty"`
		Authorization string `json:"authorization,omitempty"`
	} `json:"code,omitempty"`
	GitHub struct {
		ClientID     string `json:"client_id,omitempty"`
		ClientSecret string `json:"client_secret,omitempty"`
	} `json:"git_hub,omitempty"`
	Site struct {
		Name string `json:"name,omitempty"`
		Desc string `json:"desc,omitempty"`
	} `json:"site,omitempty"`
}
