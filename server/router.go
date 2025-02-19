package server

import (
	"context"
	"github.com/Out-Of-India-Theory/oit-go-commons/app"
	"github.com/Out-Of-India-Theory/prarthana-automated-script/configuration"
	"github.com/Out-Of-India-Theory/prarthana-automated-script/controller/auth"
	"github.com/Out-Of-India-Theory/prarthana-automated-script/controller/ingestion"
	"github.com/Out-Of-India-Theory/prarthana-automated-script/middleware"
	"github.com/Out-Of-India-Theory/prarthana-automated-script/service/facade"
	"github.com/Out-Of-India-Theory/prarthana-automated-script/service/zoho"
	"net/http"
)

func registerRoutes(ctx context.Context, app *app.App, service facade.Service, configuration *configuration.Configuration) {
	basePath := app.Engine.Group("prarthana_script")
	app.Engine.GET("/health-check", ingestion.HealthCheck)
	authRepo := zoho.InitZohoService(ctx, configuration, &http.Client{})
	authMiddleware := middleware.InitAuthMiddleware(configuration, authRepo)
	//prarthana-script
	{
		prarthanaIngestionController := ingestion.InitIngestionController(ctx, service, configuration)
		prarthanaIngestionV1 := basePath.Group("v1")
		prarthanaIngestionV1.POST("/shloks", prarthanaIngestionController.ShlokIngestion)
		prarthanaIngestionV1.POST("/stotras", prarthanaIngestionController.StotraIngestion)
		prarthanaIngestionV1.POST("/prarthanas", prarthanaIngestionController.PrarthanaIngestion)
		prarthanaIngestionV1.POST("/deities", prarthanaIngestionController.DeityIngestion)
	}

	//auth
	{
		authController := auth.InitZohoAuthController(ctx, service, configuration)
		authGroup := basePath.Group("auth/v1")
		authGroup.GET("/zoho/auth", authController.GetAuthorizationURL)
		authGroup.GET("/zoho/callback", authController.HandleAuthCallback)
		authGroup.Use(authMiddleware.AuthMiddleware(ctx, authRepo), authController.CheckTokenMiddleware())
		authGroup.GET("/read-sheet/:sheetID", authController.ReadZohoSheet)
	}
}
