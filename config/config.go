package config

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var config *viper.Viper

// Init is an exported method that takes the environment starts the viper
// (external lib) and returns the configuration struct.
func Init(env string) {
	var err error
	config = viper.New()
	config.SetConfigType("yml")
	config.SetConfigName(env)
	config.AddConfigPath("../config/")
	config.AddConfigPath("config/")
	// dbConn, err := ConnectToDB()
	// if err != nil {
	// 	log.Fatal("Error connecting to database. Error: ", err)
	// 	return
	// }
	// config.Set("db", dbConn)
	err = config.ReadInConfig()
	if err != nil {
		log.Fatal("error on parsing configuration file")
	}
}

func ConnectToDBOld() (*sql.DB, error) {
	// var err error
	// dsn := os.Getenv("DB_DSN")

	// db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// if err != nil {
	//  log.Fatal("Error connecting to database. Error: ", err)
	// }

	// return db
	// "user:password@tcp(127.0.0.1:3306)/database_name"
	dsn := config.GetString("database.dsn")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func ConnectToDB() (*gorm.DB, error) {
	var err error
	configEnv := GetConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", configEnv.GetString("database.user"), configEnv.GetString("database.password"), configEnv.GetString("database.host"), configEnv.GetString("database.port"), configEnv.GetString("database.database"))

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database. Error: ", err)
		return nil, err
	}

	return db, err
	// "user:password@tcp(127.0.0.1:3306)/database_name"
}

// func relativePath(basedir string, path *string) {
// 	p := *path
// 	if len(p) > 0 && p[0] != '/' {
// 		*path = filepath.Join(basedir, p)
// 	}
// }

func GetConfig() *viper.Viper {
	return config
}