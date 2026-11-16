package response

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// ReadJSON проверяет статус-код ответа и десериализует JSON-тело в target, если он передан
func ReadJSON(t *testing.T, resp *http.Response, expectedStatusCode int, target any) {
	t.Helper()
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "Не удалось прочитать тело ответа")

	require.Equalf(t, expectedStatusCode, resp.StatusCode,
		"HTTP status code должен быть %d, тело ответа: %s", expectedStatusCode, string(body))

	if target == nil {
		return
	}

	require.NoErrorf(t, json.Unmarshal(body, target), "Не удалось распарсить тело ответа: %s", string(body))
}
