package database

import (
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	database, err := gorm.Open(mysql.Open(os.Getenv("DB_DSN")), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(err)
	} else {
		log.Default().Println("Database initialized")
	}
	return database
}
