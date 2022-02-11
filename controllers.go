package main

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (a *App) getDbConnection() (*mongo.Database, *mongo.Client, error) {
	uri, exists := os.LookupEnv("APP_DB_URI")
	if !exists {
		a.API.Logger.Fatal("missing database connection string")
	}
	var err error
	var db *mongo.Database
	var client *mongo.Client
	for i := 1; i <= 10; i++ {
		db, client, err = connectOrFail(uri)
		if err == nil {
			break
		}
		a.API.Logger.Printf("connecting to database failed, iteration: %d, err: %s", i, err)
		time.Sleep(10 * time.Second)
	}
	return db, client, err
}

func examinations(db *mongo.Database) ([]Examination, error) {
	var exams []Examination
	queryOptions := options.Find()
	queryOptions.SetSort(bson.D{{"exam", 1}})
	queryOptions.SetLimit(300)
	colName, founded := os.LookupEnv("APP_DB_COL")
	if founded {
		cursor, err := db.Collection(colName).Find(context.TODO(), bson.D{{}}, queryOptions)
		if err != nil {
			return []Examination{}, err
		}
		defer func(cursor *mongo.Cursor, ctx context.Context) {
			err := cursor.Close(ctx)
			if err != nil {
				log.Printf("closing cursor failed: %s", err)
			}
		}(cursor, context.TODO())
		for cursor.Next(context.TODO()) {
			var e Examination
			if err = cursor.Decode(&e); err != nil {
				log.Printf("decoding examination failed: %s", err.Error())
			}
			if err != nil {
				return []Examination{}, err
			}
			exams = append(exams, e)
		}
		return exams, nil
	} else {
		return []Examination{}, errors.New("collection name missing, define APP_DB_COL variable")
	}
}
