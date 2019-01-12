/**
 * Copyright (c) 2019 Andreas Bauer
 *
 * SPDX-License-Identifier: MIT
 */

package rest

import (
	"context"
	"github.com/andreas-bauer/simple-go-user-service/pkg/mongo"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
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
	srv.db.Connect(*mongo.DefaultConnection)
	srv.db.CreateDefaultAdminUserIfNotExist()

	// Startup HTTP
	logrus.Info("Http Server starting with address ", port)
	srv.httpServer = &http.Server{Addr: port, Handler: Router(srv)}
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
