package main

import (
	"context"
	config "github.com/adminoid/vuego/internal/config"
	"github.com/adminoid/vuego/internal/entities/user"
	"github.com/adminoid/vuego/pkg/clients/postgresql"
	"github.com/adminoid/vuego/pkg/logging"
	"github.com/adminoid/vuego/pkg/project_path"
	"github.com/julienschmidt/httprouter"
	"net"
	"net/http"
	"time"
)

func main() {

	logger := logging.GetLogger()
	logger.Println(project_path.Root)

	cfg := config.NewConfig()
	logger.Printf("root path: %s", project_path.Root)
	logger.Printf("config: %v", cfg)

	postgresqlClient, err := postgresql.NewClient(context.TODO(), 3, cfg)
	if err != nil {
		logger.Errorf("err: %v", err)
	}
	logger.Printf("postgresqlClient: %v", postgresqlClient)

	repository := user.NewRepository(postgresqlClient, logger)

	logger.Info("register user handler")
	userHandler := user.NewHandler(repository, logger)
	router := httprouter.New()
	userHandler.Register(router)

	start(router, cfg)
}

func start(router *httprouter.Router, cfg config.Config) {
	logger := logging.GetLogger()
	logger.Info("start application")

	var listener net.Listener
	var listenErr error

	logger.Info("listen tcp")
	listener, listenErr = net.Listen("tcp", cfg.BindAddr)
	logger.Infof("server is listening %s", cfg.BindAddr)

	if listenErr != nil {
		logger.Fatal(listenErr)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Fatal(server.Serve(listener))
}
