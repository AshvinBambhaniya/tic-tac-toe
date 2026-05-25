package routes

import (
	"fmt"
	"sync"

	"go.uber.org/zap"

	"github.com/AshvinBambhaniya/tic-tac-toe/config"
	"github.com/AshvinBambhaniya/tic-tac-toe/constants"
	controller "github.com/AshvinBambhaniya/tic-tac-toe/controllers/api/v1"
	"github.com/AshvinBambhaniya/tic-tac-toe/middlewares"
	"github.com/doug-martin/goqu/v9"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

var mu sync.Mutex

// Setup func
func Setup(app *fiber.App, goqu *goqu.Database, logger *zap.Logger, config config.AppConfig) error {
	mu.Lock()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     config.FrontendURL,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, x-workspace-id, X-Requested-With",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	}))

	// TODO: add swagger docs
	// app.Use(swagger.New(swagger.Config{
	// 	BasePath: "/api/v1/",
	// 	FilePath: "./assets/swagger.json",
	// 	Path:     "docs",
	// 	Title:    "Swagger API Docs",
	// }))

	router := app.Group("/api")
	v1 := router.Group("/v1")

	middlewares := middlewares.NewMiddleware(config, logger)

	err := setupAuthController(v1, goqu, logger, middlewares, config)
	if err != nil {
		return err
	}

	err = setupUserController(v1, goqu, logger, middlewares)
	if err != nil {
		return err
	}

	err = healthCheckController(app, goqu, logger)
	if err != nil {
		return err
	}

	mu.Unlock()
	return nil
}

func setupAuthController(v1 fiber.Router, goqu *goqu.Database, logger *zap.Logger, middlewares middlewares.Middleware, config config.AppConfig) error {
	authController, err := controller.NewAuthController(goqu, logger, config)
	if err != nil {
		return err
	}
	v1.Post("/login", authController.DoAuth)
	v1.Get("/logout", authController.DoLogout)

	return nil
}

func setupUserController(v1 fiber.Router, goqu *goqu.Database, logger *zap.Logger, middlewares middlewares.Middleware) error {
	userController, err := controller.NewUserController(goqu, logger)
	if err != nil {
		return err
	}

	userRouter := v1.Group("/users")
	userRouter.Post("/", userController.CreateUser)
	userRouter.Get("/me", middlewares.Authenticated, userController.GetMe)
	userRouter.Get(fmt.Sprintf("/:%s", constants.ParamUid), middlewares.Authenticated, userController.GetUser)
	return nil
}

func healthCheckController(app *fiber.App, goqu *goqu.Database, logger *zap.Logger) error {
	healthController, err := controller.NewHealthController(goqu, logger)
	if err != nil {
		return err
	}

	healthz := app.Group("/healthz")
	healthz.Get("/", healthController.Overall)
	healthz.Get("/db", healthController.Db)
	return nil
}
