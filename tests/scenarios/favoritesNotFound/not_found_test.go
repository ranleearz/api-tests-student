package favoritesNotFound

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/constants/messages"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/utils"

	base "api-tests-template/tests"
)

type TestSuite struct {
	suite.Suite
	loginData models.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Регистрация нового пользователя")
	s.loginData = auth.RegisterRandomUser(s.T())
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

// Negative: Проверяем 404 ошибку при добавлении в избранное несуществующего объявления
func (s *TestSuite) TestAddNonExistentAdvertisementToFavorites() {
	advertisementId := utils.RandomUUID()

	base.Step(fmt.Sprintf("Добавляем в избранное несуществующее объявление %s", advertisementId))
	errorResponse := favorites.AddFavoriteWithError(s.T(), s.loginData.Token, advertisementId, http.StatusNotFound)

	assertions.ErrorResponseEquals(s.T(), errorResponse, messages.ErrorNotFound, messages.MessageAdvertisementNotFound)
	base.Check(fmt.Sprintf("Получена ошибка 404: %s", messages.MessageAdvertisementNotFound))
}

// Negative: Проверяем 404 ошибку при удалении из избранного несуществующего объявления
func (s *TestSuite) TestRemoveNonExistentAdvertisementFromFavorites() {
	advertisementId := utils.RandomUUID()

	base.Step(fmt.Sprintf("Удаляем из избранного несуществующее объявление %s", advertisementId))
	errorResponse := favorites.RemoveFavoriteWithError(s.T(), s.loginData.Token, advertisementId, http.StatusNotFound)

	assertions.ErrorResponseEquals(s.T(), errorResponse, messages.ErrorNotFound, messages.MessageFavoriteNotFound)
	base.Check(fmt.Sprintf("Получена ошибка 404: %s", messages.MessageFavoriteNotFound))
}
