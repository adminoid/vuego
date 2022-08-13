package main

import (
	"context"
	"fmt"
	config "github.com/adminoid/vuego/internal/config"
	"github.com/adminoid/vuego/pkg/clients/postgresql"
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

	//repository := author.NewRepository(postgreSQLClient, logger)
	//a := author2.Author{
	//	Name: "OK",
	//}
	//err = repository.Create(context.TODO(), &a)
	//if err != nil {
	//	logger.Fatal(err)
	//}

}
