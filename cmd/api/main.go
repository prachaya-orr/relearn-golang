package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/prachaya-orr/relearn-golang/docs" // Import generated docs
	"github.com/prachaya-orr/relearn-golang/internal/domain"
	"github.com/prachaya-orr/relearn-golang/internal/handler"
	"github.com/prachaya-orr/relearn-golang/internal/middleware"
	"github.com/prachaya-orr/relearn-golang/internal/repository"
	"github.com/prachaya-orr/relearn-golang/internal/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title Go CRUD API
// @version 1.0
// @description This is a sample server for a Go CRUD API.
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// 0. Load Environment Variables
	envFile := flag.String("env", ".env", "Path to environment file")
	flag.Parse()

	if err := godotenv.Load(*envFile); err != nil {
		log.Printf("No %s file found, relying on environment variables", *envFile)
	}

	// 1. Setup Logger (Colorized)
	// newLogger := logger.New(
	// 	log.New(os.Stdout, "\r\n", log.LstdFlags),
	// 	logger.Config{
	// 		SlowThreshold:             time.Second,
	// 		LogLevel:                  logger.Info,
	// 		IgnoreRecordNotFoundError: true,
	// 		ParameterizedQueries:      false,
	// 		Colorful:                  true,
	// 	},
	// )

	// 2. Setup Database Connection
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Logger: newLogger,
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 3. Auto Migrate
	if err := db.AutoMigrate(&domain.Todo{}, &domain.User{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	fmt.Println("Database migrated successfully.")

	// 4. Dependency Injection
	repo := repository.NewTodoRepository(db)
	svc := service.NewTodoService(repo)
	h := handler.NewTodoHandler(svc)

	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)

	// 5. Setup Router
	gin.ForceConsoleColor()
	r := gin.Default()
	// Fix "You trusted all proxies" warning
	r.SetTrustedProxies(nil)

	// Swagger Route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Middleware
	r.Use(middleware.ResponseInterceptor())

	// 5. Register Routes
	// Auth Routes
	r.POST("/signup", userHandler.SignUp)
	r.POST("/login", userHandler.Login)
	r.POST("/refresh-token", userHandler.RefreshToken)

	// Todo Routes (Protected)
	todoRoutes := r.Group("/todos")
	todoRoutes.Use(middleware.AuthMiddleware())
	{
		todoRoutes.POST("", h.Create)
		todoRoutes.GET("", h.FindAll)
		todoRoutes.GET("/:id", h.FindByID)
		todoRoutes.PUT("/:id", h.Update)
		todoRoutes.DELETE("/:id", h.Delete)
		todoRoutes.DELETE("", h.DeleteAll)
	}

	// 6. Start Server with Graceful Shutdown
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
