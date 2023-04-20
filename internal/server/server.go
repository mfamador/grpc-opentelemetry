// Package server contains the common HTTP server logic, e.g.:
//
//	routes, validators, error handling
package server

import (
	"fmt"
	"github.com/mfamador/grpc-input/internal/serverapi"
	"github.com/mfamador/grpc-input/pkg/serverv1"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"net"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

// Config defines the handler configuration
type Config struct {
	GrpcPort int `yaml:"grpcPort"`
}

// Start runs the gRPC server
func Start(grpcServer *grpc.Server, list net.Listener) error {
	log.Info().Msg("Test API gRPC server")
	return grpcServer.Serve(list)
}

// GetGRPCServer returns the gRPC server
func GetGRPCServer(conf Config) (*grpc.Server, net.Listener, error) {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()))
	serverv1.RegisterServiceServer(grpcServer, serverapi.NewServerService())

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.GrpcPort))
	if err != nil {
		return nil, nil, fmt.Errorf("error creating listener: %v", err)
	}
	return grpcServer, listener, nil
}

// RunApp runs the API
func RunApp(sConf Config) error {
	grpcServer, lis, err := GetGRPCServer(sConf)
	if err != nil {
		return fmt.Errorf("failed to build grpcServer: %v", err)
	}

	if e := Start(grpcServer, lis); e != nil {
		log.Fatal().Msgf("Failed to start the gRPC server")
	}

	return nil
}
