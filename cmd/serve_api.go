package cmd

import (
	"context"
	"fmt"
	"net"
	"os/signal"
	"syscall"

	"github.com/kondohiroki/go-grpc-boilerplate/config"
	"github.com/kondohiroki/go-grpc-boilerplate/internal/logger"
	"github.com/kondohiroki/go-grpc-boilerplate/internal/server"
	"github.com/kondohiroki/go-grpc-boilerplate/pkg/validation"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func init() {
	rootCmd.AddGroup(&cobra.Group{ID: "serve", Title: "Serve:"})
	rootCmd.AddCommand(serveGRPCAPICmd)
}

var serveGRPCAPICmd = &cobra.Command{
	Use:     "serve:grpc-api",
	Short:   "Start the gRPC API",
	GroupID: "serve",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Setup all the required dependencies
		setupAll()

		// Create gRPC server
		gRPCServer, err := server.NewGRPCServer()
		if err != nil {
			return fmt.Errorf("failed to create gRPC server: %w", err)
		}

		// Create validator instance
		validation.InitValidator()

		// Get port from config
		port := config.GetConfig().GRPCServer.Port

		// Create a listener
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			return fmt.Errorf("failed to listen: %w", err)
		}

		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		localIP, _ := getLocalIP()
		go func() {
			logger.Log.Info(fmt.Sprintf("Starting gRPC server on port %d", port))
			logger.Log.Info(fmt.Sprintf("Local: grpc://localhost:%d", port))
			logger.Log.Info(fmt.Sprintf("Network: grpc://%s:%d", localIP, port))
			logger.Log.Info("waiting for requests...")

			// Start gRPC server
			if err := gRPCServer.Serve(lis); err != nil && err != grpc.ErrServerStopped {
				logger.Log.Fatal(fmt.Sprintf("failed to serve: %s\n", err))
			}
		}()

		<-ctx.Done()
		stop()
		logger.Log.Info("\nShutting down gracefully, press Ctrl+C again to force")

		// Gracefully stop the gRPC server
		gRPCServer.GracefulStop()

		return nil
	},
}

// ... (rest of the code remains unchanged)

func getLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("local IP not found")
}
