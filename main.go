package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"strconv"
	"time"
)

const (
	resultsFile = "results.json"
)

const (
	easyMode   = "easy"
	mediumMode = "medium"
	hardMode   = "hard"
)

const (
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorReset  = "\033[0m"
)

var (
	errInvalidMode  = errors.New("invalid mode, please choose from easy, medium, or hard")
	errInvalidInput = errors.New("invalid input, please enter a valid number")
)

type GameMode struct {
	name     string
	attempts int
	min      int
	max      int
}

type Game struct {
	attempts       []int
	startedAt      time.Time
	secretNumber   int
	DifficultyMode GameMode
}

type GameStats struct {
	Attempts       int           `json:"attempts"`
	StartedAt      time.Time     `json:"started_at"`
	TimeElapsed    time.Duration `json:"time_elapsed"`
	Winner         bool          `json:"winner"`
	DifficultyMode string        `json:"difficulty_mode"`
}

func initGame(mode string) (game *Game, err error) {
	playerMode := mode
	switch playerMode {
	case easyMode:
		randomNumber := generateRandomNumber(50)
		return &Game{secretNumber: randomNumber, startedAt: time.Now(), DifficultyMode: GameMode{name: easyMode, attempts: 15, min: 1, max: 50}}, nil
	case mediumMode:
		randomNumber := generateRandomNumber(100)
		return &Game{secretNumber: randomNumber, startedAt: time.Now(), DifficultyMode: GameMode{name: mediumMode, attempts: 10, min: 1, max: 100}}, nil
	case hardMode:
		randomNumber := generateRandomNumber(200)
		return &Game{secretNumber: randomNumber, startedAt: time.Now(), DifficultyMode: GameMode{name: hardMode, attempts: 5, min: 1, max: 200}}, nil
	default:
		return &Game{}, errInvalidMode
	}
}

func main() {
	fmt.Println("Добро пожаловать в игру 'Угадай число'!")
	fmt.Println("Выберите режим игры: easy(15 попыток, 1-50), medium(10 попыток, 1-100), hard(5 попыток, 1-200)")
	fmt.Print("Введите режим игры: ")
	var mode string
	for {
		fmt.Scan(&mode)
		playerMode := mode
		if playerMode == easyMode || playerMode == mediumMode || playerMode == hardMode {
			break
		}
		fmt.Println("Неверный режим игры. Пожалуйста, выберите из easy, medium или hard.")
	}
	game, err := initGame(mode)
	if err != nil {
		fmt.Println(err)
		return
	}
	game.playGame(game.secretNumber, game.DifficultyMode.attempts)
	for {
		fmt.Print("Сыграть еще раз (1) Посмотреть статистику (2) Выйти (3): ")
		var choice int
		fmt.Scan(&choice)
		switch choice {
		case 1:
			game, err = initGame(mode)
			if err != nil {
				fmt.Println(err)
				return
			}
			game.playGame(game.secretNumber, game.DifficultyMode.attempts)
		case 2:
			stats, err := LoadGameStats()
			if err != nil {
				fmt.Println("Ошибка при загрузке статистики:", err)
				continue
			}
			if len(stats) == 0 {
				fmt.Println("Сохраненных игр пока нет.")
				continue
			}
			fmt.Println("Список всех игр:")
			for i, stat := range stats {
				result := "проигрыш"
				if stat.Winner {
					result = "победа"
				}
				fmt.Printf("%d. Режим: %s, результат: %s, попыток: %d, время: %s, длительность: %s\n",
					i+1, stat.DifficultyMode, result, stat.Attempts,
					stat.StartedAt.Format(time.RFC3339), stat.TimeElapsed)
			}
		case 3:
			fmt.Println("Спасибо за игру! До свидания!")
			return
		}
	}
}

func generateRandomNumber(number int) int {
	return rand.IntN(number) + 1
}

func validateInput(input, min, max int) error {
	if input < min || input > max {
		return fmt.Errorf("Число должно быть между %d и %d", min, max)
	}
	return nil
}

func (g Game) getUserInput() (int, error) {
	var input string
	fmt.Printf("Введите число между %d и %d: \n", g.DifficultyMode.min, g.DifficultyMode.max)
	_, err := fmt.Scan(&input)
	if err != nil {
		return 0, errInvalidInput
	}
	inputNum, err := strconv.Atoi(input)
	if err != nil {
		return 0, errInvalidInput
	}
	if err := validateInput(inputNum, g.DifficultyMode.min, g.DifficultyMode.max); err != nil {
		return 0, err
	}
	fmt.Printf(ColorYellow+"Вы ввели: %d\n"+ColorReset, inputNum)
	return inputNum, nil
}

func (g *Game) playGame(secretNumber, attemptNumber int) {
	var attempts []int
	var winner bool
	for !winner && len(attempts) < attemptNumber {
		fmt.Printf(ColorYellow+"Попытка %d из %d\n"+ColorReset, len(attempts)+1, attemptNumber)
		guess, err := g.getUserInput()
		if err != nil {
			fmt.Println(err)
		} else {
			attempts = append(attempts, guess)
			winner = Hint(secretNumber, guess)
		}
		fmt.Printf("Осталось попыток: %d\n", attemptNumber-len(attempts))
	}
	if winner {
		fmt.Println(ColorGreen + "Поздравляем! Вы выиграли!" + ColorReset)
	} else {
		fmt.Println(ColorRed + "Вы исчерпали все попытки. Игра окончена." + ColorReset)
	}
	fmt.Printf(ColorYellow+"Загаданное число было: %d\n"+ColorReset, secretNumber)
	for i, v := range attempts {
		fmt.Printf(ColorYellow+"Попытка %d: %d\n"+ColorReset, i+1, v)
	}
	g.attempts = attempts
	SaveGameStats(g.getGameStats(winner))
}

func Hint(secretNumber, guess int) bool {
	if guess == secretNumber {
		return true
	} else {
		difference := secretNumber - guess
		absDifference := math.Abs(float64(difference))
		if int(absDifference) <= 5 {
			fmt.Print("Горячо! 🔥  ")
		} else if int(absDifference) <= 15 {
			fmt.Print("Тепло! 🌟  ")
		} else {
			fmt.Print("Холодно! ❄️  ")
		}
		if difference > 0 {
			fmt.Println("Загаданное число больше вашего числа.")
		} else if difference < 0 {
			fmt.Println("Загаданное число меньше вашего числа.")
		}
		return false
	}
}

func (g *Game) getGameStats(winner bool) GameStats {
	return GameStats{
		Attempts:       len(g.attempts),
		StartedAt:      g.startedAt,
		TimeElapsed:    time.Since(g.startedAt),
		Winner:         winner,
		DifficultyMode: g.DifficultyMode.name,
	}
}

func SaveGameStats(stats GameStats) {
	allStats, err := LoadGameStats()
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("Ошибка при загрузке статистики игры:", err)
		return
	}
	allStats = append(allStats, stats)
	jsonData, err := json.Marshal(allStats)
	if err != nil {
		fmt.Println("Ошибка при сохранении статистики игры:", err)
		return
	}
	err = os.WriteFile(resultsFile, jsonData, 0644)
	if err != nil {
		fmt.Println("Ошибка при записи статистики игры в файл:", err)
		return
	}
	fmt.Println("Статистика игры успешно сохранена в файл results.json")
}

func LoadGameStats() ([]GameStats, error) {
	stats := make([]GameStats, 0)
	jsonData, err := os.ReadFile(resultsFile)
	if err != nil {
		return stats, err
	}
	err = json.Unmarshal(jsonData, &stats)
	if err != nil {
		return stats, err
	}
	return stats, nil
}
