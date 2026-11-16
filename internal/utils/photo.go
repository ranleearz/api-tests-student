package utils

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	photoSize            = 10
	photoFileNamePattern = "photo-*.png"
)

// CreateTempPhoto создаёт PNG-изображение во временной директории теста и возвращает путь к нему
func CreateTempPhoto(t *testing.T) string {
	t.Helper()

	img := image.NewRGBA(image.Rect(0, 0, photoSize, photoSize))
	for x := range photoSize {
		for y := range photoSize {
			img.Set(x, y, color.RGBA{R: 255, A: 255})
		}
	}

	file, err := os.CreateTemp(t.TempDir(), photoFileNamePattern)
	require.NoError(t, err, "Не удалось создать файл фотографии")
	defer file.Close()

	require.NoError(t, png.Encode(file, img), "Не удалось записать фотографию")

	return file.Name()
}

// CreateTempPhotos создаёт указанное количество PNG-изображений и возвращает пути к ним
func CreateTempPhotos(t *testing.T, count int) []string {
	t.Helper()

	photoPaths := make([]string, 0, count)
	for range count {
		photoPaths = append(photoPaths, CreateTempPhoto(t))
	}

	return photoPaths
}
