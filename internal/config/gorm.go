package config

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	usernameEnvName              = "DB_USERNAME"
	passwordEnvName              = "DB_PASSWORD"
	hostEnvName                  = "DB_HOST"
	portEnvName                  = "DB_PORT"
	databaseEnvName              = "DB_NAME"
	idleConnectionEnvName        = "DB_POOL_IDLE"
	maxConnectionEnvName         = "DB_POOL_MAX"
	maxLifeTimeConnectionEnvName = "DB_POOL_MAX_LIFE_TIME"
)

func NewDatabase() *gorm.DB {
	username := os.Getenv(usernameEnvName)
	password := os.Getenv(passwordEnvName)
	host := os.Getenv(hostEnvName)
	port := os.Getenv(portEnvName)
	database := os.Getenv(databaseEnvName)
	idleConnection, err := strconv.Atoi(os.Getenv(idleConnectionEnvName))
	if err != nil {
		panic("Failed to parse database idle connection")
	}
	maxConnection, err := strconv.Atoi(os.Getenv(maxConnectionEnvName))
	if err != nil {
		panic("Failed to parse database max connection")
	}
	maxLifeTime, err := strconv.Atoi(os.Getenv(maxLifeTimeConnectionEnvName))
	if err != nil {
		panic("Failed to parse database max life time")
	}
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable TimeZone=UTC", host, port, username, database, password)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return nil
	}

	connection, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return nil
	}

	connection.SetMaxIdleConns(idleConnection)
	connection.SetMaxOpenConns(maxConnection)
	connection.SetConnMaxLifetime(time.Duration(maxLifeTime) * time.Second)

	return db
}
