package database

import (
	"fmt"
	"os"

	"github.com/jexroid/gopi/pkg/models"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func Init() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_DATABASE"),
		os.Getenv("DB_PORT"))

	var db, err = gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{
			SkipDefaultTransaction: true,
			PrepareStmt:            true,
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		},
	)

	if err != nil {
		logrus.Panic("Database connection error: ", err)
	}

	// Migrate all your models
	migratingAuthError := db.AutoMigrate(
		&models.User{},
	)

	if migratingAuthError != nil {
		logrus.Panic("Migration error: ", migratingAuthError)
	}

	logrus.Info("Database connected and migrated successfully")
	return db
}
