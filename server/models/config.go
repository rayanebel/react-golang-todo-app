package models

type Config struct {
	DatabaseEngine        string `envconfig:"DB_ENGINE"`
	MongoInstanceURI      string `envconfig:"MONGO_URI"`
	MongoInstanceUsername string `envconfig:"MONGO_USERNAME"`
	MongoInstancePassword string `envconfig:"MONGO_PASSWORD"`
}
