package db

import (
	"context"

	"github.com/athunlal/bookNowBooking-svc/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDataBase(cfg config.Config) (*mongo.Database, error) {
	// Create a MongoDB client
	clientOptions := options.Client().ApplyURI("mongodb://" + cfg.DBHost + ":" + cfg.DBPort)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Check if the connection was successful
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	// Set the database name
	db := client.Database(cfg.DBName)

	// Optionally, you can add any additional configurations or settings here if needed

	return db, nil
}
