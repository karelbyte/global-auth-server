package main

import (
	"fmt"
	"os"

	"global-auth-server/controllers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	r.Use(gin.Logger(), gin.Recovery())
	r.LoadHTMLGlob("templates/*")
	// Set trusted proxies (for production, set your proxy IPs or use "localhost" for local dev)
	r.SetTrustedProxies([]string{"127.0.0.1", "localhost"})

	r.GET("/", controllers.Home)

	// Route group with prefix /API
	api := r.Group("/api")
	{
		api.POST("/auth/login", controllers.Login)
	}

	addr := fmt.Sprintf("%s:%s", host, port)
	fmt.Printf("Server running at http://%s\n", addr)
	r.Run(addr)
}
