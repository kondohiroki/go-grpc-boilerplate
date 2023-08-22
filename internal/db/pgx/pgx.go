package pgx

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kondohiroki/go-grpc-boilerplate/config"
	"github.com/kondohiroki/go-grpc-boilerplate/internal/logger"
	"go.uber.org/zap"
)

var (
	readPgxPool  *pgxpool.Pool
	writePgxPool *pgxpool.Pool
	m            sync.Mutex
)

func InitPgConnectionPool(cfg config.Postgres) error {
	m.Lock()
	defer m.Unlock()

	// If no read config is provided
	// OR both configs point to the same host,
	// use a single connection pool.
	if cfg.Read.Host == "" || cfg.Read.Host == cfg.Write.Host {
		singlePool, err := initSinglePool(cfg.Write)
		if err != nil {
			return err
		}

		readPgxPool = singlePool
		writePgxPool = singlePool
		return nil
	}

	// For distinct read/write databases
	readPool, err := initSinglePool(cfg.Read)
	if err != nil {
		return err
	}
	writePool, err := initSinglePool(cfg.Write)
	if err != nil {
		return err
	}

	readPgxPool = readPool
	writePgxPool = writePool
	return nil
}

func GetReadPgxPool() *pgxpool.Pool {
	m.Lock()
	defer m.Unlock()

	if readPgxPool == nil {
		logger.Log.Info("Initializing readPgxPool again")
		var err error
		readPgxPool, err = initSinglePool(config.GetConfig().Postgres.Read)
		if err != nil {
			logger.Log.Error("Failed to initialize readPgxPool", zap.Error(err))
			return nil // added return to avoid using an uninitialized pool
		}
		logger.Log.Info("readPgxPool initialized")
	}
	return readPgxPool
}

func GetWritePgxPool() *pgxpool.Pool {
	m.Lock()
	defer m.Unlock()

	if writePgxPool == nil {
		logger.Log.Info("Initializing writePgxPool again")
		var err error
		writePgxPool, err = initSinglePool(config.GetConfig().Postgres.Write)
		if err != nil {
			logger.Log.Error("Failed to initialize writePgxPool", zap.Error(err))
			return nil // added return to avoid using an uninitialized pool
		}
		logger.Log.Info("writePgxPool initialized")
	}
	return writePgxPool
}

// initSinglePool initializes a single pool without acquiring a lock
func initSinglePool(postgresConfig config.PostgresConfig) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable search_path=%s",
		postgresConfig.Host,
		postgresConfig.Port,
		postgresConfig.Username,
		postgresConfig.Password,
		postgresConfig.Database,
		postgresConfig.Schema,
	)

	connConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		fmt.Println("Failed to parse config:", err)
		return nil, err
	}

	// Set maximum number of connections
	connConfig.MaxConns = postgresConfig.MaxConnections

	pgxPool, err := pgxpool.NewWithConfig(context.Background(), connConfig)
	if err != nil {
		return nil, err
	}
	return pgxPool, nil
}

func InitSchema(ctx context.Context, postgresConfig config.PostgresConfig, schema string) (err error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		postgresConfig.Host,
		postgresConfig.Port,
		postgresConfig.Username,
		postgresConfig.Password,
		postgresConfig.Database,
	)

	pgConn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return err
	}
	defer pgConn.Close(ctx)

	// Create schema if it doesn't exist
	// Ignore error if schema already exists or if the user doesn't have permission to create schema
	pgConn.Exec(
		ctx,
		fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS %s`, schema),
	)

	// Set search path to schema so that we don't have to specify the schema name
	_, err = pgConn.Exec(
		ctx,
		fmt.Sprintf(`SET search_path TO %s`, schema),
	)
	if err != nil {
		return err
	}

	return nil
}

func ClosePgxPool() {
	m.Lock()
	defer m.Unlock()

	if readPgxPool != nil {
		fmt.Println("Closing readPgxPool")
		readPgxPool.Close()
		readPgxPool = nil
		fmt.Println("readPgxPool closed")
	}

	if writePgxPool != nil {
		fmt.Println("Closing writePgxPool")
		writePgxPool.Close()
		writePgxPool = nil
		fmt.Println("writePgxPool closed")
	}
}
