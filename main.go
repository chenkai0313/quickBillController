package main

import (
	"flag"

	"quickBillController/app"
	"quickBillController/config"
	"quickBillController/internal"
	"quickBillController/models"
)

var debug bool

func init() {
	flag.BoolVar(&debug, "debug", false, "Is it a development environment?")
}

func main() {
	flag.Parse()
	config.LoadConfig(debug)
	app.InitLogger()
	app.InitDatabase()

	models.Migrate()

	internal.RunApi()
}
