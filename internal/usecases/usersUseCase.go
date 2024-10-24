package usecases

import (
	"errors"
	"fmt"
	"github.com/DoktorGhost/golibrary-clients/internal/entities"
	"github.com/DoktorGhost/golibrary-clients/internal/repositories/postgres/dao"
	"github.com/DoktorGhost/golibrary-clients/internal/services"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"unicode"
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
	fullName, err := valid(userData.Name, userData.Surname, userData.Patronymic)
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

// вспомогалтельные функции
func validStr(name string) bool {
	for i, char := range name {
		// Проверяем, что символ является буквой или дефисом
		if !unicode.IsLetter(char) && char != '-' {
			return false
		}

		// Проверяем, что дефис не в начале или в конце строки
		if char == '-' && (i == 0 || i == len(name)-1) {
			return false
		}
	}
	return true
}

func valid(name, surname, patronymic string) (string, error) {
	if !validStr(name) {
		return "", fmt.Errorf("имя не должно содержать цифры или символы")
	}
	if !validStr(surname) {
		return "", fmt.Errorf("фамилия не должна содержать цифры или символы")
	}
	if !validStr(patronymic) {
		return "", fmt.Errorf("отчество не должно содержать цифры или символы")
	}
	fullName := strings.TrimSpace(name + " " + surname + " " + patronymic)
	return fullName, nil
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
