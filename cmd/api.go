package main

import (
	"fmt"
	config "github.com/adminoid/vuego/internal/config"
	"github.com/adminoid/vuego/pkg/project_path"
)

func main() {
	fmt.Println("----------------")
	fmt.Println(project_path.Root)
	fmt.Println("----------------")
	cfg := config.NewConfig()
	fmt.Println("cfg:")
	fmt.Println(cfg)
}
