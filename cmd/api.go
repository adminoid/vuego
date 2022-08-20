package main

import (
	"context"
	"fmt"
	config "github.com/adminoid/vuego/internal/config"
	"github.com/adminoid/vuego/internal/entities/user"
	"github.com/adminoid/vuego/pkg/clients/postgresql"
	"github.com/adminoid/vuego/pkg/logging"
	"github.com/adminoid/vuego/pkg/project_path"
)

func main() {
	fmt.Println("----------------")
	fmt.Println(project_path.Root)
	fmt.Println("----------------")
	cfg := config.NewConfig()
	fmt.Println("cfg:")
	fmt.Println(cfg)

	postgresqlClient, err := postgresql.NewClient(context.TODO(), 3, cfg)
	if err != nil {
		fmt.Errorf("%v", err)
	}

	fmt.Println(postgresqlClient)

	logger := logging.GetLogger()

	repository := user.NewRepository(postgresqlClient, logger)
	users, err1 := repository.FindAll(context.TODO())
	if err1 != nil {
		logger.Fatal(err1)
	}
	fmt.Println(users)
}
