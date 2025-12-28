package main

import (
	"web_service/internal/app/ds"
	"web_service/internal/app/dsn"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(
		&ds.Acid{},
		&ds.Carbonate{},
		&ds.CarbonateAcid{},
		&ds.Users{},
	)
	if err != nil {
		panic("cant migrate db")
	}
}
