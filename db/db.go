package db

import (
	"fmt"
	"go-htmx-test/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(host string, user string, password string, dbname string, port int, timezone string) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%s", host, user, password, dbname, port, timezone)
	d, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	// AutoMigrate
	if err := d.AutoMigrate(&models.Item{}); err != nil {
		panic("failed to auto migrate")
	}

	DB = d

	fmt.Println("Connected to database")
}
