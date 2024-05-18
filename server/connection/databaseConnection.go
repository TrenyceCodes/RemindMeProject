package connection

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoDatabaseConnection() (*mongo.Client, string, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://priceboom811:kxmpgzEmeEqZigU1@cluster0.xm9rk9a.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0").SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return client, "", err
	}
	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		return client, "", err
	}
	return client, "Pinged your deployment. You successfully connected to MongoDB!", err
}
