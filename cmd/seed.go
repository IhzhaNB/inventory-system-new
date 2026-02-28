package cmd

import (
	"context"
	"log"

	"inventory-system/internal/config"
	"inventory-system/internal/model"
	"inventory-system/pkg/utils"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Insert dummy data into database",
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Setup Infra (Sama kayak di root.go)
		cfg, err := config.LoadConfig(".")
		if err != nil {
			log.Fatalf("Failed to load configuration: %v", err)
		}

		logger := config.InitLogger(cfg.App.Env)
		defer logger.Sync()

		dbPool, err := config.ConnectDB(cfg.DB.GetURL())
		if err != nil {
			logger.Fatal("Failed to connect to database", zap.Error(err))
		}
		defer dbPool.Close()

		logger.Info("Starting database seeding...")

		// 2. Bikin Password Hash (Biar bisa login beneran)
		hashedPassword, err := utils.HashPassword("password123") // Ini password aslinya
		if err != nil {
			logger.Fatal("Failed to hash password", zap.Error(err))
		}

		// 3. Insert Data User (Admin)
		adminID := uuid.New()
		adminName := "Super Admin"
		adminEmail := "superadmin@gmail.com"
		adminRole := model.RoleSuperAdmin
		query := `
			INSERT INTO users (id, name, email, password_hash, role)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (email) DO NOTHING
		`

		_, err = dbPool.Exec(context.Background(), query, adminID, adminName, adminEmail, hashedPassword, adminRole)
		if err != nil {
			logger.Fatal("Failed to seed admin user", zap.Error(err))
		}

		logger.Info("✅ Successfully seeded super admin user!")
		logger.Info("📧 Email: superadmin@gmail.com")
		logger.Info("🔑 Password: password123")
	},
}

func init() {
	// Daftarin command seed ke root command
	rootCmd.AddCommand(seedCmd)
}
