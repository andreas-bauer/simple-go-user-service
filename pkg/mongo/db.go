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

type Datastore interface {
	FindAll() ([]*model.User, error)
	FindByEmail(email string) (*model.User, error)
	Delete(email string) error
	Save(email string) error
}

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

func (db *DB) FindAll() (results []model.User, err error) {
	cur, err := db.collection.Find(db.ctx, nil)

	if err != nil {
		logrus.WithError(err).Error()
	}

	defer cur.Close(db.ctx)
	for cur.Next(db.ctx) {
		var user model.User
		err = cur.Decode(&user)
		results = append(results, user)
	}
	return
}

func (db *DB) FindByEmail(email string) (result *model.User, err error) {
	filter := bson.D{{"email", email}}
	err = db.collection.FindOne(db.ctx, filter).Decode(&result)
	return
}

func (db *DB) ContainsUserWithEmail(email string) bool {
	user, _ := db.FindByEmail(email)
	return user != nil
}

func (db *DB) Delete(email string) (err error) {
	filter := bson.D{{"email", email}}
	_, err = db.collection.DeleteOne(db.ctx, filter)
	return
}

func (db *DB) Save(user model.User) (err error) {
	_, err = db.collection.InsertOne(db.ctx, user)
	return
}
