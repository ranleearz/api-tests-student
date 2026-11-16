package favoritesUnauthorized

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/constants/messages"
	commonModels "api-tests-template/internal/managers/common/models"
	"api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/utils"

	base "api-tests-template/tests"
)

type TestSuite struct {
	suite.Suite
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

// Negative: Проверяем 401 ошибку на ручках избранного без токена или с неверным токеном
func (s *TestSuite) TestFavoritesUnauthorized() {
	advertisementId := utils.RandomUUID()

	operations := []struct {
		name    string
		request func(token string) commonModels.ErrorResponse
	}{
		{"GET /favorites", func(token string) commonModels.ErrorResponse {
			return favorites.GetFavoritesWithError(s.T(), token, http.StatusUnauthorized)
		}},
		{"POST /favorites", func(token string) commonModels.ErrorResponse {
			return favorites.AddFavoriteWithError(s.T(), token, advertisementId, http.StatusUnauthorized)
		}},
		{"DELETE /favorites/{ad_id}", func(token string) commonModels.ErrorResponse {
			return favorites.RemoveFavoriteWithError(s.T(), token, advertisementId, http.StatusUnauthorized)
		}},
	}

	tokens := []struct {
		name            string
		token           string
		expectedMessage string
	}{
		{"без токена", "", messages.MessageAuthorizationRequired},
		{"с неверным токеном", messages.InvalidToken, messages.MessageInvalidAuthorization},
	}

	for _, operation := range operations {
		for _, tc := range tokens {
			s.Run(fmt.Sprintf("%s %s", operation.name, tc.name), func() {
				base.Step(fmt.Sprintf("Выполняем %s %s", operation.name, tc.name))
				errorResponse := operation.request(tc.token)

				assertions.ErrorResponseEquals(s.T(), errorResponse, messages.ErrorUnauthorized, tc.expectedMessage)
				base.Check(fmt.Sprintf("Получена ошибка 401: %s", tc.expectedMessage))
			})
		}
	}
}
