package db

import (
	"fmt"
	"time"
	"ua-service/models"
	"ua-service/utils/log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	DBName   string
}

func NewAdapter(config DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", config.Username, config.Password, config.Host, config.Port, config.DBName)

	var db *gorm.DB
	var err error

	for i := 0; i < 5; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		} else {
			str := fmt.Sprintf("Failed to connect to database %d, Retrying ...", i)
			log.Warn(str)
			time.Sleep(5 * time.Second)
		}
	}
	if err != nil {
		log.Error("Failed to connect to database")
		return nil, err
	}

	err = testDB(db)
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.User{})
	log.Info("Database tasks migrated successfully")

	return db, nil
}

func testDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	err = sqlDB.Ping()
	if err != nil {
		return err
	}

	log.Info("Pinged database successfully !")
	return nil
}

func CloseConnection(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Error("Failed to close database connection")
	}

	sqlDB.Close()
	log.Info("Database connection closed successfully")
}
