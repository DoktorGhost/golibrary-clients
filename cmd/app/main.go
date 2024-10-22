package main

import (
	"fmt"
	proto "github.com/DoktorGhost/external-api/src/go/pkg/grpc/clients/api/grpc/protobuf/clients_v1"
	"github.com/DoktorGhost/golibrary-clients/config"
	"github.com/DoktorGhost/golibrary-clients/internal/app"
	"github.com/DoktorGhost/golibrary-clients/internal/delivery/controllers/handlers"
	deliveryGrpc "github.com/DoktorGhost/golibrary-clients/internal/delivery/grpc/server"
	"github.com/DoktorGhost/golibrary-clients/internal/delivery/http/server"
	"github.com/DoktorGhost/platform/logger"
	"github.com/DoktorGhost/platform/storage/psg"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//инициализация логгера
	log, err := logger.GetLogger()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer log.Sync()

	//загрзка переменных окружения
	viper.AutomaticEnv()

	// Конвертируем в конфигурацию, которую ожидает InitStorage
	psgConfig := config.ConvertToPsgDBConfig(config.LoadConfig().DBConfig)

	// Инициализируем подключение к БД
	pgsqlConnector, err := psg.InitStorage(psgConfig)
	if err != nil {
		log.Error(err.Error())
		return
	}
	log.Info("соединение с БД установлено")
	defer pgsqlConnector.Close()

	cont := app.Init(pgsqlConnector)

	r := handlers.SetupRoutes(cont.UseCaseProvider)

	//старт grpc сервера
	lis, err := net.Listen("tcp", ":"+config.LoadConfig().ProviderConfig.Provider_port)
	if err != nil {
		log.Fatal("failed to listen: %v", "err", err)
	}

	grpcServer := grpc.NewServer()
	userGRPCServer := deliveryGrpc.NewUserGRPCServer(cont.UseCaseProvider.UserUseCase)

	proto.RegisterClientsServiceServer(grpcServer, userGRPCServer)
	reflection.Register(grpcServer)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal("failed to serve gRPC:", "err", err)
		}
	}()

	//старт http-сервера
	httpServer := server.NewHttpServer(r, ":8080")
	httpServer.Serve()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT)

	select {
	case killSignal := <-interrupt:
		log.Info("Выключение сервера", "signal", killSignal)
	case err = <-httpServer.Notify():
		log.Error("Ошибка сервера", "error", err)
	}

	httpServer.Shutdown()

}
