package database

import (
	"fmt"
	"log"
	"sportsbook-backend/internal/config"
	"sportsbook-backend/internal/types"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

var SPORTS []*types.SportItem
var COUNTRIES []*types.CountryItem
var TOURNAMENTS []*types.TournamentItem

func InitPostgresDB(cfg *config.Config) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUsername, cfg.DBPassword, cfg.DBName)

	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database %v:\n", err)
	}
	DB = db
}
