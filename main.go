package main

import (
	"autoelys_backend/database"
	"autoelys_backend/internal/handlers"
	"autoelys_backend/internal/middleware"
	"autoelys_backend/internal/repository"
	"autoelys_backend/internal/services"
	"autoelys_backend/internal/validation"
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "autoelys_backend/docs"
)

// @title AutoElys Backend API
// @version 1.0
// @description Backend API for AutoElys application
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@autoelys.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	if len(os.Args) > 1 {
		handleMigrationCommand()
		return
	}

	config := database.NewConfig()
	db, err := database.Connect(config)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()

	validate := validator.New()
	if err := validation.RegisterCustomValidators(validate); err != nil {
		log.Fatalf("Failed to register custom validators: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	passwordRepo := repository.NewPasswordResetRepository(db)
	brandRepo := repository.NewBrandRepository(db)
	automobileRepo := repository.NewAutomobileRepository(db)
	vehicleRepo := repository.NewVehicleRepository(db)
	emailService := services.NewEmailService()
	authHandler := handlers.NewAuthHandler(userRepo, passwordRepo, emailService, validate)
	brandHandler := handlers.NewBrandHandler(brandRepo, automobileRepo)
	vehicleHandler := handlers.NewVehicleHandler(vehicleRepo, validate)

	rateLimiter := middleware.NewRateLimiter(10, 5)

	router := gin.Default()

	// CORS configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello World"})
	})

	// Serve static files (logos, images, etc.)
	router.Static("/public", "./public")
	router.Static("/uploads", "./uploads")

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", rateLimiter.Limit(), authHandler.Register)
			auth.POST("/login", rateLimiter.Limit(), authHandler.Login)
			auth.POST("/forgot-password", rateLimiter.Limit(), authHandler.ForgotPassword)
			auth.POST("/reset-password", rateLimiter.Limit(), authHandler.ResetPassword)
			auth.GET("/me", middleware.AuthRequired(), authHandler.GetMe)
			auth.PUT("/me", middleware.AuthRequired(), authHandler.UpdateProfile)
		}

		brands := api.Group("/brands")
		{
			brands.GET("", brandHandler.GetAllBrands)
			brands.GET("/:id/automobiles", brandHandler.GetAutomobilesByBrand)
		}

		vehicles := api.Group("/vehicles")
		{
			vehicles.GET("", vehicleHandler.GetAllVehicles)
			vehicles.GET("/recommended", vehicleHandler.GetRecommendedVehicles)
			vehicles.GET("/:slug", vehicleHandler.GetVehicle)
		}

		userVehicles := api.Group("/user/vehicles")
		userVehicles.Use(middleware.AuthRequired())
		{
			userVehicles.GET("", vehicleHandler.GetUserVehicles)
			userVehicles.POST("", vehicleHandler.CreateVehicle)
			userVehicles.GET("/:uuid", vehicleHandler.GetVehicleByUUID)
			userVehicles.PUT("/:uuid", vehicleHandler.UpdateVehicle)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server is running on http://localhost:%s", port)
	log.Printf("Swagger documentation available at http://localhost:%s/swagger/index.html", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func handleMigrationCommand() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	if len(os.Args) < 2 {
		printMigrationUsage()
		os.Exit(1)
	}

	config := database.NewConfig()
	db, err := database.Connect(config)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()

	migrator := database.NewMigrator(db, "./migrations")

	command := os.Args[1]
	switch command {
	case "migrate:up":
		if err := migrator.Up(); err != nil {
			log.Fatalf("Migration up failed: %v", err)
		}
	case "migrate:down":
		if err := migrator.Down(); err != nil {
			log.Fatalf("Migration down failed: %v", err)
		}
	case "migrate:status":
		version, dirty, err := migrator.Version()
		if err != nil {
			log.Fatalf("Could not get migration version: %v", err)
		}
		fmt.Printf("Current version: %d, Dirty: %t\n", version, dirty)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printMigrationUsage()
		os.Exit(1)
	}
}

func printMigrationUsage() {
	fmt.Println("Usage:")
	fmt.Println("  go run main.go migrate:up      - Run all pending migrations")
	fmt.Println("  go run main.go migrate:down    - Rollback last migration")
	fmt.Println("  go run main.go migrate:status  - Show current migration version")
	fmt.Println("  go run main.go                 - Start the server")
}
