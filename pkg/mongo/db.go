package mongo

import (
	"context"
	"fmt"
	"github.com/andreas-bauer/simple-go-user-service/pkg/model"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/readpref"
	"github.com/sirupsen/logrus"
	"time"
)

type DB struct {
	collection *mongo.Collection
	ctx        context.Context
}

type MongoConnection struct {
	host     string
	database string
	username string
	password string
}

func (con *MongoConnection) GetUri() (uri string) {
	uri = fmt.Sprintf(`mongodb://%s:%s@%s/%s`,
		con.username,
		con.password,
		con.host,
		con.database,
	)
	return
}

var DefaultConnection = &MongoConnection{
	host:     "localhost:27017",
	username: "admin",
	password: "admin",
	database: "admin",
}

func (db *DB) Connect(con MongoConnection) {
	logrus.Info("Connect to mongo DB ...")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	db.ctx = ctx

	uri := con.GetUri()

	fmt.Println(uri)
	client, err := mongo.Connect(ctx, uri)
	if err != nil {
		logrus.WithError(err).Error("Unable to establish DB connection to ", uri)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		logrus.WithError(err).Error("Unable to ping database")
	}

	db.collection = client.Database("userservice").Collection("user")
}

func (db *DB) FindAll() {
	cur, err := db.collection.Find(db.ctx, nil)

	if err != nil {
		logrus.WithError(err).Error()
	}
	defer cur.Close(db.ctx)
	for cur.Next(db.ctx) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			logrus.WithError(err).Error()
		}
		// do something with result....
		fmt.Println(result)
	}
	if err := cur.Err(); err != nil {
		logrus.WithError(err).Error()
	}
}

func (r *DB) FindByEmail(email string) {
	//TODO implement
}

func (r *DB) Delete(email string) {
	//TODO implement
}

func (r *DB) Save(user model.User) {
	//TODO implement
}
