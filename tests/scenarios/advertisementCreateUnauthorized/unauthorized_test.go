package advertisementCreateUnauthorized

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/constants/messages"
	"api-tests-template/internal/managers/advertisements"
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

// Negative: Проверяем 401 ошибку при создании объявления без токена или с неверным токеном
func (s *TestSuite) TestCreateAdvertisementUnauthorized() {
	cases := []struct {
		name            string
		token           string
		expectedMessage string
	}{
		{"без токена", "", messages.MessageAuthorizationRequired},
		{"с неверным токеном", messages.InvalidToken, messages.MessageInvalidAuthorization},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			fields := advertisements.NewRandomAdvertisementFields()
			photoPaths := []string{utils.CreateTempPhoto(s.T())}

			base.Step(fmt.Sprintf("Создаём объявление %s", tc.name))
			errorResponse := advertisements.CreateAdvertisementWithError(s.T(), tc.token, fields, photoPaths,
				http.StatusUnauthorized)

			assertions.ErrorResponseEquals(s.T(), errorResponse, messages.ErrorUnauthorized, tc.expectedMessage)
			base.Check(fmt.Sprintf("Получена ошибка 401: %s", tc.expectedMessage))
		})
	}
}
