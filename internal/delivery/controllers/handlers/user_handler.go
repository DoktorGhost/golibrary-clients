package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/DoktorGhost/golibrary-clients/internal/entities"
	"github.com/DoktorGhost/golibrary-clients/internal/providers"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strconv"
)

func handlerAddUser(useCaseProvider *providers.UseCaseProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Неправильный метод", http.StatusMethodNotAllowed)
			return
		}

		// Чтение тела запроса
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Ошибка чтения тела запроса", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Декодирование JSON из тела запроса
		var user entities.RegisterData
		if err := json.Unmarshal(body, &user); err != nil {
			http.Error(w, "Ошибка декодирования JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Вызов метода добавления автора из юзкейса
		id, err := useCaseProvider.UserUseCase.AddUser(user)
		if err != nil {
			http.Error(w, "Ошибка при добавлении пользователя: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Успешный ответ
		w.WriteHeader(http.StatusCreated)
		responseMessage := "Пользователь успешно добавлен, ID: " + strconv.Itoa(id)
		w.Write([]byte(responseMessage))

	}
}

func handlerLogin(useCaseProvider *providers.UseCaseProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Неправильный метод", http.StatusMethodNotAllowed)
			return
		}

		var loginData entities.Login
		err := json.NewDecoder(r.Body).Decode(&loginData)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, "Ошибка декодирования", http.StatusBadRequest)
			return
		}

		user, err := useCaseProvider.UserUseCase.Login(loginData)
		if err != nil {
			http.Error(w, "Ошибка аутентификации", http.StatusBadRequest)
			return
		}

		// Успешная аутентификация
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	}
}

func handlerGetUSerById(useCaseProvider *providers.UseCaseProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Неправильный метод", http.StatusMethodNotAllowed)
			return
		}

		idStr := chi.URLParam(r, "id")

		userID, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Неверный формат ID", http.StatusBadRequest)
			return
		}

		username, err := useCaseProvider.UserUseCase.GetUserById(userID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка получения пользователя: %v", err), http.StatusInternalServerError)
			return
		}

		if username == "" {
			http.Error(w, fmt.Sprintf("Пользователь не найден: %v", err), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"username": username})

	}
}
