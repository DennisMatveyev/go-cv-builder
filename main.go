package main

import (
	"context"
	"cv-builder/auth"
	"cv-builder/config"
	"cv-builder/db"
	"cv-builder/media"
	"cv-builder/user"
	"os/signal"
	"syscall"
	"time"

	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	db := db.InitializeDB()
	cfg := config.LoadConfig()

	app := fiber.New(fiber.Config{AppName: "CV Builder"})
	app.Use(cors.New())

	auth.SetupRoutes(app.Group("/auth"), db, cfg.JWTSecret)

	userGroup := app.Group("/user")
	userGroup.Use(auth.Middleware(db, cfg.JWTSecret))
	user.SetupRoutes(userGroup, db)

	mediaGroup := app.Group("/media")
	mediaGroup.Use(auth.Middleware(db, cfg.JWTSecret))
	media.SetupRoutes(mediaGroup, db, cfg)

	if err := os.MkdirAll(cfg.UploadsDir, 0755); err != nil {
		log.Fatal("Failed to start server, could not make uploads dir: ", err)
	}
	if err := os.MkdirAll(cfg.DownloadsDir, 0755); err != nil {
		log.Fatal("Failed to start server, could not make downloads dir: ", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := app.Listen(":" + cfg.AppPort); err != nil {
			log.Fatal("Failed to start server: ", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exited gracefully")
}
