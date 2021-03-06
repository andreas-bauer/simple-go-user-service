/**
 * Copyright (c) 2019 Andreas Bauer
 *
 * SPDX-License-Identifier: MIT
 */

package rest

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/andreas-bauer/simple-go-user-service/pkg/mongo"
	"github.com/andreas-bauer/simple-go-user-service/pkg/user"
	"github.com/sirupsen/logrus"
)

type Instance struct {
	db         *mongo.DB
	httpServer *http.Server
}

const (
	port = ":8080"
)

func (srv *Instance) Start() {
	// Startup DB
	srv.db = &mongo.DB{}
	srv.db.Connect(getDBConnectionSettings())
	srv.db.CreateDefaultAdminUserIfNotExist()

	// Startup HTTP
	logrus.Info("Http Server starting with address ", port)

	service := user.NewService(srv.db)
	srv.httpServer = &http.Server{Addr: port, Handler: Router(service)}
	err := srv.httpServer.ListenAndServe()

	if err != http.ErrServerClosed {
		logrus.WithError(err).Error("Http Server stopped unexpected")
		srv.Shutdown()
	} else {
		logrus.WithError(err).Info("Http Server stopped")
	}
}

func (srv *Instance) Shutdown() {
	if srv.httpServer != nil {
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		err := srv.httpServer.Shutdown(ctx)
		if err != nil {
			logrus.WithError(err).Error("Failed to shutdown http server gracefully")
		} else {
			srv.httpServer = nil
		}
	}
}

func getDBConnectionSettings() (con mongo.Connection) {
	con = *mongo.DefaultConnection

	host, hostExist := os.LookupEnv("MONGO_HOSTNAME")
	if hostExist {
		con.Host = host
	}

	username, usernameExist := os.LookupEnv("MONGO_USERNAME")
	if usernameExist {
		con.Username = username
	}

	pw, pwExist := os.LookupEnv("MONGO_PASSWORD")
	if pwExist {
		con.Password = pw
	}

	authDb, authDbExist := os.LookupEnv("MONGO_AUTH_DB")
	if authDbExist {
		con.Database = authDb
	}

	return
}
