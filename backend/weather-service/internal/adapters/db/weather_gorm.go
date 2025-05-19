package db

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
	"weather-service/internal/core/domain"
	"weather-service/internal/core/ports"
)

type weatherRecord struct {
	gorm.Model
	Location    string `gorm:"uniqueIndex"`
	Temperature float64
	Condition   string
	UpdatedAt   int64
}

type weatherRepo struct {
	db *gorm.DB
}

func NewWeatherRepo(db *gorm.DB) ports.WeatherRepo {
	repo := &weatherRepo{db: db}
	repo.init()
	return repo
}

func (r *weatherRepo) init() {
	if err := r.db.AutoMigrate(&weatherRecord{}); err != nil {
		panic(err)
	}
}

func (r *weatherRepo) Save(w *domain.Weather) error { //todo test
	record := weatherRecord{
		Location:    w.Location,
		Temperature: w.Temperature,
		Condition:   w.Condition,
		UpdatedAt:   w.UpdatedAt.Unix(),
	}
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "location"}},
		UpdateAll: true,
	}).Create(&record).Error
}

func (r *weatherRepo) Get(location string) (*domain.Weather, error) {
	var record weatherRecord
	if err := r.db.Where("location = ?", location).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &domain.Weather{
		Location:    record.Location,
		Temperature: record.Temperature,
		Condition:   record.Condition,
		UpdatedAt:   time.Unix(record.UpdatedAt, 0),
	}, nil
}
