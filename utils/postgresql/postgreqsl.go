package postgresql

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func NewPostgreSQL(host string, port string, user string, password string, dbname string, logLevel string) (db *gorm.DB, err error) {
	// Формируем строку подключения
	uri := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host,     // Хост PostgreSQL
		port,     // Порт PostgreSQL
		user,     // Имя пользователя
		password, // Пароль
		dbname,   // Имя базы данных
	)

	var gormLogLevel logger.LogLevel
	switch logLevel {
	case "info":
		gormLogLevel = logger.Info
	case "warn":
		gormLogLevel = logger.Warn
	case "error":
		gormLogLevel = logger.Error
	case "silent":
		gormLogLevel = logger.Silent
	default:
		gormLogLevel = logger.Warn // Уровень по умолчанию
	}

	// Конфигурация логгера
	pgLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // Настройка вывода
		logger.Config{
			SlowThreshold:             time.Second,  // Порог медленных запросов
			LogLevel:                  gormLogLevel, // Уровень логирования
			IgnoreRecordNotFoundError: true,         // Игнорировать ошибку "record not found"
			Colorful:                  true,         // Цветной вывод
		},
	)

	db, err = gorm.Open(postgres.Open(uri), &gorm.Config{
		Logger: pgLogger, // Логируем запросы для отладки
	})
	if err != nil {
		log.Fatalf("Unable to connect to PostgreSQL with GORM: %v", err)
	}

	fmt.Println("PostgreSQL connected with GORM")

	// Проверяем соединение
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Unable to retrieve raw DB connection: %v", err)
	}
	err = sqlDB.Ping()
	if err != nil {
		log.Fatalf("Unable to ping PostgreSQL: %v", err)
	}
	return
}
