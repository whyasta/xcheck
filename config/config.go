package config

import (
	"bigmind/xcheck-be/config/dblogger"
	"fmt"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type AppConfig struct {
	DatabaseHost       string `mapstructure:"DATABASE_HOST"`
	DatabasePort       string `mapstructure:"DATABASE_PORT"`
	DatabaseUser       string `mapstructure:"DATABASE_USER"`
	DatabasePassword   string `mapstructure:"DATABASE_PASSWORD"`
	DatabaseName       string `mapstructure:"DATABASE_NAME"`
	DatabaseQueryDebug bool   `mapstructure:"DATABASE_QUERY_DEBUG"`
	AuthJwtSecret      string `mapstructure:"AUTH_JWT_SECRET"`
	AuthJwtExpire      int    `mapstructure:"AUTH_JWT_EXPIRE"`
	AuthJwtIssuer      string `mapstructure:"AUTH_JWT_ISSUER"`
	ServerAddress      string `mapstructure:"SERVER_ADDRESS"`
	ServerPort         string `mapstructure:"SERVER_PORT"`
	AppEnv             string `mapstructure:"APP_ENV"`
	RedisHost          string `mapstructure:"REDIS_HOST"`
	RedisPort          string `mapstructure:"REDIS_PORT"`
	RedisQueue         string `mapstructure:"REDIS_QUEUE"`
	CloudBaseURL       string `mapstructure:"CLOUD_BASE_URL"`
	MinioEndpoint      string `mapstructure:"MINIO_ENDPOINT"`
	MinioBucket        string `mapstructure:"MINIO_BUCKET"`
	MinioAccessKey     string `mapstructure:"MINIO_ACCESSKEY"`
	MinioSecretKey     string `mapstructure:"MINIO_SECRETKEY"`
}

var config *viper.Viper

func Init(env string) {
	var err error
	config = viper.New()
	log.Println("Init Config", env)

	if env != "local" {
		config.SetConfigType("env")
		config.SetConfigName(".env")
	} else {
		config.SetConfigType("env")
		config.SetConfigName(".env.local")
	}
	config.AddConfigPath(".")
	config.AutomaticEnv()
	err = config.ReadInConfig()
	if err != nil {
		log.Println("error on parsing configuration file")
		log.Println("Trying to load from env variable")
		// scaffoldLocal()
		var c AppConfig
		if e := parseEnv(&c); e != nil {
			log.Fatal("error on parsing configuration file or env variable", e)
			panic(e)
		}
	}
}

func scaffoldLocal() {
	os.Setenv("DATABASE_HOST", "101.50.2.153")
	os.Setenv("DATABASE_PORT", "3306")
	os.Setenv("DATABASE_USER", "user_xcheck")
	os.Setenv("DATABASE_PASSWORD", "xcheck@2024")
	os.Setenv("DATABASE_NAME", "xcheck")
	os.Setenv("DATABASE_QUERY_DEBUG", "false")
	os.Setenv("AUTH_JWT_SECRET", "bmok2024")
	os.Setenv("AUTH_JWT_EXPIRE", "60")
	os.Setenv("AUTH_JWT_ISSUER", "bigmind")
	os.Setenv("SERVER_ADDRESS", "localhost")
	os.Setenv("SERVER_PORT", "9052")
	os.Setenv("APP_ENV", "cloud")
	os.Setenv("REDIS_HOST", "")
	os.Setenv("REDIS_PORT", "6379")
	os.Setenv("REDIS_QUEUE", "xcheck")
	os.Setenv("CLOUD_BASE_URL", "")
}

func parseEnv(i interface{}) error {
	r := reflect.TypeOf(i)
	for r.Kind() == reflect.Ptr {
		r = r.Elem()
	}
	for i := 0; i < r.NumField(); i++ {
		env := r.Field(i).Tag.Get("mapstructure")
		if err := viper.BindEnv(env); err != nil {
			return err
		}
	}
	return viper.Unmarshal(i)
}

func ConnectToDB() (*gorm.DB, error) {
	var err error
	configEnv := GetConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", configEnv.GetString("DATABASE_USER"), configEnv.GetString("DATABASE_PASSWORD"), configEnv.GetString("DATABASE_HOST"), configEnv.GetString("DATABASE_PORT"), configEnv.GetString("DATABASE_NAME"))

	log.Println("Debug Query: ", configEnv.GetBool("DATABASE_QUERY_DEBUG"))

	gormConfig := &gorm.Config{}
	if configEnv.GetBool("DATABASE_QUERY_DEBUG") {
		dbLogger := dblogger.New(zap.L())
		dbLogger.SetAsDefault()
		dbLogger.LogMode(logger.Info)
		gormConfig = &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
			// Logger: dbLogger,
			// Logger: dbLogger,
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
		DatabaseHost:       GetConfig().GetString("DATABASE_HOST"),
		DatabasePort:       GetConfig().GetString("DATABASE_PORT"),
		DatabaseUser:       GetConfig().GetString("DATABASE_USER"),
		DatabasePassword:   GetConfig().GetString("DATABASE_PASSWORD"),
		DatabaseName:       GetConfig().GetString("DATABASE_NAME"),
		DatabaseQueryDebug: GetConfig().GetBool("DATABASE_QUERY_DEBUG"),
		AuthJwtSecret:      GetConfig().GetString("AUTH_JWT_SECRET"),
		AuthJwtExpire:      GetConfig().GetInt("AUTH_JWT_EXPIRE"),
		AuthJwtIssuer:      GetConfig().GetString("AUTH_JWT_ISSUER"),
		ServerAddress:      GetConfig().GetString("SERVER_ADDRESS"),
		ServerPort:         GetConfig().GetString("SERVER_PORT"),
		AppEnv:             GetConfig().GetString("APP_ENV"),
		RedisHost:          GetConfig().GetString("REDIS_HOST"),
		RedisPort:          GetConfig().GetString("REDIS_PORT"),
		RedisQueue:         GetConfig().GetString("REDIS_QUEUE"),
		CloudBaseURL:       GetConfig().GetString("CLOUD_BASE_URL"),
		MinioEndpoint:      GetConfig().GetString("MINIO_ENDPOINT"),
		MinioBucket:        GetConfig().GetString("MINIO_BUCKET"),
		MinioAccessKey:     GetConfig().GetString("MINIO_ACCESSKEY"),
		MinioSecretKey:     GetConfig().GetString("MINIO_SECRETKEY"),
	}
}
