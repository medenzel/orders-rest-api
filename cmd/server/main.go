package main

import (
	"github.com/medenzel/orders-rest-api/internal/database"
	"github.com/medenzel/orders-rest-api/internal/order"
	transportHTTP "github.com/medenzel/orders-rest-api/internal/transport/http"
	"github.com/medenzel/orders-rest-api/migrations"
	log "github.com/sirupsen/logrus"
)

func Run() error {
	log.SetFormatter(&log.JSONFormatter{})
	log.Info("Setting app")

	store, err := database.NewDatabase()
	if err != nil {
		log.Error("failed to setup connection to database!")
		return err
	}
	err = store.Migrate(migrations.EmbedMigrations, ".")
	if err != nil {
		log.Error("failed to migrate database!")
		return err
	}

	orderService := order.NewService(store)
	handler := transportHTTP.NewHandler(orderService)

	if err := handler.Serve(); err != nil {
		log.Error("Failed to serve our app!")
		return err
	}

	return nil

}

func main() {
	if err := Run(); err != nil {
		log.Error(err)
		log.Fatal("Error starting app!")
	}
}
