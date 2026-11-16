package auth

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/auth"
	"api-tests-template/internal/helpers/response"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"
)

const (
	registeredUserPassword  = "Autotest123"
	registeredUserFirstName = "Autotest"
)

// Login выполняет авторизацию пользователя и возвращает токен с данными пользователя, ожидает 200
func Login(t *testing.T, username, password string) models.LoginOkResponse {
	require.NotEmpty(t, username, "Отсутствует логин пользователя")
	require.NotEmpty(t, password, "Отсутствует пароль пользователя")

	request := models.LoginRequest{
		Email:    username,
		Password: password,
	}

	loginRawResponse := auth.HttpPostAuthLogin(t, request)
	defer loginRawResponse.Body.Close()

	respArray, errResponse := io.ReadAll(loginRawResponse.Body)
	require.NoError(t, errResponse)

	loginResponse := models.LoginOkResponse{}
	err := json.Unmarshal(respArray, &loginResponse)
	require.NoError(t, err)

	return loginResponse
}

// RegisterRandomUser регистрирует нового пользователя со случайным email и возвращает его токен и данные, ожидает 201
func RegisterRandomUser(t *testing.T) models.LoginOkResponse {
	request := models.RegisterRequest{
		Email:     utils.RandomEmail(),
		Password:  registeredUserPassword,
		FirstName: registeredUserFirstName,
	}

	registerResponse := models.LoginOkResponse{}
	response.ReadJSON(t, auth.HttpPostAuthRegister(t, request), http.StatusCreated, &registerResponse)
	require.NotEmpty(t, registerResponse.Token, "Токен зарегистрированного пользователя пустой")

	return registerResponse
}
