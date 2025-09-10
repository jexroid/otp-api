package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"

	"github.com/jexroid/gopi/internal/router"
)

// @title Gopi API
// @version 1.0
// @description This is a sample server for Gopi API.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8000
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		ForceColors: true,
	})

	var r = chi.NewRouter()
	router.Handler(r)

	log.Info("go server is running on http://localhost:8000")

	err = http.ListenAndServe(":8000", r)
	if err != nil {
		log.Error(err)
	}
}
