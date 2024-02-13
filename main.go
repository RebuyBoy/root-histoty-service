package main

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"root-histoty-service/config"
	"root-histoty-service/db"
	"root-histoty-service/internal/controller"
	"root-histoty-service/internal/repository"
	"root-histoty-service/internal/service"
	"root-histoty-service/pkg/auth"
)

func main() {
	log := logrus.New()

	err := godotenv.Load()
	if err != nil {
		log.Error("Error loading .env file")
	}

	dbConfig := config.GetDBConfig()
	connection, err := db.ConnectPostgresDB(dbConfig)

	defer func(connection *sqlx.DB) {
		err := connection.Close()
		if err != nil {
			log.Error("Connection closed with error: ", err)
		}
	}(connection)

	userRepo := repository.NewUserRepo(connection, &dbConfig)
	err = RunMigrations(connection, dbConfig)
	if err != nil {
		log.Warning(err)
	}
	tokenManager, err := auth.NewTokenManager(config.GetServerConfig().SecretWord, 60, 180)

	serverConfig := config.GetServerConfig()
	playerService := service.NewPlayerService(userRepo, log, tokenManager)
	srv := controller.NewServer(&serverConfig, log, playerService)

	srv.RegisterRoutes()
	srv.StartRouter()
}

func RunMigrations(connection *sqlx.DB, c config.DBConfig) error {
	driver, err := pgx.WithInstance(connection.DB, &pgx.Config{})
	if err != nil {
		return fmt.Errorf("failed to get migration tool driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		c.Database,
		driver)
	if err != nil {
		return fmt.Errorf("failed to connect migration tool: %w", err)
	}

	err = m.Up()
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
