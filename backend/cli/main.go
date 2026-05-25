package cli

import (
	"github.com/AshvinBambhaniya/tic-tac-toe/config"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// Init app initialization
func Init(cfg config.AppConfig, logger *zap.Logger) error {
	migrationCmd := GetMigrationCommandDef(cfg)
	apiCmd := GetAPICommandDef(cfg, logger)

	rootCmd := &cobra.Command{Use: "tic-tac-toe"}
	rootCmd.AddCommand(&migrationCmd, &apiCmd)
	return rootCmd.Execute()
}
