/**
 * Copyright (c) 2019 Andreas Bauer
 *
 * SPDX-License-Identifier: MIT
 */
package mongo

import (
	"context"
	"fmt"

	"github.com/andreas-bauer/simple-go-user-service/pkg/user"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DB struct {
	collection *mongo.Collection
}

type Connection struct {
	Host     string
	Database string
	Username string
	Password string
}

func (con *Connection) GetUri() (uri string) {
	uri = fmt.Sprintf(`mongodb://%s:%s@%s/%s`,
		con.Username,
		con.Password,
		con.Host,
		con.Database,
	)
	return
}

var DefaultConnection = &Connection{
	Host:     "localhost:27017",
	Username: "admin",
	Password: "admin",
	Database: "admin",
}

func (db *DB) Connect(con Connection) {
	logrus.Info("Connect to MongoDB with URI ", con.GetUri())

	uri := con.GetUri()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		logrus.WithError(err).Error("Unable to establish DB connection to ", uri)
	}

	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		logrus.WithError(err).Error("Unable to ping Database")
	}

	db.collection = client.Database("userservice").Collection("user")
}

func (db *DB) FindAll() (results []*user.User, err error) {
	ctx := context.TODO()
	cur, err := db.collection.Find(ctx, nil)

	if err != nil {
		logrus.WithError(err).Error()
	}

	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var user user.User
		err = cur.Decode(&user)
		results = append(results, &user)
	}
	return
}

func (db *DB) FindByEmail(email string) (result *user.User, err error) {
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

func (db *DB) Save(user *user.User) (err error) {
	_, err = db.collection.InsertOne(context.TODO(), user)
	return
}

func (db *DB) CreateDefaultAdminUserIfNotExist() {
	defaultAdminUser := user.User{
		Name:     "GeneratedAdmin",
		Email:    "admin@adminland.de",
		Password: "$2a$10$zZeGbbtwwUC8gfpgAVx/v.hX95qMf/dIWOpgwiyZcPTcTxvNnBYN.",
		Role:     user.ADMIN}

	adminExists := db.ContainsUserWithEmail(defaultAdminUser.Email)
	if adminExists {
		return
	}

	logrus.Info("Create default admin user because it doesn't exist yet.")
	err := db.Save(&defaultAdminUser)

	if err != nil {
		logrus.WithError(err).Error("Unable to create default admin user.")
	}

}
