// Package main starts the Data API
package main

import (
	"context"
	"fmt"
	"github.com/mfamador/grpc-input/internal/config"
	"github.com/mfamador/grpc-input/pkg/serverv1"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

func main() {
	log.Info().Timestamp().Msg("Test Client gRPC Server API")
	request := serverv1.PostRequest{Message: "message xpto"}
	conn, _ := grpc.Dial(fmt.Sprintf("localhost:%d", config.Config.Server.GrpcPort), grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	client := serverv1.NewServiceClient(conn)
	_, err := client.Post(context.Background(), &request)
	if err != nil {
		fmt.Println(err)
	}
}
