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
	"github.com/alldroll/snake-game/pkg/keyboard"
)

const (
	width  = 15
	height = 15
	frame  = 150
)

var mapping = map[keyboard.Button]game.Direction{
	keyboard.UpArrow:    game.Up,
	keyboard.DownArrow:  game.Down,
	keyboard.RightArrow: game.Right,
	keyboard.LeftArrow:  game.Left,
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

func printScreen(snakeGame *game.SnakeGame, gameStart time.Time, err error) {
	fmt.Print("\x1b[2J")
	fmt.Print("\x1b[1;1f")
	fmt.Printf("Score: %d \t Time: %s\n\n", snakeGame.Score(), time.Since(gameStart).Truncate(time.Second))

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

func sleep(snakeGame *game.SnakeGame) {
	time.Sleep(frame*time.Millisecond - time.Duration(snakeGame.Score())*time.Millisecond)
}

func main() {
	eventCh := make(chan keyboard.Button)
	defer close(eventCh)

	keyboard.OnKeyDown(os.Stdin, eventCh)
	snakeGame := game.New(width, height)
	gameOver := false
	gameStart := time.Now()

	initialise()
	defer cleanup()

	for !gameOver {
		select {
		case data := <-eventCh:
			if data == keyboard.Esc {
				break
			}

			if direction, ok := mapping[data]; ok {
				snakeGame.SetDirection(direction)
			}
		default:
		}

		err := snakeGame.Move()
		printScreen(snakeGame, gameStart, err)
		gameOver = err != nil
		sleep(snakeGame)
	}
}
