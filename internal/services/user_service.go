package services

import (
	"fmt"
	"github.com/DoktorGhost/golibrary-clients/internal/repositories/postgres/dao"
)

// UserRepository определяет методы для работы с пользователями
//
//go:generate mockgen -source=$GOFILE -destination=./mock_user.go -package=${GOPACKAGE}
type UserRepository interface {
	CreateUser(user dao.UserTable) (int, error)
	GetUserByUsername(username string) (dao.UserTable, error)
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(user dao.UserTable) (int, error) {
	userID, err := s.repo.CreateUser(user)
	if err != nil {
		return 0, fmt.Errorf("ошибка создания пользователя: %v", err)
	}
	return userID, nil
}

func (s *UserService) GetUserByUsername(username string) (dao.UserTable, error) {
	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return dao.UserTable{}, err
	}
	return user, nil
}
