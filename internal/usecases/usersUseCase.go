package usecases

import (
	"errors"
	"fmt"
	"github.com/DoktorGhost/golibrary-clients/internal/entities"
	"github.com/DoktorGhost/golibrary-clients/internal/repositories/postgres/dao"
	"github.com/DoktorGhost/golibrary-clients/internal/services"
	"github.com/DoktorGhost/platform/validator"
	"golang.org/x/crypto/bcrypt"
)

type UsersUseCase struct {
	userService *services.UserService
}

func NewUsersUseCase(userService *services.UserService) *UsersUseCase {
	return &UsersUseCase{userService: userService}
}

// регистрация пользователя
func (uc *UsersUseCase) AddUser(userData entities.RegisterData) (int, error) {
	// Проверка, существует ли пользователь с таким именем
	_, err := uc.userService.GetUserByUsername(userData.Username)
	if err == nil {
		return 0, fmt.Errorf("пользователь с таким Username уже существует")
	}

	// Валидация данных пользователя
	fullName, err := validator.Validator(userData.Name, userData.Surname, userData.Patronymic)
	if err != nil {
		return 0, fmt.Errorf("ошибка валидации данных: %v", err)
	}

	// Хеширование пароля
	hash, err := hashPassword(userData.Password)
	if err != nil {
		return 0, fmt.Errorf("ошибка хеширования пароля: %v", err)
	}

	// Подготовка данных для создания пользователя
	var data dao.UserTable

	data.Username = userData.Username
	data.PasswordHash = hash
	data.FullName = fullName

	// Создание пользователя
	id, err := uc.userService.CreateUser(data)
	if err != nil {
		return 0, fmt.Errorf("ошибка при создании пользователя: %v", err)
	}

	return id, nil
}

// возвращает из таблицы данные о пользователе, если логиn и пароль подходят
func (uc *UsersUseCase) Login(userData entities.Login) (dao.UserTable, error) {
	var user dao.UserTable
	user, err := uc.userService.GetUserByUsername(userData.Username)
	if err != nil {
		return user, errors.New("user not found")
	}

	err = checkPasswordHash(userData.Password, user.PasswordHash)
	if err != nil {
		return user, errors.New("invalid password" + err.Error())
	}

	return user, nil
}

func (uc *UsersUseCase) GetUserById(userID int) (string, error) {
	username, err := uc.userService.GetUserByID(userID)
	if err != nil {
		return "", errors.New("user not found")
	}

	return username, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func checkPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return err
	}

	return nil
}
