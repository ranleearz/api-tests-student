package favoritesRestrictions

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/constants/messages"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/managers/favorites"

	base "api-tests-template/tests"
)

const invalidAdvertisementId = "not-a-uuid"

type TestSuite struct {
	suite.Suite
	ownerLoginData     models.LoginOkResponse
	favoriterLoginData models.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Авторизация владельца объявлений с кредами из переменных окружения")
	s.ownerLoginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))

	base.Precondition("Регистрация нового пользователя, который будет добавлять объявления в избранное")
	s.favoriterLoginData = auth.RegisterRandomUser(s.T())
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

// Negative: Проверяем 409 ошибку при повторном добавлении объявления в избранное
func (s *TestSuite) TestAddFavoriteTwice() {
	base.Precondition("Владелец создаёт объявление")
	advertisement := advertisements.CreateRandomAdvertisement(s.T(), s.ownerLoginData.Token)

	base.Step(fmt.Sprintf("Добавляем объявление %s в избранное", advertisement.Id))
	favorites.AddFavorite(s.T(), s.favoriterLoginData.Token, advertisement.Id)

	base.Step(fmt.Sprintf("Повторно добавляем объявление %s в избранное", advertisement.Id))
	errorResponse := favorites.AddFavoriteWithError(s.T(), s.favoriterLoginData.Token, advertisement.Id, http.StatusConflict)

	assertions.ErrorResponseEquals(s.T(), errorResponse, messages.ErrorAlreadyExists, messages.MessageFavoriteAlreadyExists)
	base.Check(fmt.Sprintf("Получена ошибка 409: %s", messages.MessageFavoriteAlreadyExists))
}

// Negative: Проверяем 403 ошибку при добавлении в избранное собственного объявления
func (s *TestSuite) TestAddOwnAdvertisementToFavorites() {
	base.Precondition("Владелец создаёт объявление")
	advertisement := advertisements.CreateRandomAdvertisement(s.T(), s.ownerLoginData.Token)

	base.Step(fmt.Sprintf("Владелец добавляет своё объявление %s в избранное", advertisement.Id))
	errorResponse := favorites.AddFavoriteWithError(s.T(), s.ownerLoginData.Token, advertisement.Id, http.StatusForbidden)

	assertions.ErrorResponseEquals(s.T(), errorResponse, messages.ErrorForbidden, messages.MessageFavoriteOwnForbidden)
	base.Check(fmt.Sprintf("Получена ошибка 403: %s", messages.MessageFavoriteOwnForbidden))
}

// Negative: Проверяем 400 ошибку при добавлении в избранное с некорректным advertisement_id
func (s *TestSuite) TestAddFavoriteInvalidAdvertisementId() {
	cases := []struct {
		name            string
		advertisementId string
	}{
		{"пустой advertisement_id", ""},
		{"advertisement_id не в формате UUID", invalidAdvertisementId},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			base.Step(fmt.Sprintf("Добавляем в избранное объявление: %s", tc.name))
			errorResponse := favorites.AddFavoriteWithError(s.T(), s.favoriterLoginData.Token, tc.advertisementId,
				http.StatusBadRequest)

			expectedMessage := fmt.Sprintf(messages.MessageInvalidUUIDLength, len(tc.advertisementId))
			assertions.ErrorResponseEquals(s.T(), errorResponse, messages.ErrorInvalidRequest, expectedMessage)
			base.Check(fmt.Sprintf("Получена ошибка 400: %s", expectedMessage))
		})
	}
}
