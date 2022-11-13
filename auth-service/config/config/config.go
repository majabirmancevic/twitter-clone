package config

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

type Config struct {
	Port       string
	AuthDBHost string
	AuthDBPort string
}

func NewConfig() *Config {
	return &Config{
		Port:       os.Getenv("AUTH_SERVICE_PORT"),
		AuthDBHost: os.Getenv("AUTH_DB_HOST"),
		AuthDBPort: os.Getenv("AUTH_DB_PORT"),
	}
}

func GetClient(host, port string) (*mongo.Client, error) {
	uri := fmt.Sprintf("mongodb://%s:%s/", host, port)
	options := options.Client().ApplyURI(uri)
	return mongo.Connect(context.TODO(), options)
}
