package game

import (
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

// Entity represents a grid entity.
type Entity byte

const (
	// Snake represents a part of a snake.
	Snake = 'o'
	// Cell represents an empty cell of the grid.
	Cell = '*'
	// Food represents a food on the grid.
	Food = 'x'
)

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
	gen    *rand.Rand
	width  int
	height int
	queue  []unit
	table  map[unit]bool
	grid   [][]Entity
	score  int
	food   int
}

// New creates a new snake game state.
func New(width int, height int) *SnakeGame {
	u := unit{x: 0, y: 0}
	grid := make([][]Entity, height)
	gen := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := range grid {
		grid[i] = make([]Entity, width)

		for j := range grid[i] {
			grid[i][j] = Cell
		}
	}

	grid[u.y][u.x] = Snake

	game := &SnakeGame{
		gen:    gen,
		food:   -1,
		width:  width,
		height: height,
		queue:  []unit{u},
		grid:   grid,
		score:  0,
	}

	game.generateFood()

	return game
}

// Move moves the snake to the given direction.
// Returns the current score. -1 means that the game is over.
func (g *SnakeGame) Move(direction Direction) int {
	n := len(g.queue)
	unit := g.queue[n-1].next(direction)

	if unit.x >= g.width || unit.x < 0 || unit.y >= g.height || unit.y < 0 {
		return -1
	}

	removeLast := true
	foodY, foodX := toXY(g.food, g.width, g.height)

	if foodX == unit.x && foodY == unit.y {
		g.score++

		if n == len(g.queue)+1 {
			return g.score
		}

		g.generateFood()
		removeLast = false
	}

	if removeLast {
		tail := g.queue[0]
		g.grid[tail.y][tail.x] = Cell
		g.queue = g.queue[1:]
	}

	if g.grid[unit.y][unit.x] == Snake {
		return -1
	}

	g.queue = append(g.queue, unit)
	g.grid[unit.y][unit.x] = Snake

	return g.score
}

// Grid returns the current grid.
func (g *SnakeGame) Grid() [][]Entity {
	return g.grid
}

func (g *SnakeGame) generateFood() {
	n := g.width * g.height

	for {
		g.food = g.gen.Intn(n)
		foodY, foodX := toXY(g.food, g.width, g.height)

		if g.grid[foodY][foodX] == Cell {
			g.grid[foodY][foodX] = Food
			return
		}
	}
}

func toXY(id int, w, h int) (int, int) {
	return id / h, id % w
}
