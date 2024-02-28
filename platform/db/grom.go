package db

import (
	"fmt"
	"log"
	"os"

	"github.com/r3tr056/ecolens_api/app/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var PostgresDB *gorm.DB

var postgresHost = os.Getenv("POSTGRES_HOST")
var postgresPort = os.Getenv("POSTGRES_PORT")
var postgresUser = os.Getenv("POSTGRES_USERNAME")
var postgresPass = os.Getenv("POSTGRES_PASS")
var dbName = os.Getenv("POSTGRES_DB")
var sslMode = os.Getenv("POSTGRES_SSLMODE")

func OpenPostgresConnection() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", postgresHost, postgresUser, postgresPass, dbName, postgresPort, sslMode)

	var err error
	PostgresDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	// Automigrate
	err = PostgresDB.AutoMigrate(&models.Brand{}, &models.ProductImage{}, &models.LCAMetrics{}, &models.EnvironmentalProductDeclaration{}, &models.Report{}, &models.Product{}, &models.MarketPlaceProduct{})
	if err != nil {
		log.Fatalf("Failed to auto migrate : %v", err)
	}

	if err := PostgresDB.SetupJoinTable(&models.EnvironmentalProductDeclaration{}, "LCAMetrics", &models.LCAMetrics{}); err != nil {
		log.Fatalf("Failed to SetupJoinTable: %v", err)
	}

	return nil
}

func DatabaseCheck() bool {
	err := PostgresDB.Exec("SELECT 1").Error
	return err == nil
}

func CustomMigrate() {
	// Create GIN Indexes for the Searchable fields
	PostgresDB.Migrator().CreateIndex(&models.Product{}, "Name")
	PostgresDB.Migrator().CreateIndex(&models.Product{}, "Description")

	PostgresDB.Migrator().CreateIndex(&models.MarketPlaceProduct{}, "Name")
	PostgresDB.Migrator().CreateIndex(&models.MarketPlaceProduct{}, "Description")

	PostgresDB.Migrator().CreateIndex(&models.Report{}, "Name")
	PostgresDB.Migrator().CreateIndex(&models.Report{}, "Summary")
}
