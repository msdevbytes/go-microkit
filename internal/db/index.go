package db

import (
	"fmt"

	"github.com/msdevbytes/go-microkit/internal/config"
	"github.com/msdevbytes/go-microkit/pkg/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Conn *gorm.DB

func OpenConnection() error {
	cfg := config.DBConfig()

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=UTC&time_zone='%%2B00:00'",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	dbConn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	Conn = dbConn

	// 🔐 Enforce UTC at session level
	sqlDB, err := dbConn.DB()
	if err != nil {
		return fmt.Errorf("error retrieving sql.DB: %w", err)
	}

	if _, err := sqlDB.Exec("SET time_zone = '+00:00'"); err != nil {
		return fmt.Errorf("failed to set session time_zone: %w", err)
	}

	logger.Info("Connected to the database")
	// Migrate the database schema
	autoMigrate()

	return nil
}

func CloseConnect() error {
	sqlDB, err := Conn.DB()
	if err != nil {
		return err
	}

	if err := sqlDB.Close(); err != nil {
		return err
	}
	logger.Info("Database connection closed")

	return nil
}
