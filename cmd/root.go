package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kondohiroki/go-grpc-boilerplate/config"
	"github.com/kondohiroki/go-grpc-boilerplate/internal/db/pgx"
	"github.com/kondohiroki/go-grpc-boilerplate/internal/db/rdb"
	"github.com/kondohiroki/go-grpc-boilerplate/internal/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const defaultConfigFile = "config/config.yaml"

var RootCmdName = "main"

var configFile string
var rootCmd = &cobra.Command{
	Use: func() string {
		return RootCmdName
	}(),
	Short: "\nThis application is made with ❤️",
	Run: func(cmd *cobra.Command, _ []string) {
		cmd.Usage()
	},
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", fmt.Sprintf("config file (default is %s)", defaultConfigFile))
}

func setupAll() {
	setUpConfig()
	setUpLogger()
	// setUpPostgres()
	// setUpRedis()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("rootCmd.Execute() Error: %v", err)
		os.Exit(1)
	}
}

func setUpConfig() {
	if configFile == "" {
		configFile = defaultConfigFile
	}

	log.Default().Printf("Using config file: %s", configFile)
	config.SetConfig(configFile)
}

func setUpLogger() {
	log.Default().Printf("Using log level: %s", config.GetConfig().Log.Level)
	logger.InitLogger("zap")
}

func setUpPostgres() {
	// Create the database connection pool
	if config.GetConfig().Postgres.Host != "" {
		if config.GetConfig().Postgres.Schema == "" {
			logger.Log.Fatal("Postgres schema is not set")
		}

		// Initialize database schema if it doesn't exist
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		logger.Log.Info("Initializing database schema", zap.String("schema", config.GetConfig().Postgres.Schema))
		err := pgx.InitSchema(ctx, config.GetConfig().Postgres, config.GetConfig().Postgres.Schema)
		if err != nil {
			logger.Log.Fatal("pgx.InitSchema()", zap.Error(err))
		}

		logger.Log.Info("Initializing pgxPool")
		err = pgx.InitPgConnectionPool(config.GetConfig().Postgres)
		if err != nil {
			logger.Log.Fatal("pgx.InitPgConnectionPool()", zap.Error(err))
		}
		logger.Log.Info("pgxPool initialized")
	}

}

func setUpRedis() {
	// Create the database connection pool
	if config.GetConfig().Redis[0].Host != "" {
		logger.Log.Info("Initializing redis")
		err := rdb.InitRedisClient(config.GetConfig().Redis)
		if err != nil {
			logger.Log.Fatal("rdb.InitRedisClient()", zap.Error(err))
		}
		logger.Log.Info("redis initialized")
	}

}
