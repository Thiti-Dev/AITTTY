package database

import (
	"context"
	"fmt"

	"github.com/Thiti-Dev/AITTTY/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connect -> is a (fn) that uses for connecting to the db instance
func Connect() (*mongo.Database,error){
	clientOptions := options.Client()
	clientOptions.ApplyURI(config.LoadConfig("MONGO_HOST"))
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	var ctx = context.Background()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to MongoDB!")

	return client.Database(config.LoadConfig("MONGO_DB_NAME")), nil
}