package app

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

var (
	Setting Config
	Main    *Application
)

type Application struct {
	Dlp *DlpClass
}

type Config struct {
	ProjectID      string `yaml:"project_id" env:"PROJECT_ID"`
	ServiceAccount string `yaml:"service_account" env:"SERVICE_ACCOUNT" env-default:"config/service_account.json"`
}

func Init() {
	// read yaml config
	err := cleanenv.ReadConfig("config/config.yaml", &Setting)
	if err != nil {
		log.Fatalf("could not read the config file. %v", err)
	}

	// set GCP client
	err = os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", Setting.ServiceAccount)
	if err != nil {
		log.Fatalf("could not set environment variable. %v", err)
	}

	Main = &Application{}
	Main.Dlp = &DlpClass{}
	log.Println("Setting application")
}
