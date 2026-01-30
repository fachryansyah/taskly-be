package main

import (
	"fmt"
	"log"
	"tasklybe/internal/task"
	"tasklybe/internal/user"
	"tasklybe/pkg/db"

	_ "tasklybe/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

// @title Taskly API
// @version 1.0
// @description Taskly backend API
// @host localhost:3000
// @BasePath /
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file loaded; using environment variables")
	}

	app := fiber.New()

	fmt.Println("Connecting to database...")
	db.Connect()

	fmt.Println("Migrating table...")
	if err = db.DB.AutoMigrate(&task.Task{}); err != nil {
		log.Fatal("failed to migrate table:", err)
	}
	if err = db.DB.AutoMigrate(&user.User{}); err != nil {
		log.Fatal("failed to migrate table:", err)
	}
	fmt.Println("Table migrated!")

	task.RegisterTaskRoutes(app)
	user.RegisterUserRoutes(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Static("/openapi.json", "./docs/swagger.json")
	app.Static("/docs", "./public")

	log.Fatal(app.Listen(":4002"))
}
