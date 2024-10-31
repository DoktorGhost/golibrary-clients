package server

import (
	"context"
	"net"
	"time"

	"google.golang.org/grpc"
)

const (
	_defaultShutdownTimeout = 5 * time.Second // Увеличим время для корректного завершения активных соединений
)

// GRPCServer обертка над grpc.Server для управления его жизненным циклом
type GRPCServer struct {
	server *grpc.Server
	errors chan error
	lis    net.Listener
}

// NewGRPCServer создаёт новый GRPC сервер
func NewGRPCServer(lis net.Listener, grpcServer *grpc.Server) *GRPCServer {
	return &GRPCServer{
		server: grpcServer,
		errors: make(chan error, 1),
		lis:    lis,
	}
}

// Serve запускает сервер в отдельной горутине и записывает ошибку в канал ошибок
func (s *GRPCServer) Serve() {
	go func() {
		s.errors <- s.server.Serve(s.lis)
		close(s.errors)
	}()
}

// Notify возвращает канал для получения ошибок запуска сервера
func (s *GRPCServer) Notify() <-chan error {
	return s.errors
}

// Shutdown выполняет graceful shutdown сервера, завершает активные соединения в течении timeout
func (s *GRPCServer) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), _defaultShutdownTimeout)
	defer cancel()

	done := make(chan struct{})
	go func() {
		s.server.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		s.server.Stop()
		return ctx.Err()
	}
}
