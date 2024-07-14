package server

import (
	"bigmind/xcheck-be/config"
	"bigmind/xcheck-be/internal/services"
	"fmt"
	"log"

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
	log.Printf("Starting server " + configEnv.GetString("SERVER_ADDRESS") + " at port :" + configEnv.GetString("SERVER_PORT") + "\n")
	r.Run(configEnv.GetString("SERVER_ADDRESS") + ":" + configEnv.GetString("SERVER_PORT"))
}
