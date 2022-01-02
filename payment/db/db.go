package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/hellgrenj/sagas/payment/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Logger interface {
	Info(msg string)
	Error(msg string)
}
type dba struct {
	conn *mongo.Client
}

func NewDBAccess(logger Logger) *dba {
	conn := tryConnectToMongo(1, logger)
	db := &dba{conn: conn}
	coll := db.conn.Database("payment").Collection("processedmessages")
	_, err := coll.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "messageid", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		panic(err)
	}
	return db
}
func tryConnectToMongo(connectionAttempt int, logger Logger) *mongo.Client {
	const uri = "mongodb://mongoadmin:mongopwd@paymentdb:27017/?maxPoolSize=20&w=majority"
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		logger.Info(fmt.Sprintf("Unable to connect to database: %v\n", err))
		if connectionAttempt < 5 {
			connectionAttempt++
			logger.Info(fmt.Sprintf("Trying again in 4 seconds attempt %v of 5\n", connectionAttempt))
			time.Sleep(4 * time.Second)
			return tryConnectToMongo(connectionAttempt, logger)
		}
		os.Exit(1)
	}
	logger.Info("Successfully connected to mongo")
	return client
}
func (db *dba) InsertPayment(orderPayment models.OrderPayment) error {
	_, err := db.conn.Database("payment").Collection("payments").InsertOne(context.TODO(), orderPayment)
	return err
}
func (db *dba) TryMarkMessageAsProcessed(messageId string) (bool, error) {
	alreadyProcessed := false
	_, err := db.conn.Database("payment").Collection("processedmessages").InsertOne(context.Background(),
		bson.D{
			{Key: "messageid", Value: messageId},
		})
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			alreadyProcessed = true
			return alreadyProcessed, nil
		} else {
			return alreadyProcessed, err
		}
	}
	return alreadyProcessed, err
}
