package storage

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/golang-grpc-proxy/config"
	"github.com/golang-grpc-proxy/src/entities"
	slogGorm "github.com/orandin/slog-gorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ErrNotFound = errors.New("not found")
)

type Storage struct {
	db *gorm.DB
}

func New(config *config.Config, logger *slog.Logger) (Storage, error) {
	gormLogger := slogGorm.New(
		slogGorm.WithHandler(logger.Handler()),
		slogGorm.SetLogLevel(slogGorm.ErrorLogType, slog.LevelError),
		slogGorm.SetLogLevel(slogGorm.SlowQueryLogType, slog.LevelError),
		slogGorm.SetLogLevel(slogGorm.DefaultLogType, slog.LevelError),
	)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", config.Database.Host, config.Database.User, config.Database.Password, config.Database.Name, config.Database.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})

	if (err != nil) {
		return Storage{}, err
	}

	if config.Env == "local" {
		err := db.AutoMigrate(&entities.User{})
		if err != nil {
			panic("user migration error: " + err.Error())
		}
	}

	return Storage{db: db}, nil
}

func (s *Storage) IsUserUnique(email, phone string) (bool, error) {
	var user entities.User

	result := s.db.Where("email = ? or phone = ?", email, phone).Take(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return true, nil
	}

	if result.Error == nil {
		return false, nil
	}

	return false, result.Error
}

func (s *Storage) CreateUser(user entities.User) (entities.User, error) {
	result := s.db.Create(&user)
	return user, result.Error
}
func (s *Storage) GetUserByEmail(email string) (entities.User, error) {
	var user entities.User

	result := s.db.Where("email=?", email).Take(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return user, ErrNotFound
	}
	
	return user, result.Error
}