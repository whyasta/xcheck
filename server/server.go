package server

import (
	"bigmind/xcheck-be/config"
	"bigmind/xcheck-be/repositories"
	"bigmind/xcheck-be/services"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type Server struct {
	DB     *gorm.DB
	Router *gin.Engine
}

func (server *Server) Init() {
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

	repositories := repositories.NewRepository(db)
	services := services.NewService(repositories)

	r := NewRouter(services)
	log.Printf("Starting server at port :" + configEnv.GetString("server.port") + "\n")
	r.Run(configEnv.GetString("server.address") + ":" + configEnv.GetString("server.port"))

	server.Router = r
	server.DB = db
}