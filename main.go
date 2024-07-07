package main

import (
	"bigmind/xcheck-be/config"
	"bigmind/xcheck-be/docs"
	"bigmind/xcheck-be/server"
	"bigmind/xcheck-be/utils"

	// "bigmind/xcheck-be/utils"
	"flag"
	"fmt"
	"os"
)

var newServer = server.Server{}

// @securityDefinitions.apikey  BearerAuth
// @in							header
// @name						Authorization
// @description                 Type "Bearer" followed by a space and JWT token.

func main() {
	environment := flag.String("e", "development", "")
	flag.Parse()
	fmt.Printf("environment\t: %s\n", *environment)

	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
		os.Exit(1)
	}

	utils.InitLogger(*environment)
	config.Init(*environment)
	// utils.Init()

	// set swagger info
	docs.SwaggerInfo.Title = "Swagger XCheck API"
	docs.SwaggerInfo.Description = "XCheck API."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = config.GetConfig().GetString("server.address") + ":" + config.GetConfig().GetString("server.port")
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	newServer.Init()
}
