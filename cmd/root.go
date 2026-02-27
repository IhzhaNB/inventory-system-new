package cmd

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"inventory-system/internal/config"
	"inventory-system/internal/handler"
	"inventory-system/internal/repository"
	"inventory-system/internal/router"
	"inventory-system/internal/service"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "inventory-system",
	Short: "Backend API for Inventory System",
	Run: func(cmd *cobra.Command, args []string) {

		// 1. INFRASTRUCTURE SETUP
		cfg, err := config.LoadConfig(".")
		if err != nil {
			log.Fatalf("Failed to load configuration: %v", err)
		}

		logger := config.InitLogger(cfg.App.Env)
		defer logger.Sync()
		logger.Info("Config and Logger successfully loaded!")

		dbPool, err := config.ConnectDB(cfg.DB.GetURL())
		if err != nil {
			logger.Fatal("Failed to connect to database", zap.Error(err))
		}
		defer dbPool.Close()
		logger.Info("Successfully connected to the database!")

		// 2. DEPENDENCY INJECTION (Wiring up the app)
		repos := repository.NewRepository(dbPool)
		services := service.NewService(repos, &cfg, logger)
		handlers := handler.NewHandler(services, logger)

		// 3. ROUTING & MIDDLEWARE SETUP
		r := router.SetupRoute(handlers)

		// 4. START HTTP SERVER & GRACEFUL SHUTDOWN
		srv := &http.Server{
			Addr:    ":" + cfg.App.Port,
			Handler: r,
		}

		go func() {
			logger.Info("🚀 Server is running", zap.String("port", cfg.App.Port))
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.Fatal("Server failed to start", zap.Error(err))
			}
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		logger.Info("Shutting down server gracefully...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			logger.Fatal("Server forced to shutdown", zap.Error(err))
		}

		logger.Info("Server exited properly. Bye!")
	},
}

// Execute is called by main.go to run the CLI
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
