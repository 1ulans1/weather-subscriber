package db

import (
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"subscription-service/config"
	"time"
)

func Gorm(log *zap.Logger) *gorm.DB {
	conf := &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	db, err := gorm.Open(postgres.Open(config.Conf.DB.URL), conf)
	if err != nil {
		log.Fatal("Failed to connect to the database", zap.Error(err))
	}
	return db
}
