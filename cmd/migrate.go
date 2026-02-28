package cmd

import (
	"context"
	"log"
	"os"

	"inventory-system/internal/config"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Reset database and run latest migration",
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Setup Config & Infra
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

		ctx := context.Background()

		// ==========================================
		// 2. TOMBOL NUKLIR (Reset Database)
		// ==========================================
		logger.Warn("Dropping public schema to reset database...")

		resetQuery := `
			DROP SCHEMA public CASCADE;
			CREATE SCHEMA public;
			GRANT ALL ON SCHEMA public TO public;
		`
		_, err = dbPool.Exec(ctx, resetQuery)
		if err != nil {
			logger.Fatal("Failed to reset database schema", zap.Error(err))
		}
		logger.Info("Database reset successful!")

		// ==========================================
		// 3. BACA & JALANKAN FILE SQL TERBARU
		// ==========================================
		logger.Info("Applying new database schema...")

		// Pastikan path ini sesuai dengan lokasi file SQL lu!
		sqlFile := "migrations/0001_init_schema.sql"

		sqlBytes, err := os.ReadFile(sqlFile)
		if err != nil {
			logger.Fatal("Failed to read SQL file", zap.Error(err), zap.String("file", sqlFile))
		}

		// Eksekusi semua isi file SQL sekaligus
		_, err = dbPool.Exec(ctx, string(sqlBytes))
		if err != nil {
			logger.Fatal("Failed to execute migration", zap.Error(err))
		}

		logger.Info("✅ Migration executed successfully!")
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
