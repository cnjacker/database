package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	MongoURL      = "mongodb://mongo:mongo@127.0.0.1:27017"
	MongoDatabase = "app"
	MongoTimeout  = 1 * time.Second
)

type Mongo struct {
	Conn *mongo.Client
	Ctx  context.Context
}

func NewMongo() *Mongo {
	return NewMongoWithURL(MongoURL)
}

func NewMongoWithURL(url string) *Mongo {
	ctx := context.Background()

	if v, err := mongo.Connect(
		ctx,
		options.Client().SetConnectTimeout(MongoTimeout).ApplyURI(url),
	); err == nil {
		return &Mongo{
			Conn: v,
			Ctx:  ctx,
		}
	}

	return nil
}

func (m Mongo) Database() *mongo.Database {
	return m.Conn.Database(MongoDatabase)
}

func (m Mongo) Collection(collection string) *mongo.Collection {
	if v := m.Database(); v != nil {
		return v.Collection(collection)
	}

	return nil
}

func (m Mongo) Close() {
	m.Conn.Disconnect(m.Ctx)
}
