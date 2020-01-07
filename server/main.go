package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
	"github.com/rayanebel/react-golang-todo-app/server/config"
	"github.com/rayanebel/react-golang-todo-app/server/middleware"
	"github.com/rayanebel/react-golang-todo-app/server/models"
	"github.com/rayanebel/react-golang-todo-app/server/router"
)

const (
	defaultDatabaseEngine = "badger"
)

func DatabaseEngineIsSupported(engine string) bool {
	supportedDatabase := []string{"badger", "mongo"}
	for _, item := range supportedDatabase {
		if item == engine {
			return true
		}
	}
	return false
}

func main() {
	r := router.Router()
	var cfg models.Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	if cfg.DatabaseEngine == "" {
		cfg.DatabaseEngine = defaultDatabaseEngine
	}

	if !DatabaseEngineIsSupported(cfg.DatabaseEngine) {
		log.Fatalf("Unsupported '%s' as database engine...", cfg.DatabaseEngine)
	}

	log.Println("Database Engine:", cfg.DatabaseEngine)
	if cfg.DatabaseEngine == "badger" {
		middleware.InitBadgerDB()
		defer middleware.Database.Close()
	} else if cfg.DatabaseEngine == "mongo" {
		connectionString := fmt.Sprintf("mongodb+srv://%s:%s@%s/%s?retryWrites=true&w=majority", cfg.MongoInstanceUsername, cfg.MongoInstancePassword, cfg.MongoInstanceURI, "todos")
		middleware.InitMongoDB(connectionString)
		defer middleware.MongoDatabase.Disconnect(context.Background())
	}

	// Save config into variable to be used by every subpackage
	config.SaveConfig(cfg)

	fmt.Println("Starting server on the port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
