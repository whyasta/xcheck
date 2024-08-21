package config

import (
	"bigmind/xcheck-be/config/dblogger"
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type AppConfig struct {
	DATABASE_HOST        string
	DATABASE_PORT        string
	DATABASE_USER        string
	DATABASE_PASSWORD    string
	DATABASE_NAME        string
	DATABASE_QUERY_DEBUG bool
	AUTH_JWT_SECRET      string
	AUTH_JWT_EXPIRE      int
	AUTH_JWT_ISSUER      string
	SERVER_ADDRESS       string
	SERVER_PORT          string
	APP_ENV              string
	REDIS_HOST           string
	REDIS_PORT           string
	REDIS_QUEUE          string
	CLOUD_BASE_URL       string
}

var config *viper.Viper

func Init(env string) {
	var err error
	config = viper.New()
	config.SetConfigType("env")
	config.SetConfigName(".env")
	config.AddConfigPath(".")

	err = config.ReadInConfig()
	if err != nil {
		log.Fatal("error on parsing configuration file")
	}
}

func ConnectToDB() (*gorm.DB, error) {
	var err error
	configEnv := GetConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", configEnv.GetString("DATABASE_USER"), configEnv.GetString("DATABASE_PASSWORD"), configEnv.GetString("DATABASE_HOST"), configEnv.GetString("DATABASE_PORT"), configEnv.GetString("DATABASE_NAME"))

	gormConfig := &gorm.Config{}
	if configEnv.GetBool("DATABASE_QUERY_DEBUG") {
		dbLogger := dblogger.New(zap.L())
		dbLogger.SetAsDefault()
		dbLogger.LogMode(logger.Info)
		gormConfig = &gorm.Config{
			// Logger: logger.Default.LogMode(logger.Info),
			// Logger: dbLogger,
			Logger: dbLogger,
		}
	}
	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	// db = db.Debug()
	if err != nil {
		log.Fatal("Error connecting to database. Error: ", err)
		return nil, err
	}

	// cache
	// cachesPlugin := &caches.Caches{Conf: &caches.Config{
	// 	Easer: true,
	// }}
	// _ = db.Use(cachesPlugin)

	sqlDB, err := db.DB()

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(20)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	// sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	return db, err
}

func GetDsn() string {
	configEnv := GetConfig()
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		configEnv.GetString("DATABASE_USER"),
		configEnv.GetString("DATABASE_PASSWORD"),
		configEnv.GetString("DATABASE_HOST"),
		configEnv.GetString("DATABASE_PORT"),
		configEnv.GetString("DATABASE_NAME"),
	)
}

func GetConfig() *viper.Viper {
	return config
}

func GetAppConfig() *AppConfig {
	return &AppConfig{
		DATABASE_HOST:        GetConfig().GetString("DATABASE_HOST"),
		DATABASE_PORT:        GetConfig().GetString("DATABASE_PORT"),
		DATABASE_USER:        GetConfig().GetString("DATABASE_USER"),
		DATABASE_PASSWORD:    GetConfig().GetString("DATABASE_PASSWORD"),
		DATABASE_NAME:        GetConfig().GetString("DATABASE_NAME"),
		DATABASE_QUERY_DEBUG: GetConfig().GetBool("DATABASE_QUERY_DEBUG"),
		AUTH_JWT_SECRET:      GetConfig().GetString("AUTH_JWT_SECRET"),
		AUTH_JWT_EXPIRE:      GetConfig().GetInt("AUTH_JWT_EXPIRE"),
		AUTH_JWT_ISSUER:      GetConfig().GetString("AUTH_JWT_ISSUER"),
		SERVER_ADDRESS:       GetConfig().GetString("SERVER_ADDRESS"),
		SERVER_PORT:          GetConfig().GetString("SERVER_PORT"),
		APP_ENV:              GetConfig().GetString("APP_ENV"),
		REDIS_HOST:           GetConfig().GetString("REDIS_HOST"),
		REDIS_PORT:           GetConfig().GetString("REDIS_PORT"),
		REDIS_QUEUE:          GetConfig().GetString("REDIS_QUEUE"),
		CLOUD_BASE_URL:       GetConfig().GetString("CLOUD_BASE_URL"),
	}
}
