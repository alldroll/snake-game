package game

// Mostly inspired by https://leetcode.com/problems/design-snake-game

import (
	"errors"
	"math/rand"
	"time"
)

// Direction represents a snake next direction.
type Direction uint8

const (
	// Right tells that a snake is going to move to the right.
	Right Direction = iota
	// Down tells that a snake is going to move to the down.
	Down
	// Left tells that a snake is going to move to the left.
	Left
	// Up tells that a snake is going to move to the up.
	Up
)

// CanApply tells if the current direction can be changed to the given one.
func (d Direction) CanApply(to Direction) bool {
	if d == to {
		return true
	}

	if d == Left || d == Right {
		return to == Up || to == Down
	}

	return to == Left || to == Right
}

// Entity represents a grid entity.
type Entity byte

const (
	// Snake represents a part of a snake.
	Snake Entity = Entity(35)
	// Cell represents an empty cell of the grid.
	Cell Entity = Entity(183)
	// Food represents a food on the grid.
	Food Entity = Entity(243)
)

// ErrUnacceptableDirection tells that the given direction was not met by CanMove method.
var ErrUnacceptableDirection = errors.New("the given direction is unacceptable")

// ErrBoundriesCross means that the snake has crossed the border.
var ErrBoundriesCross = errors.New("snake has crossed the border")

// ErrSelfEating means that the snake has bitten its body.
var ErrSelfEating = errors.New("snake has bitten its body")

// ErrGameIsCompleted tells that the game was successfully completed.
var ErrGameIsCompleted = errors.New("game is completed")

var directions = [][2]int{
	{1, 0},
	{0, 1},
	{-1, 0},
	{0, -1},
}

type unit struct {
	x, y int
}

func (u unit) next(d Direction) unit {
	return unit{
		x: u.x + directions[d][0],
		y: u.y + directions[d][1],
	}
}

// SnakeGame represents a game logic.
type SnakeGame struct {
	gen       *rand.Rand
	width     int
	height    int
	queue     []unit
	table     map[unit]bool
	grid      [][]Entity
	score     int
	direction Direction
}

// New creates a new snake game state.
func New(width int, height int) *SnakeGame {
	u := unit{x: 0, y: 0}
	grid := make([][]Entity, height)

	for i := range grid {
		grid[i] = make([]Entity, width)

		for j := range grid[i] {
			grid[i][j] = Cell
		}
	}

	grid[u.y][u.x] = Snake

	game := &SnakeGame{
		gen:       rand.New(rand.NewSource(time.Now().UnixNano())),
		width:     width,
		height:    height,
		queue:     []unit{u},
		grid:      grid,
		score:     0,
		direction: Right,
	}

	game.generateFood()

	return game
}

// CanMove tells if the next move operation is allowed for the given direction.
func (g *SnakeGame) CanMove(direction Direction) bool {
	return g.direction.CanApply(direction)
}

// Move moves the snake to the given direction.
// Returns the current score. -1 means that the game is over.
func (g *SnakeGame) Move(direction Direction) error {
	n := len(g.queue)

	if !g.CanMove(direction) {
		return ErrUnacceptableDirection
	}

	g.direction = direction
	unit := g.queue[n-1].next(direction)

	// check if the boundary conditions are met.
	if unit.x >= g.width || unit.x < 0 || unit.y >= g.height || unit.y < 0 {
		return ErrBoundriesCross
	}

	needToGrow := false

	// if the next cell is a food cell, than we have to:
	// * increase the game score.
	// * try to regenerate a food.
	// * remember not to remove the snake tail (snake grow implementation).
	if g.grid[unit.y][unit.x] == Food {
		g.score++

		if n == len(g.queue)+1 {
			return ErrGameIsCompleted
		}

		g.generateFood()
		needToGrow = true
	}

	if !needToGrow {
		tail := g.queue[0]
		g.grid[tail.y][tail.x] = Cell
		g.queue = g.queue[1:]
	}

	// self-eating edge case
	if g.grid[unit.y][unit.x] == Snake {
		return ErrSelfEating
	}

	// Move snake to the next cell.
	g.queue = append(g.queue, unit)
	g.grid[unit.y][unit.x] = Snake

	return nil
}

// Grid returns the current grid.
func (g *SnakeGame) Grid() [][]Entity {
	return g.grid
}

// Score returns the current game score.
func (g *SnakeGame) Score() int {
	return g.score
}

func (g *SnakeGame) generateFood() {
	n := g.width * g.height

	for {
		food := g.gen.Intn(n)
		foodY, foodX := toXY(food, g.width, g.height)

		if g.grid[foodY][foodX] == Cell {
			g.grid[foodY][foodX] = Food
			return
		}
	}
}

func toXY(id int, w, h int) (int, int) {
	return id / h, id % w
}
