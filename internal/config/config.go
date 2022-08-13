package config

import (
	"fmt"
	"github.com/adminoid/vuego/pkg/project_path"
	"github.com/joho/godotenv"
	"github.com/mitchellh/mapstructure"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	DbUser   string
	DbPwd    string
	DbName   string
	DbHost   string
	DbPort   string
	LogLevel string
	BindAddr string
}

const EnvPath = ".env"

// NewConfig config constructor
func NewConfig() Config {
	env := getEnv()
	config := &Config{}
	if err := mapstructure.Decode(env, &config); err != nil {
		log.Fatalln(err)
	}
	return *config
}

func getEnv() map[string]string {
	var myEnv map[string]string

	fullPath := filepath.Join(project_path.Root, EnvPath)
	fmt.Printf("path is %s", fullPath)

	myEnv, err := godotenv.Read(os.ExpandEnv(fullPath))
	if err != nil {
		panic("reading env error")
	}

	return myEnv
}
