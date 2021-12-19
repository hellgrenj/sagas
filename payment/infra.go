package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type infra struct {
	db     *DBAccess
	logger Logger
}

func NewInfraHandler(db *DBAccess, logger Logger) *infra {
	// init setup
	i := &infra{db: db, logger: logger}
	coll := i.db.conn.Database("payment").Collection("processedmessages")
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
	return i
}
func (i *infra) TryMarkMessageAsProcessed(messageId string) (bool, error) {
	alreadyProcessed := false
	processedmessages := i.db.conn.Database("payment").Collection("processedmessages")
	_, err := processedmessages.InsertOne(context.Background(), bson.D{
		{Key: "messageid", Value: messageId},
	})
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			alreadyProcessed = true
			return alreadyProcessed, nil
		} else {
			return alreadyProcessed, err
		}
	} else {
		i.logger.Info().Printf("Inserted %s into processedmessages collection!\n", messageId)
		return alreadyProcessed, nil
	}
}
