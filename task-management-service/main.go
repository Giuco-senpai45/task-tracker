package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
	"tm-service/controller"
	"tm-service/db"
	"tm-service/http/routes"
	"tm-service/producer"
	"tm-service/service"
	"tm-service/utils/log"
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

	var kafkaProducer *producer.MessageProducer
	maxAttempts := 10
	for i := 0; i < maxAttempts; i++ {
		kafkaProducer, err = producer.NewKafkaProducer()
		if err != nil {
			log.Error("Error creating kafka producer: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	if kafkaProducer == nil {
		log.Fatal("Failed to create kafka producer after %d attempts", maxAttempts)
	}
	ts, err := service.NewTaskService(dbAdapter, kafkaProducer)
	tc := controller.NewTaskController(ts)

	// ts.CheckDeadlines()

	router := http.NewServeMux()
	routes.RegisterTaskRoutes(router, tc)

	appPort := fmt.Sprintf(":%s", os.Getenv("APP_PORT"))
	if appPort == "" {
		appPort = ":8080"
	}

	server := http.Server{
		Addr:    appPort,
		Handler: router,
	}

	log.Info("Starting server on port %s", appPort)
	server.ListenAndServe()
}
