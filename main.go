package main

// Inspired by https://github.com/danicat/pacgo/

import (
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

func printScreen(game *game.SnakeGame) {
	fmt.Print("\x1b[2J")
	fmt.Print("\x1b[1;1f")

	for _, line := range game.Grid() {
		for _, entity := range line {
			fmt.Printf("%c", entity)
		}

		fmt.Println()
	}
}

func main() {
	snakeGame := game.New(width, height)
	direction := game.Right
	input := make(chan string)

	initialise()
	defer cleanup()

	for {
		printScreen(snakeGame)

		go func(input chan string) {
			data, err := readInput()

			if err != nil {
				log.Fatal("error reading input:", err)
			}

			input <- data
		}(input)

		select {
		case data := <-input:
			if data == "ESC" {
				break
			}

			if d, ok := mapping[data]; ok {
				direction = d
			}
		default: //nothing
		}

		snakeGame.Move(direction)
		time.Sleep(200 * time.Millisecond)
	}
}
