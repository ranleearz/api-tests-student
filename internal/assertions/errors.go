package assertions

import (
	"testing"

	"github.com/stretchr/testify/require"

	commonModels "api-tests-template/internal/managers/common/models"
)

// ErrorResponseEquals проверяет, что тело ошибки содержит ожидаемые значения error и message
func ErrorResponseEquals(t *testing.T, actual commonModels.ErrorResponse, expectedError, expectedMessage string) {
	require.Equal(t, expectedError, actual.Error, "Некорректное значение error параметра")
	require.Equal(t, expectedMessage, actual.Message, "Некорректное значение message параметра")
}
