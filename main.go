package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/stdlib"

	"service-main/db"
	docs "service-main/docs"
	"service-main/handlers"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	defaultDBConnectTimeout = 45 * time.Second
	defaultDBRetryInterval  = 2 * time.Second
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

	dbConnectTimeout, err := envDuration("DB_CONNECT_TIMEOUT", defaultDBConnectTimeout)
	if err != nil {
		log.Fatal(err)
	}
	dbRetryInterval, err := envDuration("DB_CONNECT_RETRY_INTERVAL", defaultDBRetryInterval)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("connecting to database (timeout=%s, retry_interval=%s)", dbConnectTimeout, dbRetryInterval)
	pool, err := db.NewPoolWithRetry(ctx, connString, dbConnectTimeout, dbRetryInterval)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	sqlDB := stdlib.OpenDBFromPool(pool)
	defer sqlDB.Close()

	queries := db.New(sqlDB)

	router := gin.Default()
	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// This should be moved to a routes.go and handled there,
	// but for now this is fine

	router.POST("/officers", handlers.CreateOfficerHandler(queries))
	router.GET("/officers", handlers.GetOfficersHandler(queries))

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func envDuration(key string, fallback time.Duration) (time.Duration, error) {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback, nil
	}

	d, err := time.ParseDuration(raw)
	if err != nil {
		return 0, fmt.Errorf("%s must be a valid duration (for example: 30s or 2m): %w", key, err)
	}
	if d <= 0 {
		return 0, fmt.Errorf("%s must be greater than 0", key)
	}

	return d, nil
}
