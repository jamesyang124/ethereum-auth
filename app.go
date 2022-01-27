package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"

	"github.com/swaggo/fiber-swagger"
	_ "viveportengineering/DoubleA/ethereum-auth/docs"
	"viveportengineering/DoubleA/ethereum-auth/handlers"

	"log"
)

func loadNonEmptyEnv(key string, l *log.Logger) string {
	v := os.Getenv(key)
	if v == "" {
		l.Fatal("missing " + key + " env")
	} else {
		l.Printf("%s=%s", key, v)
	}

	return v
}

// @title           Ethereum Nonce-based Authentication Micro Service
// @version         1.0.0
// @description     Nonce-based auth with Ethereum digital wallet, datastore backed by redis
// @contact.name    James Yang
// @contact.email   james_yang@htc.com
// @host            localhost:3030
// @BasePath        /
func main() {
	// setup logger
	l := log.New(os.Stdout, "15:04:05 | ", 0)

	// load envs
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		godotenv.Load(".env")
	}

	appVersion := os.Getenv("APP_VERSION")
	if appVersion == "" {
		appVersion = "experiment"
	}

	signInTextTemplate := loadNonEmptyEnv("SIGNIN_TEXT_TEMPLATE", l)
	downstreamAuthUri := loadNonEmptyEnv("DOWNSTREAM_AUTH_URI", l)
	redisHost := loadNonEmptyEnv("REDIS_CACHE_HOST", l)
	redisPort := loadNonEmptyEnv("REDIS_CACHE_PORT", l)
	redisTTL := loadNonEmptyEnv("REDIS_CACHE_TTL_SECONDS", l)

	// setup redis client
	var ctx = context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: "",
	})

	// init and inject middleware before route/handler registration
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())

	app.Get("/swagger/*", fiberSwagger.WrapHandler)
	app.Get("/version", handlers.VersionHandler(appVersion))
	app.Get("/health", handlers.HealthCheckHandler)
	app.Get("/api/ethereum-auth/v1/metadata", handlers.MetadataHandler(signInTextTemplate))
	app.Post("/api/ethereum-auth/v1/nonce", handlers.NonceHandler(ctx, rdb, l, redisTTL))
	app.Post("/api/ethereum-auth/v1/login", handlers.LoginHandler(
		ctx,
		rdb,
		l,
		redisTTL,
		signInTextTemplate,
		downstreamAuthUri))

	app.Listen(":3030")
}
