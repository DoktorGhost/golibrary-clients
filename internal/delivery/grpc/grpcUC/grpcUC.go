package grpcUC

import (
	"context"
	proto "github.com/DoktorGhost/external-api/src/go/pkg/grpc/clients/api/grpc/protobuf/clients_v1"
	"github.com/DoktorGhost/golibrary-clients/internal/entities"
	"github.com/DoktorGhost/golibrary-clients/internal/usecases"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserGRPCServer struct {
	uc *usecases.UsersUseCase
	proto.ClientsServiceServer
}

func NewUserGRPCServer(uc *usecases.UsersUseCase) *UserGRPCServer {
	return &UserGRPCServer{uc: uc}
}

// / controllers
func (s *UserGRPCServer) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	user := entities.RegisterData{
		Username:   req.Username,
		Password:   req.Password,
		Name:       req.Name,
		Surname:    req.Surname,
		Patronymic: req.Patronymic,
	}

	// Вызов метода юзкейса
	id, err := s.uc.AddUser(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add user: %v", err)
	}

	return &proto.RegisterResponse{Id: int64(id)}, nil
}

func (s *UserGRPCServer) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	loginData := entities.Login{
		Username: req.Username,
		Password: req.Password,
	}
	user, err := s.uc.Login(loginData)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "login failed: %v", err)
	}

	return &proto.LoginResponse{Id: int64(user.ID), Username: user.Username, Password: user.PasswordHash, Fullname: user.FullName}, nil
}

func (s *UserGRPCServer) GetUserByID(ctx context.Context, req *proto.UserID) (*proto.Username, error) {
	id := int(req.Id)
	username, err := s.uc.GetUserById(id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}

	return &proto.Username{Username: username}, nil
}
