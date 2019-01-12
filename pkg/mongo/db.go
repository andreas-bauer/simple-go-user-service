package mongo

import (
	"context"
	"fmt"
	"github.com/andreas-bauer/simple-go-user-service/pkg/model"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/readpref"
	"github.com/sirupsen/logrus"
)

type Datastore interface {
	FindAll() ([]*model.User, error)
	FindByEmail(email string) (*model.User, error)
	Delete(email string) error
	Save(email string) error
}

type DB struct {
	collection *mongo.Collection
}

type Connection struct {
	host     string
	database string
	username string
	password string
}

func (con *Connection) GetUri() (uri string) {
	uri = fmt.Sprintf(`mongodb://%s:%s@%s/%s`,
		con.username,
		con.password,
		con.host,
		con.database,
	)
	return
}

var DefaultConnection = &Connection{
	host:     "mongodb:27017",
	username: "admin",
	password: "admin",
	database: "admin",
}

func (db *DB) Connect(con Connection) {
	logrus.Info("Connect to mongo DB ", con.host)

	uri := con.GetUri()

	client, err := mongo.Connect(context.TODO(), uri)
	if err != nil {
		logrus.WithError(err).Error("Unable to establish DB connection to ", uri)
	}

	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		logrus.WithError(err).Error("Unable to ping database")
	}

	db.collection = client.Database("userservice").Collection("user")
}

func (db *DB) FindAll() (results []model.User, err error) {
	ctx := context.TODO()
	cur, err := db.collection.Find(ctx, nil)

	if err != nil {
		logrus.WithError(err).Error()
	}

	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var user model.User
		err = cur.Decode(&user)
		results = append(results, user)
	}
	return
}

func (db *DB) FindByEmail(email string) (result *model.User, err error) {
	filter := bson.D{{"email", email}}
	err = db.collection.FindOne(context.TODO(), filter).Decode(&result)
	return
}

func (db *DB) ContainsUserWithEmail(email string) bool {
	user, _ := db.FindByEmail(email)
	return user != nil
}

func (db *DB) Delete(email string) (err error) {
	filter := bson.D{{"email", email}}
	_, err = db.collection.DeleteOne(context.TODO(), filter)
	return
}

func (db *DB) Save(user model.User) (err error) {
	_, err = db.collection.InsertOne(context.TODO(), user)
	return
}

func (db *DB) CreateDefaultAdminUserIfNotExist() {
	defaultAdminUser := model.User{
		Name:     "GeneratedAdmin",
		Email:    "admin@adminland.de",
		Password: "$2a$10$zZeGbbtwwUC8gfpgAVx/v.hX95qMf/dIWOpgwiyZcPTcTxvNnBYN.",
		Role:     model.Enum.ADMIN}

	adminExists := db.ContainsUserWithEmail(defaultAdminUser.Email)
	if adminExists {
		return
	}

	logrus.Info("Create default admin user because it doesn't exist yet.")
	err := db.Save(defaultAdminUser)

	if err != nil {
		logrus.WithError(err).Error("Unable to create default admin user.")
	}

}
