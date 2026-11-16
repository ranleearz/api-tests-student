package advertisementCreateValidation

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/constants/messages"
	"api-tests-template/internal/managers/advertisements"
	advertisementsModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"

	base "api-tests-template/tests"
)

const (
	maxTitleLength       = 50
	maxDescriptionLength = 500
	maxPhotosCount       = 3
	validPhotosCount     = 1
	negativeValue        = "-1"
	zeroValue            = "0"
	notNumberValue       = "abc"
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

	base.Precondition("Авторизация пользователя с кредами из переменных окружения и получение его параметров")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

// Negative: Проверяем 400 ошибку при создании объявления без обязательных или с некорректными параметрами
func (s *TestSuite) TestCreateAdvertisementInvalidParams() {
	cases := []struct {
		name            string
		modifyFields    func(fields map[string]string)
		photosCount     int
		expectedMessage string
	}{
		{
			name:            "без названия",
			modifyFields:    func(fields map[string]string) { delete(fields, advertisementsModels.FieldTitle) },
			photosCount:     validPhotosCount,
			expectedMessage: messages.MessageTitleRequired,
		},
		{
			name:            "пустое название",
			modifyFields:    func(fields map[string]string) { fields[advertisementsModels.FieldTitle] = "" },
			photosCount:     validPhotosCount,
			expectedMessage: messages.MessageTitleRequired,
		},
		{
			name: "слишком длинное название",
			modifyFields: func(fields map[string]string) {
				fields[advertisementsModels.FieldTitle] = utils.RandomString(maxTitleLength + 1)
			},
			photosCount:     validPhotosCount,
			expectedMessage: fmt.Sprintf(messages.MessageTitleTooLong, maxTitleLength, maxTitleLength+1),
		},
		{
			name:            "без описания",
			modifyFields:    func(fields map[string]string) { delete(fields, advertisementsModels.FieldDescription) },
			photosCount:     validPhotosCount,
			expectedMessage: messages.MessageDescriptionRequired,
		},
		{
			name: "слишком длинное описание",
			modifyFields: func(fields map[string]string) {
				fields[advertisementsModels.FieldDescription] = utils.RandomString(maxDescriptionLength + 1)
			},
			photosCount:     validPhotosCount,
			expectedMessage: fmt.Sprintf(messages.MessageDescriptionLength, maxDescriptionLength),
		},
		{
			name:            "отрицательная цена",
			modifyFields:    func(fields map[string]string) { fields[advertisementsModels.FieldPrice] = negativeValue },
			photosCount:     validPhotosCount,
			expectedMessage: messages.MessagePriceRange,
		},
		{
			name:            "нечисловая цена",
			modifyFields:    func(fields map[string]string) { fields[advertisementsModels.FieldPrice] = notNumberValue },
			photosCount:     validPhotosCount,
			expectedMessage: fmt.Sprintf(messages.MessagePriceNotNumber, notNumberValue),
		},
		{
			name:            "нулевое количество",
			modifyFields:    func(fields map[string]string) { fields[advertisementsModels.FieldQuantity] = zeroValue },
			photosCount:     validPhotosCount,
			expectedMessage: messages.MessageQuantityRange,
		},
		{
			name:            "отрицательное количество",
			modifyFields:    func(fields map[string]string) { fields[advertisementsModels.FieldQuantity] = negativeValue },
			photosCount:     validPhotosCount,
			expectedMessage: messages.MessageQuantityRange,
		},
		{
			name:            "без фотографий",
			photosCount:     0,
			expectedMessage: messages.MessagePhotoRequired,
		},
		{
			name:            "слишком много фотографий",
			photosCount:     maxPhotosCount + 1,
			expectedMessage: fmt.Sprintf(messages.MessageTooManyPhotos, maxPhotosCount),
		},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			fields := advertisements.NewRandomAdvertisementFields()
			if tc.modifyFields != nil {
				tc.modifyFields(fields)
			}
			photoPaths := utils.CreateTempPhotos(s.T(), tc.photosCount)

			base.Step(fmt.Sprintf("Создаём объявление: %s", tc.name))
			errorResponse := advertisements.CreateAdvertisementWithError(s.T(), s.loginData.Token, fields, photoPaths,
				http.StatusBadRequest)

			assertions.ErrorResponseEquals(s.T(), errorResponse, messages.ErrorInvalidRequest, tc.expectedMessage)
			base.Check(fmt.Sprintf("Получена ошибка 400: %s", tc.expectedMessage))
		})
	}
}
