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
	postgresConfig := config.GetConfig().Postgres

	// Check if either Read or Write host is provided
	if postgresConfig.Read.Host != "" || postgresConfig.Write.Host != "" {
		// Ensure both Read and Write hosts are provided
		if postgresConfig.Read.Schema == "" || postgresConfig.Write.Schema == "" {
			logger.Log.Fatal("Both Read and Write schema must be set for Postgres")
		}

		// Initialize database schema if it doesn't exist
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		logger.Log.Info(fmt.Sprintf("Initializing schema: %s", postgresConfig.Write.Schema))
		err := pgx.InitSchema(ctx, postgresConfig.Write, postgresConfig.Write.Schema)
		if err != nil {
			logger.Log.Fatal("Failed in pgx.InitSchema()", zap.Error(err))
		}

		// Initialize the connection pool
		logger.Log.Info("Initializing pgxPool")
		err = pgx.InitPgConnectionPool(postgresConfig)
		if err != nil {
			logger.Log.Fatal("Failed in pgx.InitPgConnectionPool()", zap.Error(err))
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
