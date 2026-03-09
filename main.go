package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"service-main/db"
)

func main() {
	ctx := context.Background()

	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	pool, err := db.NewPool(ctx, connString)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	_ = db.New(pool)

	router := gin.Default()
	router.Run(":8080")
}
