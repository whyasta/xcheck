package server

import (
	"bigmind/xcheck-be/config"
	"bigmind/xcheck-be/internal/services"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// type Server struct {
// 	DB     *gorm.DB
// 	Router *gin.Engine
// }

func Init() {
	configEnv := config.GetConfig()

	// setup db connection
	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", configEnv.GetString("database.user"), configEnv.GetString("database.password"), configEnv.GetString("database.host"), configEnv.GetString("database.port"), configEnv.GetString("database.database"))
	// db, err := sql.Open("mysql", dsn)

	// if err != nil {
	// 	fmt.Printf("Cannot connect to %s database"+"\n", "mysql")

	// 	log.Fatal("Error:", err)
	// } else {
	// 	fmt.Printf("Connected to %s database successfully"+"\n", "mysql")
	// }

	db, err := config.ConnectToDB()
	if err != nil {
		fmt.Printf("Cannot connect to %s database"+"\n", "mysql")

		log.Fatal("Error:", err)
	} else {
		fmt.Printf("Connected to %s database successfully"+"\n", "mysql")
	}

	// setup services
	services := services.RegisterServices(db)

	r := NewRouter(services)
	log.Printf("Starting server " + configEnv.GetString("APP_ENV") + ":" + configEnv.GetString("SERVER_ADDRESS") + " at port :" + configEnv.GetString("SERVER_PORT") + "\n")
	//r.Run(configEnv.GetString("SERVER_ADDRESS") + ":" + configEnv.GetString("SERVER_PORT"))
	srv := &http.Server{
		Addr:         configEnv.GetString("SERVER_ADDRESS") + ":" + configEnv.GetString("SERVER_PORT"),
		WriteTimeout: time.Second * 60,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}
	srv.ListenAndServe()
}
