package db

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var once sync.Once
var client *mongo.Client

func Connect() *mongo.Client {
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Fatal(err)
		}

		uri := os.Getenv("DB_URI")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")

		credentials := options.Credential{
			Username: user,
			Password: password,
		}
		clientOptions := options.Client().ApplyURI(uri).SetAuth(credentials)
		client, err = mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			log.Fatal(err)
		}

		err = client.Ping(context.TODO(), nil)
		if err != nil {
			log.Fatal(err)
		}

		indexQuery := []string{"email"}

		CreateIndex("users", indexQuery, options.Index().SetUnique(true))
	})

	return client
}

func CreateIndex(collectionName string, indexQuery []string, option *options.IndexOptions) {
	collection := Collection(collectionName)
	index := []mongo.IndexModel{}

	for _, val := range indexQuery {
		index = append(index, mongo.IndexModel{
			Keys: bson.D{
				{
					Key:   val,
					Value: 1,
				},
			},
			Options: option,
		})
	}

	collection.Indexes().CreateMany(
		context.Background(),
		index,
	)
}

func Collection(name string) *mongo.Collection {
	return client.Database("command").Collection(name)
}

func Disconnect() {
	client.Disconnect(context.TODO())
}
