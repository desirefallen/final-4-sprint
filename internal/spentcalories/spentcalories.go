package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parsedString := strings.Split(data, ",")
	if len(parsedString) != 3 {
		return 0, "", time.Duration(0), fmt.Errorf("строка %s имеет некорректный формат", data)
	}
	steps, err := strconv.Atoi(parsedString[0])
	if err != nil {
		return 0, "", time.Duration(0), err
	}
	if steps <= 0 {
		log.Printf("Количество шагов меньше или равно 0, steps = %d", steps)
		return 0, "", time.Duration(0), fmt.Errorf("количество шагов меньше или равно 0, steps = %d", steps)
	}
	activityType := parsedString[1]
	duration, err := time.ParseDuration(parsedString[2])
	if duration <= 0 {
		log.Printf("Значение времени меньше или равно 0, duration = %.f", duration.Hours())
		return 0, "", time.Duration(0), fmt.Errorf("значение времени меньше или равно 0, duration = %.2f", duration.Hours())
	}
	if err != nil {
		return 0, "", time.Duration(0), err
	}
	return steps, activityType, duration, nil
}

func distance(steps int, height float64) float64 {
	if steps <= 0 {
		log.Printf("Количество шагов меньше или равно 0, steps = %d", steps)
		return 0
	}
	if height <= 0 {
		log.Printf("Рост меньше или равен 0, height = %.2f", height)
		return 0
	}
	return (height * stepLengthCoefficient * float64(steps) / mInKm)
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	return (distance(steps, height) / (duration.Hours()))
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activityType, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
	}
	switch activityType {
	case "Ходьба":
		distance := distance(steps, height)
		meanSpeed := meanSpeed(steps, height, duration)
		spentCalories, err := WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", "Ходьба", duration.Hours(), distance, meanSpeed, spentCalories), nil
	case "Бег":
		distance := distance(steps, height)
		meanSpeed := meanSpeed(steps, height, duration)
		spentCalories, err := RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", "Бег", duration.Hours(), distance, meanSpeed, spentCalories), nil
	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("параметр steps некорректен")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("параметр weight некорректен")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("параметр duration некорректен")
	}
	if height <= 0 {
		return 0, fmt.Errorf("параметр height некорректен")
	}
	durationInMinutes := duration.Minutes()
	return ((weight * meanSpeed(steps, height, duration) * durationInMinutes) / minInH), nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("параметр steps некорректен")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("параметр weight некорректен")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("параметр duration некорректен")
	}
	if height <= 0 {
		return 0, fmt.Errorf("параметр height некорректен")
	}
	durationInMinutes := duration.Minutes()
	return (((weight * meanSpeed(steps, height, duration) * durationInMinutes) / minInH) * walkingCaloriesCoefficient), nil
}
