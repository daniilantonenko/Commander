package database

import (
	"commander/config"
	"commander/internal/models"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Return new Postgresql db instance
func NewDB(c *config.Config) *gorm.DB {
	var err error
	dbuser := viper.GetString("DB_USER")
	dbpassword := viper.GetString("DB_PASS")
	dbname := viper.GetString("DB_NAME")
	dbhost := viper.GetString("DB_HOST")

	dsn := "host=" + dbhost + " user=" + dbuser + " password=" + dbpassword + " dbname=" + dbname + " sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&models.Command{})

	log.Println("Connected to database")

	return db
}
