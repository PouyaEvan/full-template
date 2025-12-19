package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	_ "github.com/youruser/yourproject/docs" // Import generated docs
)

// @title           Go Fiber Clean Architecture API
// @version         1.0
// @description     This is a production-ready boilerplate server.
// @termsOfService  http://swagger.io/terms/

// @contact.name    API Support
// @contact.email   support@swagger.io

// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html

// @host            localhost:8080
// @BasePath        /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func SetupSwagger(app *fiber.App) {
	app.Get("/swagger/*", swagger.HandlerDefault)
}
