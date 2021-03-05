package repository

import (
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type MongoClient struct {
	client          *mongo.Client
	movieCollection *mongo.Collection
}

func (m *MongoClient) initialize() error {
	var err error
	ctx := context.Background()
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	m.client, err = mongo.Connect(ctx, clientOptions)

	if err != nil {
		return errors.Wrap(err, "Could not connect to MongoDB")
	}

	if err = m.client.Ping(ctx, nil); err != nil {
		return errors.Wrap(err, "Failed to ping MongoDB")
	}

	log.Println("Connected to MongoDB!")

	m.movieCollection = m.client.Database("golang-training").Collection("movies")
	return nil
}

func (m *MongoClient) insertMovieBatch(batch []*Movie, errorChan chan error) {
	var ui []interface{}
	for _, movie := range batch {
		if movie != nil {
			ui = append(ui, movie)
		}
	}
	_, err := m.movieCollection.InsertMany(context.TODO(), ui)
	if err != nil {
		log.Printf("%+v", batch)
	}
	errorChan <- err
}
