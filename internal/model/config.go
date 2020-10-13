package model

type Config struct {
	Code struct {
		Endpoint      string
		Authorization string
	}
	GitHub struct {
		ClientID     string
		ClientSecret string
	}
}
