package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

const (
	WIDTH  = 20
	HEIGHT = 10
)

type Point struct {
	x, y int
}

func clearScreen() {
	cmd := exec.Command("clear")
	if os.Getenv("OS") == "Windows_NT" {
		cmd = exec.Command("cmd", "/c", "cls")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func createBoard() [][]rune {
	board := make([][]rune, HEIGHT)
	for i := range board {
		board[i] = make([]rune, WIDTH)
		for j := range board[i] {
			board[i][j] = ' '
		}
	}
	return board
}

func printBoard(board [][]rune, score int) {
	clearScreen()
	fmt.Println("+" + string(make([]rune, WIDTH, WIDTH+2)))
	for _, row := range board {
		fmt.Print("|")
		for _, cell := range row {
			if cell == 0 {
				fmt.Print(" ")
			} else {
				fmt.Printf("%c", cell)
			}
		}
		fmt.Println("|")
	}
	fmt.Println("+" + string(make([]rune, WIDTH, WIDTH+2)))
	fmt.Printf("Score: %d\n", score)
}

func placeFood(board [][]rune, snake []Point) Point {
	for {
		x := rand.Intn(WIDTH)
		y := rand.Intn(HEIGHT)
		isOnSnake := false
		for _, p := range snake {
			if p.x == x && p.y == y {
				isOnSnake = true
				break
			}
		}
		if !isOnSnake {
			board[y][x] = '*'
			return Point{x, y}
		}
	}
}

func moveSnake(snake []Point, direction string) []Point {
	head := snake[0]
	var newHead Point
	switch direction {
	case "UP":
		newHead = Point{head.x, head.y - 1}
	case "DOWN":
		newHead = Point{head.x, head.y + 1}
	case "LEFT":
		newHead = Point{head.x - 1, head.y}
	case "RIGHT":
		newHead = Point{head.x + 1, head.y}
	}
	newSnake := append([]Point{newHead}, snake...)
	return newSnake[:len(snake)]
}

func checkCollision(snake []Point) bool {
	head := snake[0]
	if head.x < 0 || head.x >= WIDTH || head.y < 0 || head.y >= HEIGHT {
		return true
	}
	for _, p := range snake[1:] {
		if head == p {
			return true
		}
	}
	return false
}

func gameLoop() {
	snake := []Point{{WIDTH / 2, HEIGHT / 2}}
	direction := "RIGHT"
	board := createBoard()
	food := placeFood(board, snake)
	score := 0

	// Goroutine za kontrolu unosa korisnika
	inputChan := make(chan string)
	go func() {
		var key string
		for {
			fmt.Scanln(&key)
			inputChan <- key
		}
	}()

	for {
		// Resetiraj ploču i postavi zmiju i hranu
		board = createBoard()
		for _, p := range snake {
			board[p.y][p.x] = 'O'
		}
		board[food.y][food.x] = '*'

		// Prikaži ploču
		printBoard(board, score)

		// Provjera unosa korisnika
		select {
		case key := <-inputChan:
			if key == "w" && direction != "DOWN" {
				direction = "UP"
			} else if key == "s" && direction != "UP" {
				direction = "DOWN"
			} else if key == "a" && direction != "RIGHT" {
				direction = "LEFT"
			} else if key == "d" && direction != "LEFT" {
				direction = "RIGHT"
			}
		default:
			// Nastavi bez unosa
		}

		// Pomakni zmiju
		snake = moveSnake(snake, direction)
		if checkCollision(snake) {
			fmt.Println("Game Over! Your score:", score)
			break
		}

		// Provjeri je li zmija pojela hranu
		if snake[0].x == food.x && snake[0].y == food.y {
			snake = append(snake, Point{}) // Dodaj segment
			food = placeFood(board, snake) // Postavi novu hranu
			score++
		}

		time.Sleep(200 * time.Millisecond)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	gameLoop()
}
