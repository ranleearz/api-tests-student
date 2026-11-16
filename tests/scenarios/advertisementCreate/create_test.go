package advertisementCreate

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"

	base "api-tests-template/tests"
)

const (
	advertisementDescription = "Объявление со всеми параметрами, созданное автотестом"
	advertisementPrice       = 1500.5
	advertisementQuantity    = 3
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

// Positive: Проверяем создание объявления со всеми возможными параметрами
func (s *TestSuite) TestCreateAdvertisementWithAllParams() {
	title := advertisements.RandomTitle()
	fields := advertisements.NewAdvertisementFields(title, advertisementDescription, advertisementPrice, advertisementQuantity)
	photoPath := utils.CreateTempPhoto(s.T())

	base.Step(fmt.Sprintf("Создаём объявление %q со всеми параметрами и фотографией", title))
	created := advertisements.CreateAdvertisement(s.T(), s.loginData.Token, fields, []string{photoPath})
	s.T().Cleanup(func() {
		base.Step(fmt.Sprintf("Удаляем созданное объявление %s", created.Id))
		advertisements.DeleteAdvertisementIfExists(s.T(), s.loginData.Token, created.Id)
	})
	require.NotEmpty(s.T(), created.Photos, "В ответе на создание нет фотографий")

	s.Run("Объявление доступно по id и содержит переданные параметры", func() {
		base.Step(fmt.Sprintf("Получаем объявление по id %s", created.Id))
		advertisement := advertisements.GetAdvertisement(s.T(), created.Id)

		require.Equal(s.T(), created.Id, advertisement.Id, "id не совпадает")
		require.Equal(s.T(), title, advertisement.Title, "Название не совпадает")
		require.Equal(s.T(), advertisementDescription, advertisement.Description, "Описание не совпадает")
		require.InDelta(s.T(), advertisementPrice, advertisement.Price, 0, "Цена не совпадает")
		require.Equal(s.T(), advertisementQuantity, advertisement.Quantity, "Количество не совпадает")
		require.Equal(s.T(), s.loginData.User.Id, advertisement.UserId, "user_id не совпадает")
		base.Check("Параметры объявления совпадают с переданными при создании")
	})

	s.Run("Объявление находится в поиске по полному названию", func() {
		base.Step(fmt.Sprintf("Ищем объявления по названию %q", title))
		list := advertisements.GetAdvertisementsList(s.T(), s.loginData.Token, title)

		assertions.ContainsAdvertisementByTitle(s.T(), list.Items, title)
		base.Check("Объявление найдено в поиске")
	})

	s.Run("Фотография загружена и доступна", func() {
		base.Step(fmt.Sprintf("Получаем фотографии объявления %s", created.Id))
		photos := advertisements.GetAdvertisementPhotos(s.T(), s.loginData.Token, created.Id)

		require.Len(s.T(), photos, len(created.Photos), "Количество фотографий не совпадает с загруженным")
		require.Equal(s.T(), created.Photos[0].Id, photos[0].Id, "id фотографии не совпадает")
		require.NotEmpty(s.T(), photos[0].Url, "url фотографии пустой")
		base.Check("Фотография объявления доступна")
	})
}
