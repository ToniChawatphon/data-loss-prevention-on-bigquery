package main

import (
	"log"

	"github.com/ToniChawatphon/data-loss-prevention-on-bigquery/app"
)

func main() {
	app.Init()
	log.Println(app.Setting.projectID)
}
