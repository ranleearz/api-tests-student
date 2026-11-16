package tests

import (
	"log"

	"api-tests-template/internal/utils"
)

// SetupSuite действия перед запуском сьюта с тестами
func SetupSuite() {
	log.Println("Init environment variables")
	utils.LoadEnv()
}

// TearDownSuite действия после запуска сьюта с тестами
func TearDownSuite() {
	log.Println("Tear down suite")
}

// Precondition выводит текст в консоль с меткой [Precondition]
func Precondition(text string) {
	utils.LogWithLabelAndTimeStamp("Precondition", text)
}

// Step выводит текст шага теста в консоль с меткой [Step]
func Step(text string) {
	utils.LogWithLabelAndTimeStamp("Step", text)
}

// Check выводит результат успешной проверки в консоль с меткой [Check]
func Check(text string) {
	utils.LogWithLabelAndTimeStamp("Check", text)
}
