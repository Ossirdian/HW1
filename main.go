package main

import (
	"errors"
	"fmt"
	"math/rand/v2"
)

var (
	errInvalidInput = errors.New("invalid input, please provide a number between 1 and 100")
	errOutOfRange   = errors.New("number out of range, please provide a number between 1 and 100")
	errInvalidMode  = errors.New("invalid mode, please choose from easy, medium, or hard")
)

var (
	EasyMode   = "easy"
	MediumMode = "medium"
	HardMode   = "hard"
)

type Game struct {
	secretNumber   int
	attemptsNumber int
}

func initGame(mode string) (game Game, err error) {
	switch mode {
	case EasyMode:
		attempts := 15
		randomNumber := generateRandomNumber(50)
		return Game{secretNumber: randomNumber, attemptsNumber: attempts}, nil
	case MediumMode:
		attempts := 10
		randomNumber := generateRandomNumber(100)
		return Game{secretNumber: randomNumber, attemptsNumber: attempts}, nil
	case HardMode:
		attempts := 5
		randomNumber := generateRandomNumber(200)
		return Game{secretNumber: randomNumber, attemptsNumber: attempts}, nil
	default:
		return Game{}, errInvalidMode
	}
}

func main() {
	fmt.Println("Добро пожаловать в игру 'Угадай число'!")
	fmt.Println("Выберите режим игры: easy(15 попыток, 1-50), medium(10 попыток, 1-100), hard(5 попыток, 1-200)")
	fmt.Print("Введите режим игры: ")
	var mode string
	fmt.Scan(&mode)
	game, err := initGame(mode)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := playGame(game.secretNumber, game.attemptsNumber); err != nil {
		fmt.Println(err)
	}

}

func generateRandomNumber(number int) int {
	return rand.IntN(number) + 1
}

func validateInput(input int) error {
	if input < 1 || input > 100 {
		return errOutOfRange
	}
	return nil
}

func getUserInput() (int, error) {
	var input int
	fmt.Print("Enter a number between 1 and 100: ")
	_, err := fmt.Scan(&input)
	if err != nil {
		return 0, errInvalidInput
	}
	if err := validateInput(input); err != nil {
		return 0, err
	}
	return input, nil
}

func playGame(secretNumber, attemptNumber int) error {
	var guess int
	var attempts []int
	for guess != secretNumber && len(attempts) < attemptNumber {
		var err error
		guess, err = getUserInput()
		if err != nil {
			return err
		}
		attempts = append(attempts, guess)
	}
	fmt.Printf("Вы угадали!🙌")
	return nil
}
