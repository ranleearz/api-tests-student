package advertisements

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/helpers/response"
	"api-tests-template/internal/managers/advertisements/models"
	commonModels "api-tests-template/internal/managers/common/models"
	"api-tests-template/internal/utils"
)

const (
	randomTitlePrefix    = "autotest "
	randomTitleLength    = 12
	defaultDescription   = "Объявление, созданное автотестом"
	defaultPrice         = 1000
	defaultQuantity      = 1
	priceFormatPrecision = -1
	priceFormatBitSize   = 64
)

func GetAdvertisementsList(t *testing.T, token, text string) models.AdvertisementsList {
	advertisementsRawResponse := advertisements.HttpGetAdvertisements(t, token, text)
	defer advertisementsRawResponse.Body.Close()

	respArray, errResponse := io.ReadAll(advertisementsRawResponse.Body)
	require.NoError(t, errResponse)

	advertisementsResponse := models.AdvertisementsList{}
	err := json.Unmarshal(respArray, &advertisementsResponse)
	require.NoError(t, err)

	return advertisementsResponse
}

// RandomTitle генерирует уникальное название объявления
func RandomTitle() string {
	return randomTitlePrefix + utils.RandomString(randomTitleLength)
}

// NewAdvertisementFields формирует поля формы создания объявления
func NewAdvertisementFields(title, description string, price float64, quantity int) map[string]string {
	return map[string]string{
		models.FieldTitle:       title,
		models.FieldDescription: description,
		models.FieldPrice:       strconv.FormatFloat(price, 'f', priceFormatPrecision, priceFormatBitSize),
		models.FieldQuantity:    strconv.Itoa(quantity),
	}
}

// NewRandomAdvertisementFields формирует валидные поля формы создания объявления со случайным названием
func NewRandomAdvertisementFields() map[string]string {
	return NewAdvertisementFields(RandomTitle(), defaultDescription, defaultPrice, defaultQuantity)
}

// CreateAdvertisement создаёт объявление с переданными полями и фотографиями, ожидает 201
func CreateAdvertisement(t *testing.T, token string, fields map[string]string, photoPaths []string) models.Advertisement {
	advertisementResponse := models.Advertisement{}
	response.ReadJSON(t, advertisements.HttpPostAdvertisement(t, token, fields, photoPaths), http.StatusCreated, &advertisementResponse)
	require.NotEmpty(t, advertisementResponse.Id, "id созданного объявления пустой")

	return advertisementResponse
}

// CreateAdvertisementWithError отправляет запрос на создание объявления и возвращает тело ошибки, ожидает expectedStatusCode
func CreateAdvertisementWithError(t *testing.T, token string, fields map[string]string, photoPaths []string,
	expectedStatusCode int) commonModels.ErrorResponse {
	errorResponse := commonModels.ErrorResponse{}
	response.ReadJSON(t, advertisements.HttpPostAdvertisement(t, token, fields, photoPaths), expectedStatusCode, &errorResponse)

	return errorResponse
}

// CreateRandomAdvertisement создаёт валидное объявление со случайным названием и одной фотографией и удаляет его после теста
func CreateRandomAdvertisement(t *testing.T, token string) models.Advertisement {
	advertisement := CreateAdvertisement(t, token, NewRandomAdvertisementFields(), []string{utils.CreateTempPhoto(t)})
	t.Cleanup(func() {
		DeleteAdvertisementIfExists(t, token, advertisement.Id)
	})

	return advertisement
}

// GetAdvertisement возвращает объявление по id, ожидает 200
func GetAdvertisement(t *testing.T, id string) models.Advertisement {
	advertisementResponse := models.Advertisement{}
	response.ReadJSON(t, advertisements.HttpGetAdvertisement(t, id), http.StatusOK, &advertisementResponse)

	return advertisementResponse
}

// GetAdvertisementPhotos возвращает список фотографий объявления, ожидает 200
func GetAdvertisementPhotos(t *testing.T, token, id string) []models.Photo {
	var photos []models.Photo
	response.ReadJSON(t, advertisements.HttpGetAdvertisementPhotos(t, token, id), http.StatusOK, &photos)

	return photos
}

// DeleteAdvertisement удаляет объявление по id, ожидает 204
func DeleteAdvertisement(t *testing.T, token, id string) {
	response.ReadJSON(t, advertisements.HttpDeleteAdvertisement(t, token, id), http.StatusNoContent, nil)
}

// DeleteAdvertisementIfExists удаляет объявление по id, если оно ещё не было удалено
func DeleteAdvertisementIfExists(t *testing.T, token, id string) {
	getRawResponse := advertisements.HttpGetAdvertisement(t, id)
	require.NoError(t, getRawResponse.Body.Close())
	if getRawResponse.StatusCode == http.StatusNotFound {
		return
	}
	DeleteAdvertisement(t, token, id)
}
