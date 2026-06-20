package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"

	c "go-darts-tgbot/internal/config"
)

type Database struct {
	Dsn string
	DB  *gorm.DB
}

func InitDatabase() *Database {
	return &Database{
		Dsn: c.GlobalConfig.Dsn,
	}
}

func Connect(d *Database) error {
	db, err := gorm.Open(postgres.Open(d.Dsn), &gorm.Config{})

	if err != nil {
		slog.Error("db connection error", slog.Any("err", err))
		return err
	}

	sqlDB, err := db.DB()

	if err != nil {
		slog.Error("failed to get DB instance", slog.Any("err", err))
		return err
	}

	err = sqlDB.Ping()
	if err != nil {
		slog.Error("db connection check failed", slog.Any("err", err))
		return err
	}

	slog.Info("db connection established")

	d.DB = db
	return nil
}
