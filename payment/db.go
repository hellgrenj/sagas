package main

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBAccess struct {
	conn *mongo.Client
}

func NewDBAccess(logger *Logger) *DBAccess {
	conn := TryConnectToMongo(1, logger)
	return &DBAccess{conn: conn}
}
func TryConnectToMongo(connectionAttempt int, logger *Logger) *mongo.Client {
	const uri = "mongodb://mongoadmin:mongopwd@paymentdb:27017/?maxPoolSize=20&w=majority"
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		logger.info.Printf("Unable to connect to database: %v\n", err)
		if connectionAttempt < 5 {
			connectionAttempt++
			logger.info.Printf("Trying again in 4 seconds attempt %v of 5\n", connectionAttempt)
			time.Sleep(4 * time.Second)
			return TryConnectToMongo(connectionAttempt, logger)
		}
		os.Exit(1)
	}
	logger.info.Println("Successfully connected to mongo")
	return client
}
