package server

import "github.com/DoktorGhost/golibrary-clients/internal/usecases"

import (
	"context"
	"errors"
	proto "github.com/DoktorGhost/external-api/src/go/pkg/grpc/clients/api/grpc/protobuf/clients_v1"
	"github.com/DoktorGhost/golibrary-clients/internal/entities"
)

type UserGRPCServer struct {
	uc *usecases.UsersUseCase
	proto.ClientsServiceServer
}

func NewUserGRPCServer(uc *usecases.UsersUseCase) *UserGRPCServer {
	return &UserGRPCServer{uc: uc}
}

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
		return nil, errors.New("error adding user: " + err.Error())
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
		return &proto.LoginResponse{}, err
	}

	return &proto.LoginResponse{Id: int64(user.ID), Username: user.Username, Password: user.PasswordHash, Fullname: user.FullName}, nil
}
