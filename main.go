package main

import (
	"log"

	"service-main/auth"
	"service-main/db"
	docs "service-main/docs"
	"service-main/handlers"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	queries, cleanup, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	router := gin.Default()
	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// This should be moved to a routes.go and handled there,
	// but for now this is fine

	store := auth.SessionStore{Sessions: make(map[string]auth.Session)}

	go store.SessionCleanJob()

	{
		authed := router.Group("/auth")
		authed.Use(store.SessionGuard())
		authed.POST("/officers", handlers.CreateOfficerHandler(queries))
	}

	router.GET("/officers", handlers.GetOfficersHandler(queries))

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
