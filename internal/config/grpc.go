package config

import (
	"errors"
	"net"
	"os"
)

const (
	grpcHostName = "GRPC_HOST"
	grpcPortName = "GRPC_PORT"
)

type GRPCConfig interface {
	Address() string
}

type grpcConfig struct {
	host string
	port string
}

func (cfg *grpcConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}

func NewGRPCConfig() (GRPCConfig, error) {
	host := os.Getenv(grpcHostName)
	if len(host) == 0 {
		return nil, errors.New("environment variable GRPC_HOST is not set")
	}

	port := os.Getenv(grpcPortName)
	if len(port) == 0 {
		return nil, errors.New("environment variable GRPC_PORT is not set")
	}

	return &grpcConfig{host: host, port: port}, nil
}
