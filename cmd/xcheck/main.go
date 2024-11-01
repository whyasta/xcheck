package main

import (
	"bigmind/xcheck-be/cmd/xcheck/server"
	"bigmind/xcheck-be/config"
	"log"

	// "bigmind/xcheck-be/utils"
	"flag"
	"fmt"
	"os"
	// "github.com/swaggo/swag/example/basic/docs"
)

func main() {
	environment := flag.String("e", "development", "")
	flag.Parse()
	fmt.Printf("environment\t: %s\n", *environment)

	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
		os.Exit(1)
	}

	config.InitLogger(*environment)
	config.Init(*environment)
	// utils.Init()

	// set swagger info
	// docs.SwaggerInfo.Title = "Swagger XCheck API"
	// docs.SwaggerInfo.Description = "XCheck API."
	// docs.SwaggerInfo.Version = "1.0"
	// docs.SwaggerInfo.Host = config.GetConfig().GetString("SERVER_ADDRESS") + ":" + config.GetConfig().GetString("SERVER_PORT")
	// docs.SwaggerInfo.BasePath = "/"
	// docs.SwaggerInfo.Schemes = []string{"http", "https"}

	serverEnv := config.GetConfig().GetString("APP_ENV")
	log.Printf("server environment: %s\n", serverEnv)

	server.Init()
}
