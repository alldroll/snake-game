package main

// Inspired by https://github.com/danicat/pacgo/

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/alldroll/snake-game/pkg/game"
)

const (
	width  = 10
	height = 10
	frame  = 200
)

var mapping = map[string]game.Direction{
	"UP":    game.Up,
	"DOWN":  game.Down,
	"RIGHT": game.Right,
	"LEFT":  game.Left,
}

func initialise() {
	cbTerm := exec.Command("stty", "cbreak", "-echo")
	cbTerm.Stdin = os.Stdin

	if err := cbTerm.Run(); err != nil {
		log.Fatalln("unable to activate cbreak mode:", err)
	}
}

func cleanup() {
	cookedTerm := exec.Command("stty", "-cbreak", "echo")
	cookedTerm.Stdin = os.Stdin

	if err := cookedTerm.Run(); err != nil {
		log.Fatalln("unable to activate cooked mode:", err)
	}
}

func readInput() (string, error) {
	buffer := make([]byte, 100)
	cnt, err := os.Stdin.Read(buffer)

	if err != nil {
		return "", err
	}

	if cnt == 1 && buffer[0] == 0x1b {
		return "ESC", nil
	}

	if buffer[0] == 0x1b && buffer[1] == '[' {
		switch buffer[2] {
		case 'A':
			return "UP", nil
		case 'B':
			return "DOWN", nil
		case 'C':
			return "RIGHT", nil
		case 'D':
			return "LEFT", nil
		}
	}

	return "", nil
}

func printScreen(snakeGame *game.SnakeGame, err error) {
	fmt.Print("\x1b[2J")
	fmt.Print("\x1b[1;1f")
	fmt.Printf("Score: %d\n\n", snakeGame.Score())

	for _, line := range snakeGame.Grid() {
		for _, entity := range line {
			fmt.Printf(" %c ", entity)
		}

		fmt.Println()
	}

	if errors.Is(err, game.ErrGameIsCompleted) {
		fmt.Print("\nGame is completed\n")
	} else if err != nil {
		fmt.Printf("\nGame over: %s\n", err.Error())
	}
}

func main() {
	snakeGame := game.New(width, height)
	input := make(chan string)
	direction := game.Right
	gameOver := false

	initialise()
	defer cleanup()

	go func(input chan string) {
		for {
			data, err := readInput()

			if err != nil {
				log.Fatal("error reading input:", err)
			}

			input <- data
		}
	}(input)

	timer := time.NewTimer(frame * time.Millisecond)

	for !gameOver {
		select {
		case data := <-input:
			if data == "ESC" {
				break
			}

			if d := mapping[data]; snakeGame.CanMove(d) {
				direction = d
			}
		case <-timer.C:
			err := snakeGame.Move(direction)
			printScreen(snakeGame, err)
			gameOver = err != nil
			timer.Reset(frame*time.Millisecond - 5*time.Duration(snakeGame.Score()))
		}
	}
}
