package daysteps

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	spentCalories "github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	parsedString := strings.Split(data, ",")
	if len(parsedString) != 2 {
		log.Printf("строка %s имеет некорректный формат", data)
		return 0, time.Duration(0), fmt.Errorf("строка %s имеет некорректный формат", data)
	}
	stepsCount, err := strconv.Atoi(parsedString[0])
	if err != nil {
		return 0, time.Duration(0), err
	}
	if stepsCount <= 0 {
		log.Printf("Количество шагов меньше или равно 0, stepsCount = %d", stepsCount)
		return 0, time.Duration(0), fmt.Errorf("количество шагов меньше или равно 0, stepsCount = %d", stepsCount)
	}
	walkDuration, err := time.ParseDuration(parsedString[1])
	if err != nil {
		return 0, time.Duration(0), err
	}
	if walkDuration <= 0 {
		log.Printf("Значение времени меньше или равно 0, walkDuration = %.f", walkDuration.Hours())
		return 0, time.Duration(0), fmt.Errorf("значение времени меньше или равно 0, walkDuration = %.2f", walkDuration.Hours())
	}
	return stepsCount, walkDuration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	stepsCount, walkDuration, err := parsePackage(data)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	walkDistanceInM := float32(stepsCount) * stepLength
	walkDistanceIKm := walkDistanceInM / mInKm
	walkSpentCalories, err := spentCalories.WalkingSpentCalories(stepsCount, weight, height, walkDuration)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", stepsCount, walkDistanceIKm, walkSpentCalories)
}
