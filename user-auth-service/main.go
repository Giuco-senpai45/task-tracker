package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"ua-service/controller"
	"ua-service/db"
	"ua-service/http/routes"
	"ua-service/service"
	"ua-service/utils/log"
)

func main() {
	log.Instantiate()

	dbUserName := os.Getenv("DB_USER")
	dbUserPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	dbName := os.Getenv("DB_NAME")

	dbConfig := db.DatabaseConfig{
		Host:     dbHost,
		Port:     dbPort,
		Username: dbUserName,
		Password: dbUserPassword,
		DBName:   dbName,
	}

	dbAdapter, err := db.NewAdapter(dbConfig)
	if err != nil {
		log.Error("Error connecting to database: %v", err)
	}
	defer db.CloseConnection(dbAdapter)

	as := service.NewAuthService(dbAdapter)
	ac := controller.NewAuthController(as)

	router := http.NewServeMux()
	routes.RegisterAuthRoutes(router, ac)

	appPort := fmt.Sprintf(":%s", os.Getenv("APP_PORT"))
	if appPort == "" {
		appPort = ":8081"
	}

	server := http.Server{
		Addr:    appPort,
		Handler: router,
	}

	log.Info("Starting server on port %s", appPort)
	server.ListenAndServe()
}
