package config

import "os"

type Config struct {
	Port string
	URI  string
}

func NewConfig() *Config {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb+srv://cluster0.fhmwctr.mongodb.net/?authSource=%24external&authMechanism=MONGODB-X509&retryWrites=true&w=majority&appName=Cluster0&tlsCertificateKeyFile=X509-cert-211311170312867980.pem"
	}
	return &Config{Port: port, URI: uri}
}
