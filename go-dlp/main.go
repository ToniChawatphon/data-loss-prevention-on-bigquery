package main

import (
	"github.com/ToniChawatphon/data-loss-prevention-on-bigquery/app"
)

func main() {
	var input string
	input = "BTC"

	app.Init()
	app.Main.Dlp.Scan(input, app.Setting.ProjectID)
}
