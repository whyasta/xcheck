package server

import (
	"bigmind/xcheck-be/config"
	"bigmind/xcheck-be/repositories"
	"bigmind/xcheck-be/services"
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Server struct {
	DB     *sql.DB
	Router *gin.Engine
}

func (server *Server) Init() {
	config := config.GetConfig()

	// setup db connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.GetString("database.user"), config.GetString("database.password"), config.GetString("database.host"), config.GetString("database.port"), config.GetString("database.database"))
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		fmt.Printf("Cannot connect to %s database"+"\n", "mysql")

		log.Fatal("Error:", err)
	} else {
		fmt.Printf("Connected to %s database successfully"+"\n", "mysql")
	}

	repositories := repositories.NewRepository(db)
	services := services.NewService(repositories)

	r := NewRouter(services)
	log.Printf("Starting server at port :" + config.GetString("server.port") + "\n")
	r.Run(config.GetString("server.address") + ":" + config.GetString("server.port"))

	server.Router = r
	server.DB = db
}
