// Package main ...
// @title Global Auth Server API
// @version 1.0
// @description This is a sample server for the Global Auth Server API.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:4000
// @BasePath /api

// @externalDocs.description Swagger Open API Specification
// @externalDocs.url https://swagger.io/specification/

package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "global-auth-server/docs"
	"global-auth-server/routes"
	"os"
)

func main() {

	_ = godotenv.Load()

	// Set Gin mode from environment variable or default to release
	mode := os.Getenv("GIN_MODE")
	if mode == "" {
		mode = gin.ReleaseMode
	}
	gin.SetMode(mode)

	port := os.Getenv("PORT")
	host := os.Getenv("HOSTDEPLOY")
	if port == "" {
		port = "8080"
	}
	if host == "" {
		host = "0.0.0.0"
	}

	r := gin.New()
	r.Use(cors.Default())
	r.Use(gin.Logger(), gin.Recovery())
	r.LoadHTMLGlob("templates/*")
	// Set trusted proxies (for production, set your proxy IPs or use "localhost" for local dev)
	r.SetTrustedProxies([]string{"127.0.0.1", "localhost"})

	// App routes
	routes.RegisterRoutes(r)

	addr := fmt.Sprintf("%s:%s", host, port)
	fmt.Printf("Server running at http://%s\n", addr)
	r.Run(addr)
}
