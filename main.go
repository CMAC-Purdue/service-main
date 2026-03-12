// @title Service Main API
// @version 1.0
// @description CMAC's Go API
// @BasePath /
package main

import (
	"log"
	"os"

	"service-main/auth"
	"service-main/db"
	docs "service-main/docs"
	"service-main/handlers"
	"service-main/util"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {

	if err := util.LoadDotEnv(".env"); err != nil {
		log.Fatal("Failed to load .env")
	}

	queries, cleanup, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	router := gin.Default()
	router.SetTrustedProxies(nil) // gin throws a warning if this is not explicitly disabled

	store := auth.SessionStore{Sessions: make(map[string]auth.Session)}
	admin_phrase, exists := os.LookupEnv("ADMIN_PHRASE")

	if !exists {
		log.Fatal("ADMIN_PHRASE is not set")
	}

	go store.SessionCleanJob()

	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// This should be moved to a routes.go and handled there,
	// but for now this is fine

	{
		authed := router.Group("/auth")
		authed.Use(store.SessionGuard())
		authed.POST("/officers", handlers.CreateOfficerHandler(queries))

		authed.GET("/sessions", handlers.DisplaySessionsHandler(&store))

	}

	router.GET("/officers", handlers.GetOfficersHandler(queries))
	router.POST("/opme", handlers.AdminSessionLogin(&store, admin_phrase))

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
