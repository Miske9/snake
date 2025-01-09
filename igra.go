package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"

	"github.com/eiannone/keyboard"
)

// Konstantne vrijednosti za veličinu ploče
const (
	WIDTH  = 20
	HEIGHT = 10
)

// Struktura koja predstavlja točku (x, y) na ploči
type Point struct {
	x, y int
}

// clearScreen briše ekran
func clearScreen() {
	cmd := exec.Command("clear")
	if os.Getenv("OS") == "Windows_NT" {
		cmd = exec.Command("cmd", "/c", "cls")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// createBoard kreira dvodimenzionalnu ploču s okvirima (#)
func createBoard() [][]rune {
	board := make([][]rune, HEIGHT+2)
	for i := range board {
		board[i] = make([]rune, WIDTH+2)
		for j := range board[i] {
			if i == 0 || i == HEIGHT+1 || j == 0 || j == WIDTH+1 {
				board[i][j] = '#' // Zidovi
			} else {
				board[i][j] = ' ' // Prazan prostor
			}
		}
	}
	return board
}

// printBoard iscrtava ploču i rezultat na konzoli
func printBoard(board [][]rune, score int) {
	clearScreen()
	for _, row := range board {
		for _, cell := range row {
			fmt.Printf("%c", cell)
		}
		fmt.Println()
	}
	fmt.Printf("Rezultat: %d\n", score)
}

// placeFood postavlja hranu na slučajnu poziciju na ploči koja nije zauzeta zmijom
func placeFood(board [][]rune, snake []Point) Point {
	for {
		x := rand.Intn(WIDTH) + 1
		y := rand.Intn(HEIGHT) + 1
		isOnSnake := false
		for _, p := range snake {
			if p.x == x && p.y == y {
				isOnSnake = true
				break
			}
		}
		if !isOnSnake {
			board[y][x] = '*' // Hrana
			return Point{x, y}
		}
	}
}

// moveSnake pomiče zmiju na novu poziciju prema trenutnom smjeru
func moveSnake(snake []Point, xSpeed, ySpeed int) []Point {
	head := snake[0]
	newHead := Point{head.x + xSpeed, head.y + ySpeed}
	return append([]Point{newHead}, snake[:len(snake)-1]...)
}

// checkCollision provjerava sudare zmije sa zidovima ili sobom
func checkCollision(snake []Point, board [][]rune) bool {
	head := snake[0]
	if board[head.y][head.x] == '#' { // Sudar sa zidom
		return true
	}
	for _, p := range snake[1:] {
		if head == p { // Sudar sa vlastitim tijelom
			return true
		}
	}
	return false
}

// gameLoop upravlja logikom jedne igre
func gameLoop() bool {
	// Inicijalizacija zmije, brzine i ploče
	snake := []Point{{WIDTH / 2, HEIGHT / 2}}
	xSpeed, ySpeed := 1, 0 // Početni smjer: desno
	board := createBoard()
	food := placeFood(board, snake)
	score := 0

	// Otvaranje praćenja tipkovnice
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer keyboard.Close()

	// Kreiranje kanala za unos
	inputChan := make(chan rune)
	go func() {
		for {
			if char, _, err := keyboard.GetKey(); err == nil {
				inputChan <- char
			}
		}
	}()

	for {
		// Pomiče zmiju u trenutnom smjeru
		snake = moveSnake(snake, xSpeed, ySpeed)

		// Provjerava sudar
		if checkCollision(snake, board) {
			printBoard(board, score)
			fmt.Println("Game Over!")
			fmt.Println("Double click 'r' for restart or 'q' for exit.")
			for {
				select {
				case char := <-inputChan:
					if char == 'r' {
						return true // Ponovno pokretanje
					} else if char == 'q' {
						return false // Izlaz
					}
				}
			}
		}

		// Provjerava je li zmija pojela hranu
		if snake[0].x == food.x && snake[0].y == food.y {
			snake = append(snake, snake[len(snake)-1]) // Dodaje dio zmije
			food = placeFood(board, snake)             // Nova hrana
			score++
		}

		// Resetira ploču i dodaje zmiju i hranu
		board = createBoard()
		for _, p := range snake {
			board[p.y][p.x] = 'O'
		}
		board[food.y][food.x] = '*'

		// Prikazuje ploču
		printBoard(board, score)

		// Pauza kako bi igra imala stabilnu brzinu
		time.Sleep(200 * time.Millisecond)

		// Obrada unosa za smjer kretanja
		select {
		case char := <-inputChan:
			switch char {
			case 'w':
				if ySpeed == 0 {
					xSpeed, ySpeed = 0, -1
				}
			case 's':
				if ySpeed == 0 {
					xSpeed, ySpeed = 0, 1
				}
			case 'a':
				if xSpeed == 0 {
					xSpeed, ySpeed = -1, 0
				}
			case 'd':
				if xSpeed == 0 {
					xSpeed, ySpeed = 1, 0
				}
			}
		default:
			// Nastavlja kretanje bez promjene
		}
	}
}

// main upravlja početkom i ponovnim pokretanjem igre
func main() {
	rand.Seed(time.Now().UnixNano()) // Postavlja seed za slučajne brojeve
	for gameLoop() {
		// Ponovno pokreće igru ako gameLoop vrati true
	}
}
