package app

import (
	"context"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

var (
	Setting Config
)

type Config struct {
	projectID string `yaml:"project_id" env:"PROJECT_ID"`
}

func Init() {
	ctx := context.Background()
	// read yaml config
	err := cleanenv.ReadConfig("config/config.yaml", &Setting)
	if err != nil {
		log.Fatalf("could not read the config file. %v", err)
	}

	dlp := &Dlp{}
	dlp.Connect(ctx)
}
