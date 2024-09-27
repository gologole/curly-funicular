package main

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	pb "main.go/proto" // Замените на путь к вашему сгенерированному пакету
)

const jwtSecret = "your_secret_key" // Замените на свой секретный ключ

type server struct {
	pb.UnimplementedAuthServiceServer
}

func (s *server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if req.Username == "user" && req.Password == "password" { // Простейшая аутентификация
		token, err := generateJWT(req.Username)
		if err != nil {
			return nil, err
		}
		return &pb.LoginResponse{Token: token}, nil
	}
	return nil, fmt.Errorf("invalid credentials")
}

func generateJWT(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}
func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, &server{})

	reflection.Register(grpcServer)

	log.Println("Server is running on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
