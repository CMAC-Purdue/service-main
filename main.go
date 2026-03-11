package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/stdlib"

	"service-main/db"
	"service-main/handlers"
)

func main() {
	ctx := context.Background()

	if err := loadDotEnv(".env"); err != nil {
		log.Fatal(err)
	}

	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	pool, err := db.NewPool(ctx, connString)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	sqlDB := stdlib.OpenDBFromPool(pool)
	defer sqlDB.Close()

	queries := db.New(sqlDB)

	router := gin.Default()
	router.POST("/officers", handlers.CreateOfficerHandler(queries))
	router.Run(":8080")
}
