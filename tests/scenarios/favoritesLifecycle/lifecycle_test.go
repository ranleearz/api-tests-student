package favoritesLifecycle

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/managers/favorites"

	base "api-tests-template/tests"
)

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

// Positive: Проверяем полный цикл работы с избранным авторизованного пользователя
func (s *TestSuite) TestFavoritesLifecycle() {
	base.Precondition("Владелец создаёт объявление")
	advertisement := advertisements.CreateRandomAdvertisement(s.T(), s.ownerLoginData.Token)

	s.Run("Добавленное объявление появляется в избранном", func() {
		base.Step(fmt.Sprintf("Добавляем объявление %s в избранное", advertisement.Id))
		favorites.AddFavorite(s.T(), s.favoriterLoginData.Token, advertisement.Id)

		base.Step("Получаем список избранного")
		favoritesList := favorites.GetFavorites(s.T(), s.favoriterLoginData.Token)

		assertions.ContainsFavoriteByAdvertisementId(s.T(), favoritesList.Items, advertisement.Id)
		base.Check("Объявление присутствует в избранном")
	})

	s.Run("Удалённое из избранного объявление пропадает из списка", func() {
		base.Step(fmt.Sprintf("Удаляем объявление %s из избранного", advertisement.Id))
		favorites.RemoveFavorite(s.T(), s.favoriterLoginData.Token, advertisement.Id)

		base.Step("Получаем список избранного")
		favoritesList := favorites.GetFavorites(s.T(), s.favoriterLoginData.Token)

		assertions.NotContainsFavoriteByAdvertisementId(s.T(), favoritesList.Items, advertisement.Id)
		base.Check("Объявление отсутствует в избранном")
	})

	s.Run("Удалённое владельцем объявление пропадает из избранного", func() {
		base.Step(fmt.Sprintf("Повторно добавляем объявление %s в избранное", advertisement.Id))
		favorites.AddFavorite(s.T(), s.favoriterLoginData.Token, advertisement.Id)

		base.Step(fmt.Sprintf("Владелец удаляет объявление %s", advertisement.Id))
		advertisements.DeleteAdvertisement(s.T(), s.ownerLoginData.Token, advertisement.Id)

		base.Step("Получаем список избранного")
		favoritesList := favorites.GetFavorites(s.T(), s.favoriterLoginData.Token)

		assertions.NotContainsFavoriteByAdvertisementId(s.T(), favoritesList.Items, advertisement.Id)
		base.Check("Удалённое объявление отсутствует в избранном")
	})
}
