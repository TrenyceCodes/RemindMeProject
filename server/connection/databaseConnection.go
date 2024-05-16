package connection

import (
	"context"
	"fmt"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoDatabaseConnection() (*mongo.Client, string, error) {
	envFile, _ := godotenv.Read(".env")

	ctx := context.Background()

	// Set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(envFile["DATABASE_CONNECTION"]).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, "", err
	}

	// Ping to confirm a successful connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, "", err
	}

	// Check if the MongoDB connection is nil
	if client == nil {
		fmt.Println("MongoDB connection is nil. Aborting insertion.")
		return nil, "", err
	}

	// Check if the connection is still alive before proceeding
	if err := client.Ping(context.Background(), nil); err != nil {
		fmt.Println("Error pinging MongoDB:", err)
		return nil, "", err
	}

	return client, "Pinged your deployment. You successfully connected to MongoDB!", nil
}
