package database

import (
	"fmt"
	"strconv"
	"time"

	"backend/pkg/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

var DB *gorm.DB

func ConnectDB() (*gorm.DB, error) {
	cfg := config.InitConfig()

	var dsn string
	if len(cfg.DB.URL) == 0 {
		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.DB.HOST, cfg.DB.PORT, cfg.DB.USER, cfg.DB.PASS, cfg.DB.NAME, cfg.DB.SSLMODE,
		)
	} else {
		dsn = cfg.DB.URL
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.Use(tracing.NewPlugin()); err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	maxIdle, _ := strconv.Atoi(cfg.DB.MAXIDLE)
	maxOpenConn, _ := strconv.Atoi(cfg.DB.MAXOPENCONN)
	maxLifeTime, _ := strconv.Atoi(cfg.DB.MAXLIFETIME)

	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetMaxOpenConns(maxOpenConn)
	sqlDB.SetConnMaxLifetime(time.Duration(maxLifeTime) * time.Hour)

	DB = db
	return db, nil
}
