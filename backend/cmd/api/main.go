package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"

	httphandler "github.com/youruser/yourproject/internal/adapter/handler/http"
	"github.com/youruser/yourproject/internal/adapter/handler/http/middleware"
	"github.com/youruser/yourproject/internal/adapter/payment/cardtocard"
	"github.com/youruser/yourproject/internal/adapter/payment/vandar"
	"github.com/youruser/yourproject/internal/adapter/payment/zarinpal"
	"github.com/youruser/yourproject/internal/adapter/repository/postgres"
	"github.com/youruser/yourproject/internal/adapter/sms/senator"
	"github.com/youruser/yourproject/internal/adapter/storage/s3"
	"github.com/youruser/yourproject/pkg/logger"
	"github.com/youruser/yourproject/pkg/telemetry"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func main() {
	// 1. Initialize Logger
	logger.InitLogger()
	logger.Log.Info("Starting application...")

	// 2. Initialize Telemetry (OpenTelemetry)
	tp, err := telemetry.InitTracer()
	if err != nil {
		logger.Log.Fatal("Failed to init tracer", zap.Error(err))
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			logger.Log.Error("Error shutting down tracer provider", zap.Error(err))
		}
	}()

	// 3. Initialize Infrastructure
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	// Postgres
	dbPool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Log.Fatal("Unable to connect to database", zap.Error(err))
	}
	defer dbPool.Close()

	// Repositories
	userRepo := postgres.NewUserRepository(dbPool)

	// 4. Initialize Adapters
	s3Adapter, err := s3.NewS3Adapter()
	if err != nil {
		logger.Log.Error("Failed to init S3 adapter", zap.Error(err))
	}
	
	// Payment Adapters
	zarinpalAdapter := zarinpal.NewZarinpalAdapter(os.Getenv("ZARINPAL_MERCHANT_ID"))
	vandarAdapter := vandar.NewVandarAdapter(os.Getenv("VANDAR_API_KEY"))
	cardToCardAdapter := cardtocard.NewCardToCardAdapter()

	// SMS Adapter
	smsAdapter := senator.NewSenatorAdapter()

	// Handlers
	authHandler := httphandler.NewAuthHandler(smsAdapter, rdb, userRepo)
	wsHandler := httphandler.NewWebSocketHandler()
	go wsHandler.Run()

	// 5. Initialize Fiber App
	app := fiber.New(fiber.Config{
		AppName: "Go Clean Arch Boilerplate",
	})

	// Setup Swagger
	SetupSwagger(app)

	// 6. Middleware
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(helmet.New())
	app.Use(limiter.New(limiter.Config{
		Max:        20,
		Expiration: 30 * time.Second,
	}))
	// CSRF for non-API routes (if any) or configured for API
	app.Use(csrf.New(csrf.Config{
		KeyLookup: "header:X-CSRF-Token",
	}))
	app.Use(otelfiber.Middleware()) // OpenTelemetry Middleware

	// 7. Routes
	api := app.Group("/api")

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// WebSocket Route
	app.Use("/ws", httphandler.WebSocketMiddleware())
	app.Get("/ws", websocket.New(wsHandler.HandleWebSocket))

	// Test WebSocket Broadcast
	api.Post("/broadcast", func(c *fiber.Ctx) error {
		type Msg struct {
			Message string `json:"message"`
		}
		var msg Msg
		if err := c.BodyParser(&msg); err != nil {
			return c.Status(400).SendString(err.Error())
		}
		wsHandler.BroadcastMessage(msg.Message)
		return c.SendString("Message broadcasted")
	})

	// Auth Routes
	auth := api.Group("/auth")
	auth.Post("/otp/send", authHandler.SendOTP)
	auth.Post("/otp/verify", authHandler.VerifyOTP)
	auth.Get("/google/callback", authHandler.GoogleCallback)

	// 2FA Routes
	auth.Post("/2fa/setup", middleware.Protected(), authHandler.Setup2FA)
	auth.Post("/2fa/enable", middleware.Protected(), authHandler.Enable2FA)
	auth.Post("/2fa/verify", middleware.Protected(), authHandler.Verify2FALogin)

	// Example Protected Route
	api.Get("/protected", middleware.Protected(), func(c *fiber.Ctx) error {
		userID := c.Locals("user_id")
		return c.JSON(fiber.Map{"message": "Access granted", "user_id": userID})
	})

	// Example Upload Route (Simplified for boilerplate)
	api.Post("/upload", middleware.Protected(), func(c *fiber.Ctx) error {
		if s3Adapter == nil {
			return c.Status(500).JSON(fiber.Map{"error": "Storage not configured"})
		}

		fileHeader, err := c.FormFile("file")
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "No file uploaded"})
		}

		file, err := fileHeader.Open()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to open file"})
		}
		defer file.Close()

		// Magic Bytes Check
		buffer := make([]byte, 512)
		if _, err := file.Read(buffer); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to read file"})
		}
		// Reset file pointer
		if _, err := file.Seek(0, 0); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to reset file pointer"})
		}

		contentType := http.DetectContentType(buffer)
		// Allow specific types (expand as needed)
		if contentType != "image/jpeg" && contentType != "image/png" && contentType != "application/pdf" {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid file type: " + contentType})
		}

		key, err := s3Adapter.UploadFile(c.Context(), fileHeader.Filename, file)
		if err != nil {
			logger.Log.Error("Upload failed", zap.Error(err))
			return c.Status(500).JSON(fiber.Map{"error": "Upload failed"})
		}

		return c.JSON(fiber.Map{"url": key})
	})

	// Presigned URL Route
	api.Post("/upload/presigned", middleware.Protected(), func(c *fiber.Ctx) error {
		type Request struct {
			Filename string `json:"filename"`
		}
		var req Request
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
		}

		if s3Adapter == nil {
			return c.Status(500).JSON(fiber.Map{"error": "Storage not configured"})
		}

		key, err := s3Adapter.GeneratePresignedURL(c.Context(), req.Filename, 300)
		if err != nil {
			logger.Log.Error("Failed to generate presigned URL", zap.Error(err))
			return c.Status(500).JSON(fiber.Map{"error": "Failed to generate URL"})
		}
		return c.JSON(fiber.Map{"url": key})
	})

	// Payment Routes
	payments := api.Group("/payments", middleware.Protected())
	
	// Zarinpal Request
	payments.Post("/zarinpal/request", func(c *fiber.Ctx) error {
		type PaymentReq struct {
			Amount int64 `json:"amount"`
			Desc   string `json:"description"`
		}
		var req PaymentReq
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
		}
		
		// In prod: Get user email from context/db
		url, authority, err := zarinpalAdapter.RequestPayment(c.Context(), req.Amount, "http://localhost:3000/payment/callback", req.Desc, "user@example.com")
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		
		return c.JSON(fiber.Map{"payment_url": url, "authority": authority})
	})

	// Vandar Request
	payments.Post("/vandar/request", func(c *fiber.Ctx) error {
		type PaymentReq struct {
			Amount int64 `json:"amount"`
			Desc   string `json:"description"`
			Mobile string `json:"mobile"`
		}
		var req PaymentReq
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
		}

		url, token, err := vandarAdapter.RequestPayment(c.Context(), req.Amount, "http://localhost:3000/payment/callback/vandar", req.Desc, req.Mobile)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"payment_url": url, "token": token})
	})

	// Card to Card Submission
	payments.Post("/card-to-card", func(c *fiber.Ctx) error {
		type C2CReq struct {
			Amount      int64  `json:"amount"`
			ReceiptURL  string `json:"receipt_url"`
			Description string `json:"description"`
		}
		var req C2CReq
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
		}

		// In prod: Get User ID from JWT claims
		userID := "user-123" 
		
		txID, err := cardToCardAdapter.SubmitReceipt(c.Context(), userID, req.Amount, req.ReceiptURL, req.Description)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"transaction_id": txID, "status": "pending_approval"})
	})

	// 8. Graceful Shutdown
	go func() {
		if err := app.Listen(":8080"); err != nil {
			logger.Log.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Listen for shutdown signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info("Shutting down server...")

	// Give the server a deadline for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.Log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	// Close Redis connection
	if err := rdb.Close(); err != nil {
		logger.Log.Error("Error closing Redis connection", zap.Error(err))
	}

	logger.Log.Info("Server exited gracefully")
}
