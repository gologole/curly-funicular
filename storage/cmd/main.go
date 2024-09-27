package main

import (
	"cmd/main.go/configs"
	"cmd/main.go/internal/api"
	"cmd/main.go/internal/api/rpc"
	"cmd/main.go/internal/cache"
	db2 "cmd/main.go/internal/db"
	"cmd/main.go/internal/service"
	"cmd/main.go/pkg/logger"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"strconv"
)

func main() {
	cfg, err := configs.InitConfig("./configs/config.yaml")
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	elkLogger := logger.ConfigureLogger(cfg.HttpServer.ElkDomain)
	defer elkLogger.Sync()

	db, err := db2.NewDatabase(cfg)

	er := db.Migrate()
	if er != nil {
		panic("can't migrate :    " + er.Error())
	}
	cache := cache.NewCache(cfg)

	myservice := service.NewService(cfg, db, cache)

	// Настройка прослушивания на определенном порту
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(cfg.GrpcServer.Port)) // Например, ":50051"
	if err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}

	// Создаем новый gRPC сервер
	grpcServer := grpc.NewServer()

	// Регистрируем наш gRPC сервер
	storageServer := api.NewGrpcServer(cfg, myservice, elkLogger)
	rpc.RegisterStorageServer(grpcServer, storageServer)

	// Включаем возможность проверки зарегистрированных сервисов через gRPC reflection (для отладки)
	reflection.Register(grpcServer)

	// Запускаем gRPC сервер
	elkLogger.Info("Server registered")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Ошибка запуска gRPC сервера: %v", err)
	}
}
