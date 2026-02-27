// @title           Inventory System API
// @version         1.0
// @description     This is a sample inventory system server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1
package cmd

import (
	"context"
	"os"
	"path/filepath"

	"inventory-system/internal/config"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long:  `Executes the SQL scripts inside the migrations folder to setup the database schema.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Load Configuration
		cfg, err := config.LoadConfig(".")
		if err != nil {
			panic("Failed to load configuration: " + err.Error())
		}

		// 2. Initialize Logger
		logger := config.InitLogger(cfg.App.Env)
		defer logger.Sync()

		logger.Info("Starting database migration...")

		// 3. Connect to Database
		dbPool, err := config.ConnectDB(cfg.DB.GetURL())
		if err != nil {
			logger.Fatal("Failed to connect to database", zap.Error(err))
		}
		defer dbPool.Close()

		// 4. Read the SQL file
		migrationPath := filepath.Join("migrations", "0001_init_schema.sql")
		sqlBytes, err := os.ReadFile(migrationPath)
		if err != nil {
			logger.Fatal("Failed to read migration file", zap.Error(err), zap.String("path", migrationPath))
		}

		// 5. Execute the SQL script
		ctx := context.Background()
		_, err = dbPool.Exec(ctx, string(sqlBytes))
		if err != nil {
			logger.Fatal("Failed to execute migration", zap.Error(err))
		}

		logger.Info("Database migration completed successfully! ✅")
	},
}

// init is automatically called by Go before main()
func init() {
	// Register the migrateCmd to the rootCmd
	rootCmd.AddCommand(migrateCmd)
}
